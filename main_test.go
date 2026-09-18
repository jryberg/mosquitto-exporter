package main

import "testing"

func TestProcessUpdateMetricType(t *testing.T) {
	tests := []struct {
		topic   string
		payload string
		counter bool
	}{
		{"$SYS/broker/connections/socket/count", "42", true},
		{"$SYS/broker/clients/expired", "3", true},
		{"$SYS/broker/packet/out/count", "7", false},
		{"$SYS/broker/packet/out/bytes", "1024", false},
		{"$SYS/broker/clients/total", "5", false},
	}

	for _, tt := range tests {
		t.Run(tt.topic, func(t *testing.T) {
			processUpdate(tt.topic, tt.payload)

			_, isCounter := counterMetrics[tt.topic]
			_, isGauge := gaugeMetrics[tt.topic]
			if isCounter != tt.counter || isGauge == tt.counter {
				t.Errorf("counter=%v gauge=%v, want counter=%v", isCounter, isGauge, tt.counter)
			}
		})
	}
}
