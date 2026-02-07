// Copyright 2019 Altinity Ltd and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clickhouse

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	log "github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"

	core "k8s.io/api/core/v1"
	kube "k8s.io/client-go/kubernetes"

	api "github.com/altinity/clickhouse-operator/pkg/apis/clickhouse.altinity.com/v1"
	"github.com/altinity/clickhouse-operator/pkg/apis/common/types"
	"github.com/altinity/clickhouse-operator/pkg/apis/metrics"
	"github.com/altinity/clickhouse-operator/pkg/chop"
	chopAPI "github.com/altinity/clickhouse-operator/pkg/client/clientset/versioned"
	"github.com/altinity/clickhouse-operator/pkg/controller"
	chiNormalizer "github.com/altinity/clickhouse-operator/pkg/model/chi/normalizer"
	"github.com/altinity/clickhouse-operator/pkg/model/clickhouse"
	normalizerCommon "github.com/altinity/clickhouse-operator/pkg/model/common/normalizer"
)

// Exporter implements prometheus.Collector interface
type Exporter struct {
	collectorTimeout time.Duration

	// crInstallations is an index of watched CRs
	crInstallations crInstallationsIndex

	mutex               sync.RWMutex
	toRemoveFromWatched sync.Map
}

// Type compatibility
var _ prometheus.Collector = &Exporter{}

// NewExporter returns a new instance of Exporter type
func NewExporter(collectorTimeout time.Duration) *Exporter {
	return &Exporter{
		crInstallations:  newCRInstallationsIndex(),
		collectorTimeout: collectorTimeout,
	}
}

// getWatchedCHIs
func (e *Exporter) getWatchedCHIs() []*metrics.WatchedCR {
	return e.crInstallations.slice()
}

// Collect implements prometheus.Collector Collect method
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// Run cleanup on each collect
	e.cleanup()

	if ch == nil {
		log.Warning("Prometheus channel is closed. Unable to write metrics")
		return
	}

	start := time.Now()

	log.V(1).Info("Collect started")
	defer func() {
		log.V(1).Infof("Collect completed [%s]", time.Since(start))
	}()

	// Collection process should have limited duration
	ctx, cancel := context.WithTimeout(context.Background(), e.collectorTimeout)
	defer cancel()

	// This method may be called concurrently and must therefore be implemented in a concurrency safe way
	e.mutex.Lock()
	defer e.mutex.Unlock()

	log.V(1).Infof("Launching host collectors [%s]", time.Since(start))

	var wg = sync.WaitGroup{}
	e.crInstallations.walk(func(chi *metrics.WatchedCR, _ *metrics.WatchedCluster, host *metrics.WatchedHost) {
		wg.Add(1)
		go func(ctx context.Context, chi *metrics.WatchedCR, host *metrics.WatchedHost, ch chan<- prometheus.Metric) {
			defer wg.Done()
			e.collectHostMetrics(ctx, chi, host, ch)
		}(ctx, chi, host, ch)
	})
	wg.Wait()
}

// Describe implements prometheus.Collector Describe method
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(e, ch)
}

// enqueueToRemoveFromWatched
func (e *Exporter) enqueueToRemoveFromWatched(chi *metrics.WatchedCR) {
	e.toRemoveFromWatched.Store(chi, struct{}{})
}

// cleanup cleans all pending for cleaning
func (e *Exporter) cleanup() {
	// Clean up all pending for cleaning CHIs
	log.V(2).Info("Starting cleanup")
	e.toRemoveFromWatched.Range(func(key, value interface{}) bool {
		switch key.(type) {
		case *metrics.WatchedCR:
			e.toRemoveFromWatched.Delete(key)
			e.removeFromWatched(key.(*metrics.WatchedCR))
			log.V(1).Infof("Removed ClickHouseInstallation (%s/%s) from Exporter", key.(*metrics.WatchedCR).Name, key.(*metrics.WatchedCR).Namespace)
		}
		return true
	})
	log.V(2).Info("Completed cleanup")
}

// removeFromWatched deletes record from watched index
func (e *Exporter) removeFromWatched(chi *metrics.WatchedCR) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	log.V(1).Infof("Remove ClickHouseInstallation (%s/%s)", chi.Namespace, chi.Name)
	e.crInstallations.remove(chi.IndexKey())
}

// updateWatched updates watched index
func (e *Exporter) updateWatched(chi *metrics.WatchedCR) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	log.V(1).Infof("Update ClickHouseInstallation (%s/%s): %s", chi.Namespace, chi.Name, chi)
	e.crInstallations.set(chi.IndexKey(), chi)
}

// newFetcher returns new Metrics Fetcher for specified host
func (e *Exporter) newHostFetcher(host *metrics.WatchedHost) *MetricsFetcher {
	// Make base cluster connection params
	clusterConnectionParams := clickhouse.NewClusterConnectionParamsFromCHOpConfig(chop.Config())
	// Adjust base cluster connection params with per-host props
	switch clusterConnectionParams.Scheme {
	case api.ChSchemeAuto:
		switch {
		case types.IsPortAssigned(host.HTTPPort):
			clusterConnectionParams.Scheme = "http"
			clusterConnectionParams.Port = int(host.HTTPPort)
		case types.IsPortAssigned(host.HTTPSPort):
			clusterConnectionParams.Scheme = "https"
			clusterConnectionParams.Port = int(host.HTTPSPort)
		}
	case api.ChSchemeHTTP:
		clusterConnectionParams.Port = int(host.HTTPPort)
	case api.ChSchemeHTTPS:
		clusterConnectionParams.Port = int(host.HTTPSPort)
	}

	return NewMetricsFetcher(
		clusterConnectionParams.NewEndpointConnectionParams(host.Hostname),
		chop.Config().ClickHouse.Metrics.TablesRegexp,
	)
}

// collectHostMetrics collects metrics from one host and writes them into chan
func (e *Exporter) collectHostMetrics(ctx context.Context, chi *metrics.WatchedCR, host *metrics.WatchedHost, c chan<- prometheus.Metric) {
	collector := NewCollector(
		e.newHostFetcher(host),
		NewCHIPrometheusWriter(c, chi, host),
	)
	collector.CollectHostMetrics(ctx, host)
}

// getWatchedCHI serves HTTP request to get list of watched CHIs
func (e *Exporter) getWatchedCHI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(e.getWatchedCHIs())
}

// fetchCHI decodes chi from the request
func (e *Exporter) fetchCHI(r *http.Request) (*metrics.WatchedCR, error) {
	chi := &metrics.WatchedCR{}
	if err := json.NewDecoder(r.Body).Decode(chi); err == nil {
		if chi.IsValid() {
			return chi, nil
		}
	}

	return nil, fmt.Errorf("unable to parse CHI from request")
}

// updateWatchedCHI serves HTTP request to add CHI to the list of watched CHIs
func (e *Exporter) updateWatchedCHI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if chi, err := e.fetchCHI(r); err == nil {
		e.updateWatched(chi)
	} else {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
	}
}

// deleteWatchedCHI serves HTTP request to delete CHI from the list of watched CHIs
func (e *Exporter) deleteWatchedCHI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if chi, err := e.fetchCHI(r); err == nil {
		e.enqueueToRemoveFromWatched(chi)
	} else {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
	}
}

// DiscoveryWatchedCHIs discovers all ClickHouseInstallation objects available for monitoring and adds them to watched list
func (e *Exporter) DiscoveryWatchedCHIs(kubeClient kube.Interface, chopClient *chopAPI.Clientset) {
	// Get all CHI objects from watched namespace(s)
	watchedNamespace := chop.Config().GetInformerNamespace()
	list, err := chopClient.ClickhouseV1().ClickHouseInstallations(watchedNamespace).List(context.TODO(), controller.NewListOptions())
	if err != nil {
		log.V(1).Infof("Error read ClickHouseInstallations %v", err)
		return
	}
	if list == nil {
		return
	}

	// Walk over the list of ClickHouseInstallation objects and add them as watched
	for i := range list.Items {
		e.processDiscoveredCR(kubeClient, &list.Items[i])
	}
}

func (e *Exporter) processDiscoveredCR(kubeClient kube.Interface, chi *api.ClickHouseInstallation) {
	if !e.shouldWatchCR(chi) {
		log.V(1).Infof("Skip discovered CHI: %s/%s", chi.Namespace, chi.Name)
		return
	}

	log.V(1).Infof("Add discovered CHI: %s/%s", chi.Namespace, chi.Name)
	normalizer := chiNormalizer.New(func(namespace, name string) (*core.Secret, error) {
		return kubeClient.CoreV1().Secrets(namespace).Get(context.TODO(), name, controller.NewGetOptions())
	})

	normalized, _ := normalizer.CreateTemplated(chi, normalizerCommon.NewOptions[api.ClickHouseInstallation]())

	watchedCHI := metrics.NewWatchedCR(normalized)
	e.updateWatched(watchedCHI)
}

func (e *Exporter) shouldWatchCR(chi *api.ClickHouseInstallation) bool {
	if chi.IsStopped() {
		log.V(1).Infof("CHI %s/%s is stopped, unable to watch it", chi.Namespace, chi.Name)
		return false
	}

	return true
}
