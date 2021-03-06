package apiclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const CoxEdgeAPIBase = "https://portal.coxedge.com/api/v1"
const CoxEdgeServiceCode = "edge-services"

type Client struct {
	apiKey     string
	HTTPClient *http.Client
}

func NewClient(apiKey string) Client {
	return Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		apiKey:     apiKey,
	}
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("MC-Api-Key", c.apiKey)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
