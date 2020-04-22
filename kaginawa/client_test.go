package kaginawa

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testAPIKey = "test123"

func TestNewClient_invalidEndpoint(t *testing.T) {
	_, err := NewClient("abc.def", testAPIKey)
	if err == nil {
		t.Error("expected error, got nil.")
	}
}

func TestNewClient_validEndpoint(t *testing.T) {
	client, err := NewClient("http://localhost:3000", testAPIKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Error("expected client object, got nil.")
	}
}

func TestNewClient_emptyAPIKey(t *testing.T) {
	_, err := NewClient("http://localhost:3000", "")
	if err == nil {
		t.Error("expected error, got nil.")
	}
}

func TestListAliveNodes(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization != "token "+testAPIKey {
			w.WriteHeader(http.StatusUnauthorized)
			t.Errorf("invalid api key: %s", authorization)
			return
		}
		expected, err := ioutil.ReadFile("testdata/nodes_multiple.json")
		if err != nil {
			t.Errorf("failed to initialize testdata: %v", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(expected); err != nil {
			t.Errorf("failed to write testdata response: %v", err)
			return
		}
	}))
	defer ts.Close()
	client, err := NewClient(ts.URL, testAPIKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	reports, err := client.ListAliveNodes(context.Background(), 0)
	if err != nil {
		t.Fatalf("unexpexted error: %v", err)
	}
	if len(reports) != 2 {
		t.Errorf("expected %d record, got %d record(s)", 2, len(reports))
	}
	for i := range reports {
		if reports[i].SSHServerHost != "example.com" {
			t.Errorf("%d: SSHServerHost expected %s, got %s", i, "example.com", reports[i].SSHServerHost)
		}
	}
}

func TestListNodesByCustomID(t *testing.T) {
	testCustomID := "test-mac"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization != "token "+testAPIKey {
			w.WriteHeader(http.StatusUnauthorized)
			t.Errorf("invalid api key: %s", authorization)
			return
		}
		customIDQuery := r.URL.Query().Get("custom-id")
		if customIDQuery != testCustomID {
			w.WriteHeader(http.StatusNotFound)
			t.Errorf("invalid custom id: %s", customIDQuery)
			return
		}
		expected, err := ioutil.ReadFile("testdata/nodes_single.json")
		if err != nil {
			t.Errorf("failed to initialize testdata: %v", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(expected); err != nil {
			t.Errorf("failed to write testdata response: %v", err)
			return
		}
	}))
	defer ts.Close()
	client, err := NewClient(ts.URL, testAPIKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	reports, err := client.ListNodesByCustomID(context.Background(), testCustomID)
	if err != nil {
		t.Fatalf("unexpexted error: %v", err)
	}
	if len(reports) != 1 {
		t.Errorf("expected %d record, got %d record(s)", 1, len(reports))
	}
}
