package kaginawa

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testAPIKey = "test123"

func TestNewClientWithInvalidEndpoint(t *testing.T) {
	_, err := NewClient("abc.def", testAPIKey)
	if err == nil {
		t.Error("expected error, got nil.")
	}
}

func TestNewClientWithEmptyEndpoint(t *testing.T) {
	_, err := NewClient("", testAPIKey)
	if err == nil {
		t.Error("expected error, got nil.")
	}
}

func TestNewClientWithValidEndpoint(t *testing.T) {
	client, err := NewClient("http://localhost:3000", testAPIKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Error("expected client object, got nil.")
	}
}

func TestNewClientWithEmptyAPIKey(t *testing.T) {
	_, err := NewClient("http://localhost:3000", "")
	if err == nil {
		t.Error("expected error, got nil.")
	}
}

func TestFindNode(t *testing.T) {
	testID := "f0:18:98:eb:c7:27"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization != "token "+testAPIKey {
			w.WriteHeader(http.StatusUnauthorized)
			t.Errorf("invalid api key: %s", authorization)
			return
		}
		if !strings.HasSuffix(r.URL.Path, testID) {
			w.WriteHeader(http.StatusNotFound)
			t.Errorf("invalid path: %s", r.URL.Path)
			return
		}
		expected, err := ioutil.ReadFile("testdata/node.json")
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
	report, err := client.FindNode(context.Background(), testID)
	if err != nil {
		t.Fatalf("unexpexted error: %v", err)
	}
	if report.ID != testID {
		t.Errorf("expected ID %s, got %s", testID, report.ID)
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

func TestListHistories(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization != "token "+testAPIKey {
			w.WriteHeader(http.StatusUnauthorized)
			t.Errorf("invalid api key: %s", authorization)
			return
		}
		expected, err := ioutil.ReadFile("testdata/histories_multiple.json")
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
	reports, err := client.ListHistories(context.Background(), "b8:27:eb:36:83:e0", 1500000000, 1600000000)
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

func TestFindSSHServerByHostname(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization != "token "+testAPIKey {
			w.WriteHeader(http.StatusUnauthorized)
			t.Errorf("invalid api key: %s", authorization)
			return
		}
		expected, err := ioutil.ReadFile("testdata/servers_single.json")
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
	server, err := client.FindSSHServerByHostname(context.Background(), "example.com")
	if err != nil {
		t.Fatalf("unexpexted error: %v", err)
	}
	if server == nil {
		t.Fatal("expected nil, got an error.")
	}
	if server.Host != "example.com" {
		t.Errorf("Host expected %s, got %s", "example.com", server.Host)
	}
	if server.Port != 22 {
		t.Errorf("Port expected %d, got %d", 22, server.Port)
	}
	if server.User != "kaginawa" {
		t.Errorf("User expected %s, got %s", "kaginawa", server.User)
	}
	if server.Password != "test-pw" {
		t.Errorf("Password expected %s, got %s", "test-pw", server.Password)
	}
	if server.Key != "" {
		t.Errorf("Key expected %s, got %s", "(empty)", server.Key)
	}
}
