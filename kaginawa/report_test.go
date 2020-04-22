package kaginawa

import (
	"testing"
	"time"
)

func TestReport_Timestamp(t *testing.T) {
	tests := []struct {
		input    Report
		expected time.Time
	}{
		{
			input:    Report{ServerTime: 0},
			expected: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
		},
		{
			input:    Report{ServerTime: 1600000000},
			expected: time.Date(2020, 9, 13, 12, 26, 40, 0, time.UTC).Local(),
		},
	}
	for i, d := range tests {
		actual := d.input.Timestamp()
		if d.expected != actual {
			t.Errorf("test %d: Timestamp() expected %v, got %v", i, d.expected, actual)
		}
	}
}

func TestReport_BootTimestamp(t *testing.T) {
	tests := []struct {
		input    Report
		expected time.Time
	}{
		{
			input:    Report{BootTime: 0},
			expected: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
		},
		{
			input:    Report{BootTime: 1600000000},
			expected: time.Date(2020, 9, 13, 12, 26, 40, 0, time.UTC).Local(),
		},
	}
	for i, d := range tests {
		actual := d.input.BootTimestamp()
		if d.expected != actual {
			t.Errorf("test %d: BootTimestamp() expected %v, got %v", i, d.expected, actual)
		}
	}
}
