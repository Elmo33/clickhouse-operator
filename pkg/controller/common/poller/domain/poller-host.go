// Copyright 2019 Altinity Ltd and/or its affiliates. All rights reserved.
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

package domain

import (
	"context"
	"fmt"

	log "github.com/altinity/clickhouse-operator/pkg/announcer"
	api "github.com/altinity/clickhouse-operator/pkg/apis/clickhouse.altinity.com/v1"
	"github.com/altinity/clickhouse-operator/pkg/chop"
	"github.com/altinity/clickhouse-operator/pkg/controller/common/poller"
	"github.com/altinity/clickhouse-operator/pkg/util"
)

// PollHost polls host
func PollHost(
	ctx context.Context,
	host *api.Host,
	isDoneFn func(ctx context.Context, host *api.Host) bool,
	_opts ...*poller.Options,
) error {
	if util.IsContextDone(ctx) {
		log.V(1).Info("poll host is done")
		return nil
	}

	opts := poller.NewOptions().FromConfig(chop.Config())
	for _, opt := range _opts {
		opts = opts.Merge(opt)
	}
	functions := &poller.Functions{
		IsDone: func(_ctx context.Context, _ any) bool {
			return isDoneFn(_ctx, host)
		},
	}
	caption := fmt.Sprintf("%s/%s", host.Runtime.Address.Namespace, host.Runtime.Address.HostName)

	return poller.New(ctx, caption).
		WithOptions(opts).
		WithFunctions(functions).
		Poll()
}
