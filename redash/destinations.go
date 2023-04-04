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
	// log "github.com/sirupsen/logrus"
)

// Type of Destination
type Destination struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	Icon string `json:"icon,omitempty"`
}

// DestinationType stucture
type DestinationTypes struct {
	Destination
	ConfigurationSchema interface{}
	//interface{} `json:"configuration_schema"`
}

type EmailDestination struct {
	Destination
	ConfgurationSchema struct {
		Type       string `json:"type,omitempty"`
		Properties struct {
			Addresses struct {
				Type string `json:"type,omitempty"`
			} `json:"addresses"`
			SubjectTemplate struct {
				Type    string `json:"type,omitempty"`
				Default string `json:"default,omitempty"`
				Title   string `json:"title,omitempty"`
			} `json:"subject_template"`
		} `json:"properties"`
		Required     []string `json:"required,omitempty"`
		ExtraOptions []string `json:"extra_options,omitempty"`
	} `json:"configuration_schema,omitempty"`
}
type SlackDestination struct {
	Destination
	ConfgurationSchema struct {
		Type       string `json:"type"`
		Properties struct {
			URL struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"url,omitempty"`
			Username struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"username,omitempty"`
			IconEmoji struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"icon_emoji,omitempty"`
			IconURL struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"icon_url,omitempty"`
			Channel struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"channel,omitempty"`
		} `json:"properties,omitempty"`
		Secret []string `json:"secret,omitempty"`
	} `json:"configuration_schema,omitempty"`
}

type WebhookDestination struct {
	Destination
	ConfgurationSchema struct {
		Type       string `json:"type"`
		Properties struct {
			URL struct {
				Type string `json:"type,omitempty"`
			} `json:"url,omitempty"`
			Username struct {
				Type string `json:"type,omitempty"`
			} `json:"username,omitempty"`
			Password struct {
				Type string `json:"type,omitempty"`
			} `json:"password,omitempty"`
		} `json:"properties,omitempty"`
		Required []string `json:"required,omitempty"`
		Secret   []string `json:"secret,omitempty"`
	} `json:"configuration_schema,omitempty"`
}

type HipChatDestination struct {
	Destination
	ConfgurationSchema struct {
		Type       string `json:"type"`
		Properties struct {
			URL struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"url,omitempty"`
		} `json:"properties,omitempty"`
		Secret   []string `json:"secret,omitempty"`
		Required []string `json:"required,omitempty"`
	} `json:"configuration_schema,omitempty"`
	Deprecated bool `json:"deprecated,omitempty"`
}
type MattermostDestination struct {
	Destination
	ConfgurationSchema struct {
		Type       string `json:"type"`
		Properties struct {
			URL struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"url,omitempty"`
			Username struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"username,omitempty"`
			IconURL struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"icon_url,omitempty"`
			Channel struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"channel,omitempty"`
		} `json:"properties,omitempty"`
		Secret string `json:"secret,omitempty"`
	} `json:"configuration_schema,omitempty"`
}
type ChatWorkDestination struct {
	Destination
	ConfgurationSchema struct {
		Type       string `json:"type"`
		Properties struct {
			APIToken struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"api_token,omitempty"`
			RoomID struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"room_id,omitempty"`
			MessageTemplate struct {
				Type    string `json:"type,omitempty"`
				Default string `json:"default,omitempty"`
				Title   string `json:"title,omitempty"`
			} `json:"message_template,omitempty"`
		} `json:"properties,omitempty"`
		Secret   []string `json:"secret,omitempty"`
		Required []string `json:"required,omitempty"`
	} `json:"configuration_schema,omitempty"`
}
type PagerDutyDestination struct {
	Destination
	ConfgurationSchema struct {
		Type       string `json:"type"`
		Properties struct {
			IntegrationKey struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"integration_key,omitempty"`
			Description struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"description,omitempty"`
		} `json:"properties,omitempty"`
		Secret   []string `json:"secret,omitempty"`
		Required []string `json:"required,omitempty"`
	} `json:"configuration_schema,omitempty"`
}
type HangoutsChatDestination struct {
	Destination
	ConfgurationSchema struct {
		Type       string `json:"type,omitempty"`
		Properties struct {
			URL struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"url"`
			IconURL struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"icon_url"`
		} `json:"properties"`
		Secret   []string `json:"secret,omitempty"`
		Required []string `json:"required,omitempty"`
	} `json:"configuration_schema,omitempty"`
}

func ParseDestinationType(payload []byte) (destination interface{}, err error) {
	dst := &Destination{}
	err = json.Unmarshal(payload, dst)
	if err != nil {
		return nil, err
	}
	switch dst.Type {
	case "email":
		destination = &EmailDestination{}
	case "slack":
		destination = &SlackDestination{}
	case "webhook":
		destination = &WebhookDestination{}
	case "hipchat":
		destination = &HipChatDestination{}
	case "mattermost":
		destination = &MattermostDestination{}
	case "chatwork":
		destination = &ChatWorkDestination{}
	case "pagerduty":
		destination = &PagerDutyDestination{}
	case "hangouts_chat":
		destination = &HangoutsChatDestination{}
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
func (c *Client) GetDestinationTypes() (*[]DestinationTypes, error) {
	path := "/api/destinations/types"
	query := url.Values{}
	response, err := c.get(path, query)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	destinationTypes := []DestinationTypes{}
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
