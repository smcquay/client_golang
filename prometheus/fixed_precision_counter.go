// Copyright 2014 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package prometheus

import (
	"math"
	"sync/atomic"

	dto "github.com/prometheus/client_model/go"
)

// FixedPrecisionCounter implements a prometheus metric that uses atomic adds
// and stores for speed.
type FixedPrecisionCounter struct {
	val  int64
	prec uint

	desc *Desc
}

// NewCounter returns a populated counter.
func NewFixedPrecisionCounter(opts Opts, prec uint) Counter {
	desc := NewDesc(
		BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
		opts.Help,
		nil,
		opts.ConstLabels,
	)
	return &FixedPrecisionCounter{
		desc: desc,
		prec: uint(math.Pow10(int(prec))),
	}
}

// Desc returns a prometheus description for a counter.
func (c *FixedPrecisionCounter) Desc() *Desc {
	return c.desc
}

// Set stores the value in the counter.
func (c *FixedPrecisionCounter) Set(val float64) {
	atomic.StoreInt64(&c.val, int64(val)*int64(c.prec))
}

// Inc adds 1 to the counter.
func (c *FixedPrecisionCounter) Inc() {
	c.Add(1)
}

// Dec decrements 1 from the counter.
func (c *FixedPrecisionCounter) Dec() {
	c.Add(-1)
}

// Add generically adds delta to the value stored by counter.
func (c *FixedPrecisionCounter) Add(delta float64) {
	atomic.AddInt64(&c.val, int64(delta*float64(c.prec)))
}

// Sub is the invese of Add.
func (c *FixedPrecisionCounter) Sub(val float64) {
	c.Add(float64(val * -1))
}

func (c *FixedPrecisionCounter) Write(out *dto.Metric) error {
	f := float64(atomic.LoadInt64(&c.val)) / float64(c.prec)
	out.Counter = &dto.Counter{Value: &f}
	return nil
}

// Describe sends the counter's description to the chan
func (c *FixedPrecisionCounter) Describe(dc chan<- *Desc) {
	dc <- c.desc
}

// Collect sends the counter value to the chan
func (c *FixedPrecisionCounter) Collect(mc chan<- Metric) {
	mc <- c
}
