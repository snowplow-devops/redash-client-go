//
// Copyright (c) 2020 Snowplow Analytics Ltd. All rights reserved.
//
// This program is licensed to you under the Apache License Version 2.0,
// and you may not use this file except in compliance with the Apache License Version 2.0.
// You may obtain a copy of the Apache License Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0.
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the Apache License Version 2.0 is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the Apache License Version 2.0 for the specific language governing permissions and limitations there under.
//

package redash

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Client contains an active Redash API client
type Client struct {
	Config *Config
}

// Config holds the neccesary setup vars
type Config struct {
	RedashURI  string
	APIKey     string
	StrictMode bool
}

// NewClient returns a *Client from a valid *Config
func NewClient(config *Config) (*Client, error) {
	redashURI, err := url.ParseRequestURI(config.RedashURI)
	if err != nil {
		return nil, fmt.Errorf("Missing or invalid RedashURI")
	}

	if redashURI.Scheme != "http" && redashURI.Scheme != "https" {
		return nil, fmt.Errorf("Only HTTP(S) URIs allowed")
	}

	if config.APIKey == "" {
		return nil, fmt.Errorf("Missing APIKey")
	}

	c := &Client{Config: config}
	return c, nil
}

// IsStrict returns true if StrictMode is set. This currently causes
// data_source creates/updates to fail if extraneous properties
// are present in the payload.
func (c *Client) IsStrict() bool {
	return c.Config.StrictMode
}

func (c *Client) doRequest(method, path, body string) (*http.Response, error) {
	requestURI := strings.TrimSuffix(c.Config.RedashURI, "/") + path

	log.Debug(fmt.Sprintf("[DEBUG] %s request to %s", method, path))

	response, err := func() (*http.Response, error) {
		request, err := http.NewRequest(method, requestURI, strings.NewReader(body))
		if err != nil {
			return nil, err
		}

		request.Header.Add("Content-Type", "application/json")
		request.Header.Set("Authorization", "Key "+c.Config.APIKey)

		return http.DefaultClient.Do(request)
	}()
	if err != nil {
		return nil, err
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, fmt.Errorf("HTTP Response: %d", response.StatusCode)
	}

	return response, nil
}

func (c *Client) get(path string) (*http.Response, error) {
	return c.doRequest(http.MethodGet, path, "")
}

func (c *Client) post(path string, payload string) (*http.Response, error) {
	return c.doRequest(http.MethodPost, path, payload)
}

func (c *Client) put(path string, payload string) (*http.Response, error) {
	return c.doRequest(http.MethodPut, path, payload)
}

func (c *Client) delete(path string) (*http.Response, error) {
	return c.doRequest(http.MethodDelete, path, "")
}
