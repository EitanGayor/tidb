// Copyright 2021 PingCAP, Inc.
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

package expression

import (
	"testing"

	"github.com/pingcap/tidb/config"
	"github.com/pingcap/tidb/testkit/testmain"
	"github.com/pingcap/tidb/util/testbridge"
	"github.com/pingcap/tidb/util/timeutil"
	"github.com/tikv/client-go/v2/tikv"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	testbridge.WorkaroundGoCheckFlags()
	testmain.ShortCircuitForBench(m)

	config.UpdateGlobal(func(conf *config.Config) {
		conf.TiKVClient.AsyncCommit.SafeWindow = 0
		conf.TiKVClient.AsyncCommit.AllowedClockDrift = 0
	})
	tikv.EnableFailpoints()

	// Some test depends on the values of timeutil.SystemLocation()
	// If we don't SetSystemTZ() here, the value would change unpredictable.
	// Affected by the order whether a testsuite runs before or after integration test.
	// Note, SetSystemTZ() is a sync.Once operation.
	timeutil.SetSystemTZ("system")

	opts := []goleak.Option{
		goleak.IgnoreTopFunction("go.etcd.io/etcd/pkg/logutil.(*MergeLogger).outputLoop"),
		goleak.IgnoreTopFunction("go.opencensus.io/stats/view.(*worker).start"),
	}

	goleak.VerifyTestMain(m, opts...)
}
