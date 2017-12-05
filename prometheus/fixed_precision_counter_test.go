package prometheus

import (
	"testing"

	dto "github.com/prometheus/client_model/go"
)

func TestFixedPrecisionCounterAdd(t *testing.T) {
	counter := NewFixedPrecisionCounter(Opts{
		Name: "test",
		Help: "test help",
	}, 3).(*FixedPrecisionCounter)

	counter.Inc()
	var want int64 = 1000
	if expected, got := want, counter.val; expected != got {
		t.Errorf("Expected %f, got %f.", expected, got)
	}
	counter.Add(42.3)
	want = 43300
	if expected, got := want, counter.val; expected != got {
		t.Errorf("Expected %f, got %f.", expected, got)
	}

	counter.Sub(3.2)
	want = 40100
	if expected, got := want, counter.val; expected != got {
		t.Errorf("Expected %f, got %f.", expected, got)
	}

	m := &dto.Metric{}
	counter.Write(m)

	if expected, got := `counter:<value:40.1 > `, m.String(); expected != got {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
