package kaginawa

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"
)

func TestTimestamp(t *testing.T) {
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

func TestBootTimestamp(t *testing.T) {
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

func TestUnmarshalReport(t *testing.T) {
	raw, err := ioutil.ReadFile("testdata/node.json")
	if err != nil {
		t.Fatalf("failed to initialize testdata: %v", err)
	}
	var report Report
	if err = json.Unmarshal(raw, &report); err != nil {
		t.Fatalf("failed to unmarshal testdata: %v", err)
	}
	if report.Payload != "{\"ip\":\"126.74.231.234\"}" {
		t.Errorf("unexpected payload: %s", report.Payload)
	}
	if report.PayloadCmd != "curl https://api.ipify.org?format=json" {
		t.Errorf("unexpected payload cmd: %s", report.PayloadCmd)
	}
}
