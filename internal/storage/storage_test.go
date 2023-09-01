package storage

import "testing"

func TestSaveGauge(t *testing.T) {
	type metrics struct {
		name string
		val  float64
	}
	tests := []struct {
		name    string
		storage *Storage
		metrics metrics
		ok      bool
		want    float64
	}{
		{
			name:    "simple test on saving gauge metrics",
			storage: NewStorage(),
			metrics: metrics{
				name: "Heap",
				val:  5.0,
			},
			ok:   true,
			want: 5.0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.storage.SaveGauge(test.metrics.name, test.metrics.val)
			val, ok := test.storage.gauge[test.metrics.name]
			if ok != test.ok {
				t.Errorf("Want: %T, got: %T", test.ok, ok)
			}
			if val != test.want {
				t.Errorf("Want: %f, got: %f", test.want, val)
			}
		})
	}
}

func TestUpdateCounter(t *testing.T) {
	storage := NewStorage()
	type metrics struct {
		name string
		val  int64
	}
	tests := []struct {
		name    string
		metrics metrics
		ok      bool
		want    int64
	}{
		{
			name: "simple test on saving counter metrics",
			metrics: metrics{
				name: "PollCounter",
				val:  1,
			},
			ok:   true,
			want: 1,
		},
		{
			name: "simple test on saving counter metrics",
			metrics: metrics{
				name: "PollCounter",
				val:  2,
			},
			ok:   true,
			want: 3,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.UpdateCounter(test.metrics.name, test.metrics.val)
			val, ok := storage.counter[test.metrics.name]
			if ok != test.ok {
				t.Errorf("Want: %T, got: %T", test.ok, ok)
			}
			if val != test.want {
				t.Errorf("Want: %d, got: %d", test.want, val)
			}
		})
	}
}
