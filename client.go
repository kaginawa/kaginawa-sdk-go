package kaginawa

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	nodesResource   = "/nodes"
	serversResource = "/servers"
)

// Client is a Kaginawa Server REST API client.
type Client struct {
	endpoint   string
	apiKey     string
	client     http.Client
	closeError error
}

// NewClient will creates Kaginawa client object.
func NewClient(endpoint, apiKey string) (*Client, error) {
	if len(endpoint) == 0 {
		return nil, errors.New("most specify an endpoint")
	}
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		return nil, fmt.Errorf("not an http or https endpoint: %s", endpoint)
	}
	if len(apiKey) == 0 {
		return nil, errors.New("most specify an api key")
	}
	return &Client{
		endpoint: endpoint,
		apiKey:   apiKey,
		client:   http.Client{},
	}, nil
}

func (c *Client) request(ctx context.Context, method, url string, body io.Reader, expectedStatus int) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "token "+c.apiKey)
	req.URL.Query()
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	if resp.StatusCode != expectedStatus {
		return nil, fmt.Errorf("kaginawa server respond HTTP %s", resp.Status)
	}
	return resp, nil
}

func (c *Client) safeClose(closer io.Closer) {
	c.closeError = closer.Close()
}

// FindNode finds a report by id.
func (c *Client) FindNode(ctx context.Context, id string) (*Report, error) {
	resp, err := c.request(ctx, http.MethodGet, c.endpoint+nodesResource+"/"+id, nil, http.StatusOK)
	if err != nil {
		return nil, err
	}
	defer c.safeClose(resp.Body)
	var report Report
	if err := json.NewDecoder(resp.Body).Decode(&report); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	return &report, nil
}

// ListAliveNodes queries list of all alive nodes.
// It is considered alive if the last received time is within 5 minutes.
func (c *Client) ListAliveNodes(ctx context.Context, thresholdMin int) ([]Report, error) {
	values := url.Values{"projection": {"id"}}
	if thresholdMin > 0 {
		values.Add("minutes", strconv.Itoa(thresholdMin))
	}
	resp, err := c.request(ctx, http.MethodGet, c.endpoint+nodesResource+"?"+values.Encode(), nil, http.StatusOK)
	if err != nil {
		return nil, err
	}
	defer c.safeClose(resp.Body)
	var reports []Report
	if err := json.NewDecoder(resp.Body).Decode(&reports); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	return reports, nil
}

// ListNodesByCustomID queries list of reports by custom-id.
func (c *Client) ListNodesByCustomID(ctx context.Context, customID string) ([]Report, error) {
	values := url.Values{"custom-id": {customID}}
	resp, err := c.request(ctx, http.MethodGet, c.endpoint+nodesResource+"?"+values.Encode(), nil, http.StatusOK)
	if err != nil {
		return nil, err
	}
	defer c.safeClose(resp.Body)
	var reports []Report
	if err := json.NewDecoder(resp.Body).Decode(&reports); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	return reports, nil
}

// ListHistories queries list of histories by id.
func (c *Client) ListHistories(ctx context.Context, id string, beginTimestamp, endTimestamp int64) ([]Report, error) {
	values := url.Values{"projection": {"measurement"}}
	if beginTimestamp > 0 {
		values.Add("begin", strconv.FormatInt(beginTimestamp, 10))
	}
	if endTimestamp > 0 {
		values.Add("end", strconv.FormatInt(endTimestamp, 10))
	}
	path := fmt.Sprintf("%s%s/%s/histories?%s", c.endpoint, nodesResource, id, values.Encode())
	resp, err := c.request(ctx, http.MethodGet, path, nil, http.StatusOK)
	if err != nil {
		return nil, err
	}
	defer c.safeClose(resp.Body)
	var reports []Report
	if err := json.NewDecoder(resp.Body).Decode(&reports); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	return reports, nil
}

// FindSSHServerByHostname finds a SSH server entry by hostname.
func (c *Client) FindSSHServerByHostname(ctx context.Context, hostname string) (*SSHServer, error) {
	resp, err := c.request(ctx, http.MethodGet, c.endpoint+serversResource+"/"+hostname, nil, http.StatusOK)
	if err != nil {
		return nil, err
	}
	defer c.safeClose(resp.Body)
	var server SSHServer
	if err := json.NewDecoder(resp.Body).Decode(&server); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	return &server, nil
}
