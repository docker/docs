package apiclient

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type NginxLoadBalancerStatus struct {
	NginxServers *NginxServers `json:"servers"`
}

type NginxServers struct {
	Total       int           `json:"total"`
	Generation  int           `json:"generation"`
	NginxServer []NginxServer `json:"server"`
}

type NginxServer struct {
	Index    int    `json:"index"`
	Upstream string `json:"upstream"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Rise     int    `json:"rise"`
	Fall     int    `json:"fall"`
	Type     string `json:"type"`
	Port     int    `json:"port"`
}

const loadBalancerStatusSubroute = "/load_balancer_status"

func (c *apiClient) LoadBalancerStatus() (*NginxLoadBalancerStatus, error) {
	resp, err := c.LoadBalancerStatusResponse()
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if err := validateStatusCode(resp, http.StatusOK); err != nil {
		return nil, err
	}

	var loadBalancerStatus NginxLoadBalancerStatus
	if err := json.NewDecoder(resp.Body).Decode(&loadBalancerStatus); err != nil {
		return nil, err
	}

	return &loadBalancerStatus, err
}

func (c *apiClient) LoadBalancerStatusResponse() (*http.Response, error) {
	v := url.Values{}
	v.Add("format", "json")
	return c.makeRequest("GET", url.URL{Path: loadBalancerStatusSubroute, RawQuery: v.Encode()}, nil)
}
