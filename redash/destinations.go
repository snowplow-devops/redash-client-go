//
// Copyright (c) 2020-2022 Snowplow Analytics Ltd. All rights reserved.
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// ParseDestinationType parses payload and extracts into a suitable type
func ParseDestinationType(payload []byte) (destination interface{}, err error) {
	dst := &Destination{}
	err = json.Unmarshal(payload, dst)
	if err != nil {
		return nil, err
	}
	dstTypes := map[string]interface{}{
		"email":         &EmailDestination{},
		"slack":         &SlackDestination{},
		"webhook":       &WebhookDestination{},
		"mattermost":    &MattermostDestination{},
		"chatwork":      &ChatWorkDestination{},
		"pagerduty":     &PagerDutyDestination{},
		"hangouts_chat": &HangoutsChatDestination{},
	}
	destination, ok := dstTypes[dst.Type]
	if !ok {
		log.Errorf("Invalid destination %s", dst.Type)
	}
	err = json.Unmarshal(payload, destination)
	if err != nil {
		return nil, err
	}

	return destination, nil
}

// GetDestinations gets an array of all Destinations available
func (c *Client) GetDestinations() (*[]Destination, error) {
	path := "/api/destinations"
	query := url.Values{}
	response, err := c.get(path, query)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	destinations := []Destination{}
	err = json.Unmarshal(body, &destinations)
	if err != nil {
		return nil, err
	}

	return &destinations, nil
}

// GetDestination gets a specific Destination
func (c *Client) GetDestination(id int) (destination interface{}, err error) {
	path := "/api/destinations/" + strconv.Itoa(id)
	query := url.Values{}
	response, err := c.get(path, query)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	destination, err = ParseDestinationType(body)
	if err != nil {
		return nil, err
	}

	return destination, nil
}

// GetDestinationTypes gets all available types with configuration details
func (c *Client) GetDestinationTypes() (*[]DestinationCommon, error) {
	path := "/api/destinations/types"
	query := url.Values{}
	response, err := c.get(path, query)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	destinationTypes := []DestinationCommon{}
	err = json.Unmarshal(body, &destinationTypes)
	if err != nil {
		return nil, err
	}

	return &destinationTypes, nil
}

// CreateDestination creates a new Destination
func (c *Client) CreateDestination(payload []byte) (destination interface{}, err error) {
	path := "/api/destinations"

	_, err = ParseDestinationType(payload)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	response, err := c.post(path, string(payload), query)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	destination = Destination{}

	err = json.Unmarshal(body, &destination)
	if err != nil {
		return nil, err
	}

	return &destination, nil
}

// UpdateDestination Updates an existing Destination
func (c *Client) UpdateDestination(id int, payload []byte) (resp interface{}, err error) {
	path := "/api/destinations/" + strconv.Itoa(id)

	destination, err := ParseDestinationType(payload)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	response, err := c.post(path, fmt.Sprint(&destination), query)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	resp, err = ParseDestinationType(body)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// DeleteDestination deletes a specific Destination
func (c *Client) DeleteDestination(id int) error {
	path := "/api/destinations/" + strconv.Itoa(id)

	query := url.Values{}
	_, err := c.delete(path, query)
	if err != nil {
		return err
	}

	return nil
}
