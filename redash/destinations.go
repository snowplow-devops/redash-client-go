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
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
)

// ErrUnsupportedDestinationType is wrapped into the error returned by
// ParseDestinationType (and anything that calls it) when a destination's
// "type" isn't one of the types this client knows how to parse into a
// concrete struct. The destination is still returned, as a generic
// *Destination, alongside this error - callers that only need the common
// fields (ID, Name, Type, Icon, raw Options) can use errors.Is to treat
// this as non-fatal and use the result anyway.
var ErrUnsupportedDestinationType = errors.New("unsupported destination type")

// destinationConstructors maps a destination's "type" to a constructor for
// its concrete struct. It's a package-level map of constructors, rather
// than of shared struct pointers, so each ParseDestinationType call gets
// its own fresh instance to unmarshal into.
var destinationConstructors = map[string]func() interface{}{
	"email":         func() interface{} { return &EmailDestination{} },
	"slack":         func() interface{} { return &SlackDestination{} },
	"webhook":       func() interface{} { return &WebhookDestination{} },
	"mattermost":    func() interface{} { return &MattermostDestination{} },
	"chatwork":      func() interface{} { return &ChatWorkDestination{} },
	"pagerduty":     func() interface{} { return &PagerDutyDestination{} },
	"hangouts_chat": func() interface{} { return &HangoutsChatDestination{} },
}

// ParseDestinationType parses payload and extracts into a suitable type.
// If the payload's "type" isn't recognized, it returns a generic
// *Destination and an error wrapping ErrUnsupportedDestinationType, rather
// than failing outright.
func ParseDestinationType(payload []byte) (destination interface{}, err error) {
	dst := &Destination{}
	err = json.Unmarshal(payload, dst)
	if err != nil {
		return nil, err
	}
	if dst.Options == nil {
		return nil, errors.New("empty options field")
	}
	newDestination, ok := destinationConstructors[dst.Type]
	if !ok {
		return dst, fmt.Errorf("%w: %q", ErrUnsupportedDestinationType, dst.Type)
	}
	destination = newDestination()
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
	defer func() { _ = response.Body.Close() }()
	body, _ := io.ReadAll(response.Body)

	destinations := []Destination{}
	err = json.Unmarshal(body, &destinations)
	if err != nil {
		return nil, err
	}

	return &destinations, nil
}

// GetDestination gets a specific Destination. If the destination's type
// isn't one this client recognizes, it still returns a generic
// *Destination alongside an error wrapping ErrUnsupportedDestinationType.
func (c *Client) GetDestination(id int) (interface{}, error) {
	path := "/api/destinations/" + strconv.Itoa(id)
	query := url.Values{}
	response, err := c.get(path, query)
	if err != nil {
		return nil, err
	}

	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	destination, err := ParseDestinationType(body)
	if err != nil && !errors.Is(err, ErrUnsupportedDestinationType) {
		return nil, err
	}

	return destination, err
}

// GetDestinationTypes gets all available types with configuration details
func (c *Client) GetDestinationTypes() (*[]DestinationTypes, error) {
	path := "/api/destinations/types"
	query := url.Values{}
	response, err := c.get(path, query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = response.Body.Close() }()
	body, _ := io.ReadAll(response.Body)

	destinationTypes := []DestinationTypes{}
	err = json.Unmarshal(body, &destinationTypes)
	if err != nil {
		return nil, err
	}

	return &destinationTypes, nil
}

// CreateDestination creates a new Destination. If the destination's type
// isn't one this client recognizes, it still creates it and returns a
// generic *Destination alongside an error wrapping
// ErrUnsupportedDestinationType.
func (c *Client) CreateDestination(payload []byte) (destination interface{}, err error) {
	path := "/api/destinations"

	_, err = ParseDestinationType(payload)
	if err != nil && !errors.Is(err, ErrUnsupportedDestinationType) {
		return nil, err
	}

	query := url.Values{}
	response, err := c.post(path, string(payload), query)
	if err != nil {
		return nil, err
	}

	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	destination, err = ParseDestinationType(body)
	if err != nil && !errors.Is(err, ErrUnsupportedDestinationType) {
		return nil, err
	}

	return destination, err
}

// UpdateDestination updates an existing Destination. If the destination's
// type isn't one this client recognizes, it still updates it and returns a
// generic *Destination alongside an error wrapping
// ErrUnsupportedDestinationType.
func (c *Client) UpdateDestination(id int, payload []byte) (destination interface{}, err error) {
	path := "/api/destinations/" + strconv.Itoa(id)

	_, err = ParseDestinationType(payload)
	if err != nil && !errors.Is(err, ErrUnsupportedDestinationType) {
		return nil, err
	}

	query := url.Values{}
	response, err := c.post(path, string(payload), query)
	if err != nil {
		return nil, err
	}

	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	destination, err = ParseDestinationType(body)
	if err != nil && !errors.Is(err, ErrUnsupportedDestinationType) {
		return nil, err
	}

	return destination, err
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
