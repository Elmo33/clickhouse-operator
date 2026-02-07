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
	"sync"

	log "github.com/golang/glog"

	"github.com/altinity/clickhouse-operator/pkg/apis/metrics"
)

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

// Add adds or updates a CR in the registry
func (r *CRRegistry) Add(cr *metrics.WatchedCR) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	log.V(1).Infof("Registry: Add CR (%s/%s): %s", cr.Namespace, cr.Name, cr)
	r.index.set(cr.IndexKey(), cr)
}

// Remove removes a CR from the registry
func (r *CRRegistry) Remove(cr *metrics.WatchedCR) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	log.V(1).Infof("Registry: Remove CR (%s/%s)", cr.Namespace, cr.Name)
	r.index.remove(cr.IndexKey())
}

// EnqueueRemove enqueues a CR for removal (will be removed on next Cleanup)
func (r *CRRegistry) EnqueueRemove(cr *metrics.WatchedCR) {
	r.toRemove.Store(cr, struct{}{})
}

// Cleanup processes all CRs enqueued for removal
func (r *CRRegistry) Cleanup() {
	log.V(2).Info("Registry: Starting cleanup")
	r.toRemove.Range(func(key, value interface{}) bool {
		if cr, ok := key.(*metrics.WatchedCR); ok {
			r.toRemove.Delete(key)
			r.Remove(cr)
			log.V(1).Infof("Registry: Cleaned up CR (%s/%s)", cr.Name, cr.Namespace)
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
