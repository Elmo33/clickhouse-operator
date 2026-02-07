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
	"fmt"
	"sync"

	log "github.com/golang/glog"

	"github.com/altinity/clickhouse-operator/pkg/apis/metrics"
)

// removeType identifies what type of entity to remove
type removeType int

const (
	removeTypeCR removeType = iota
	removeTypeHost
)

// removeRequest represents a pending removal request
type removeRequest struct {
	removeType removeType
	cr         *metrics.WatchedCR
	host       *HostRequest
}

// CRRegistry is a thread-safe storage for watched Custom Resources
type CRRegistry struct {
	index    crInstallationsIndex
	mutex    sync.RWMutex
	toRemove sync.Map
}

// NewCRRegistry creates a new CRRegistry instance
func NewCRRegistry() *CRRegistry {
	return &CRRegistry{
		index: newCRInstallationsIndex(),
	}
}

// AddCR adds or updates a CR in the registry
func (r *CRRegistry) AddCR(cr *metrics.WatchedCR) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	log.V(1).Infof("Registry: Add CR (%s/%s): %s", cr.Namespace, cr.Name, cr)
	r.index.set(cr.IndexKey(), cr)
}

// AddHost adds a host to an existing CR in the registry
func (r *CRRegistry) AddHost(req *HostRequest) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	crKey := (&metrics.WatchedCR{Namespace: req.CRNamespace, Name: req.CRName}).IndexKey()
	cr, ok := r.index.get(crKey)
	if !ok || cr == nil {
		return fmt.Errorf("CR not found: %s", crKey)
	}

	// Find or create cluster
	var cluster *metrics.WatchedCluster
	for _, c := range cr.Clusters {
		if c.Name == req.ClusterName {
			cluster = c
			break
		}
	}
	if cluster == nil {
		cluster = &metrics.WatchedCluster{Name: req.ClusterName}
		cr.Clusters = append(cr.Clusters, cluster)
	}

	// Add or update host
	found := false
	for i, h := range cluster.Hosts {
		if h.Hostname == req.Host.Hostname {
			cluster.Hosts[i] = req.Host
			found = true
			break
		}
	}
	if !found {
		cluster.Hosts = append(cluster.Hosts, req.Host)
	}

	log.V(1).Infof("Registry: Add Host %s to CR (%s/%s) cluster %s", req.Host.Hostname, req.CRNamespace, req.CRName, req.ClusterName)
	return nil
}

// RemoveCR removes a CR from the registry
func (r *CRRegistry) RemoveCR(cr *metrics.WatchedCR) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	log.V(1).Infof("Registry: Remove CR (%s/%s)", cr.Namespace, cr.Name)
	r.index.remove(cr.IndexKey())
}

// RemoveHost removes a host from a CR in the registry
func (r *CRRegistry) RemoveHost(req *HostRequest) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	crKey := (&metrics.WatchedCR{Namespace: req.CRNamespace, Name: req.CRName}).IndexKey()
	cr, ok := r.index.get(crKey)
	if !ok || cr == nil {
		log.V(1).Infof("Registry: Cannot remove host, CR not found: %s", crKey)
		return
	}

	for _, cluster := range cr.Clusters {
		if cluster.Name == req.ClusterName {
			for i, h := range cluster.Hosts {
				if h.Hostname == req.Host.Hostname {
					cluster.Hosts = append(cluster.Hosts[:i], cluster.Hosts[i+1:]...)
					log.V(1).Infof("Registry: Remove Host %s from CR (%s/%s) cluster %s", req.Host.Hostname, req.CRNamespace, req.CRName, req.ClusterName)
					return
				}
			}
		}
	}
}

// EnqueueRemoveCR enqueues a CR for removal (will be removed on next Cleanup)
func (r *CRRegistry) EnqueueRemoveCR(cr *metrics.WatchedCR) {
	req := &removeRequest{removeType: removeTypeCR, cr: cr}
	r.toRemove.Store(req, struct{}{})
}

// EnqueueRemoveHost enqueues a host for removal (will be removed on next Cleanup)
func (r *CRRegistry) EnqueueRemoveHost(host *HostRequest) {
	req := &removeRequest{removeType: removeTypeHost, host: host}
	r.toRemove.Store(req, struct{}{})
}

// Cleanup processes all CRs and hosts enqueued for removal
func (r *CRRegistry) Cleanup() {
	log.V(2).Info("Registry: Starting cleanup")
	r.toRemove.Range(func(key, value interface{}) bool {
		r.toRemove.Delete(key)
		if req, ok := key.(*removeRequest); ok {
			switch req.removeType {
			case removeTypeCR:
				r.RemoveCR(req.cr)
				log.V(1).Infof("Registry: Cleaned up CR (%s/%s)", req.cr.Namespace, req.cr.Name)
			case removeTypeHost:
				r.RemoveHost(req.host)
				log.V(1).Infof("Registry: Cleaned up Host %s from CR (%s/%s)", req.host.Host.Hostname, req.host.CRNamespace, req.host.CRName)
			}
		}
		return true
	})
	log.V(2).Info("Registry: Completed cleanup")
}

// List returns all watched CRs as a slice
func (r *CRRegistry) List() []*metrics.WatchedCR {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.index.slice()
}

// Walk iterates over all hosts while holding an exclusive lock
// Use this when the iteration may modify state
func (r *CRRegistry) Walk(fn func(*metrics.WatchedCR, *metrics.WatchedCluster, *metrics.WatchedHost)) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.index.walk(fn)
}
