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
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/altinity/clickhouse-operator/pkg/apis/metrics"
)

// RESTServer provides HTTP API for managing watched CRs
type RESTServer struct {
	registry *CRRegistry
}

// NewRESTServer creates a new RESTServer instance
func NewRESTServer(registry *CRRegistry) *RESTServer {
	return &RESTServer{
		registry: registry,
	}
}

// ServeHTTP implements http.Handler interface
func (s *RESTServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/chi" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.handleGet(w, r)
	case http.MethodPost:
		s.handlePost(w, r)
	case http.MethodDelete:
		s.handleDelete(w, r)
	default:
		_, _ = fmt.Fprintf(w, "Sorry, only GET, POST and DELETE methods are supported.")
	}
}

// handleGet serves HTTP GET request to get list of watched CRs
func (s *RESTServer) handleGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(s.registry.List())
}

// handlePost serves HTTP POST request to add CR to the list of watched CRs
func (s *RESTServer) handlePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if cr, err := s.decodeCR(r); err == nil {
		s.registry.Add(cr)
	} else {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
	}
}

// handleDelete serves HTTP DELETE request to delete CR from the list of watched CRs
func (s *RESTServer) handleDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if cr, err := s.decodeCR(r); err == nil {
		s.registry.EnqueueRemove(cr)
	} else {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
	}
}

// decodeCR decodes CR from the HTTP request body
func (s *RESTServer) decodeCR(r *http.Request) (*metrics.WatchedCR, error) {
	cr := &metrics.WatchedCR{}
	if err := json.NewDecoder(r.Body).Decode(cr); err == nil {
		if cr.IsValid() {
			return cr, nil
		}
	}
	return nil, fmt.Errorf("unable to parse CR from request")
}

// StartMetricsREST starts Prometheus metrics exporter and REST API server
func StartMetricsREST(
	metricsAddress string,
	metricsPath string,
	collectorTimeout time.Duration,
	chiListAddress string,
	chiListPath string,
) (*Exporter, *CRRegistry) {
	log.V(1).Infof("Starting metrics exporter at '%s%s'\n", metricsAddress, metricsPath)

	// Create shared registry
	registry := NewCRRegistry()

	// Create and register Prometheus exporter
	exporter := NewExporter(registry, collectorTimeout)
	prometheus.MustRegister(exporter)

	// Create REST server
	restServer := NewRESTServer(registry)

	// Setup HTTP handlers
	http.Handle(metricsPath, promhttp.Handler())
	http.Handle(chiListPath, restServer)

	// Start HTTP servers
	go http.ListenAndServe(metricsAddress, nil)
	if metricsAddress != chiListAddress {
		go http.ListenAndServe(chiListAddress, nil)
	}

	return exporter, registry
}
