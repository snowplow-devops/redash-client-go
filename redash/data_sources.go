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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// DataSource struct
type DataSource struct {
	ID                 int                    `json:"id,omitempty"`
	Name               string                 `json:"name,omitempty"`
	ScheduledQueueName string                 `json:"scheduled_queue_name,omitempty"`
	QueueName          string                 `json:"queue_name,omitempty"`
	Options            map[string]interface{} `json:"options,omitempty"`
	Paused             int                    `json:"paused,omitempty"`
	PauseReason        string                 `json:"pause_reason,omitempty"`
	Type               string                 `json:"type,omitempty"`
	Syntax             string                 `json:"syntax,omitempty"`
	Groups             map[int]bool           `json:"groups,omitempty"`
}

// DataSourceType struct
type DataSourceType struct {
	Type                string `json:"type"`
	Name                string `json:"name,omitempty"`
	ConfigurationSchema struct {
		Secret     []string                               `json:"secret,omitempty"`
		Required   []string                               `json:"required,omitempty"`
		Type       string                                 `json:"type,omitempty"`
		Order      []string                               `json:"order,omitempty"`
		Properties map[string]DataSourceTypePropertyField `json:"properties,omitempty"`
	} `json:"configuration_schema,omitempty"`
}

// DataSourceTypePropertyField struct
type DataSourceTypePropertyField struct {
	Type    string
	Title   string
	Default interface{}
}

//GetDataSources gets an array of all DataSources available
func (c *Client) GetDataSources() (*[]DataSource, error) {
	path := "/api/data_sources"
	query := url.Values{}
	response, err := c.get(path, query)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	dataSources := []DataSource{}
	err = json.Unmarshal(body, &dataSources)
	if err != nil {
		return nil, err
	}

	return &dataSources, nil
}

//GetDataSource gets a specific DataSource
func (c *Client) GetDataSource(id int) (*DataSource, error) {
	path := "/api/data_sources/" + strconv.Itoa(id)
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

	dataSource := DataSource{}

	err = json.Unmarshal(body, &dataSource)
	if err != nil {
		return nil, err
	}

	return &dataSource, nil
}

//GetDataSourceTypes gets all available types with configuration details
func (c *Client) GetDataSourceTypes() ([]DataSourceType, error) {
	path := "/api/data_sources/types"
	query := url.Values{}
	response, err := c.get(path, query)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	dataSourceTypes := []DataSourceType{}
	err = json.Unmarshal(body, &dataSourceTypes)
	if err != nil {
		return nil, err
	}

	return dataSourceTypes, nil
}

// SanitizeDataSourceOptions checks the validity of the options field in a
// DataSource.Option against Redash's API and cleans up when possible
func (c *Client) SanitizeDataSourceOptions(dataSource *DataSource) (*DataSource, error) {
	dataSourceTypes, err := c.GetDataSourceTypes()
	if err != nil {
		fmt.Println(err)
	}

	for _, dst := range dataSourceTypes {
		if dst.Type == dataSource.Type {

			for _, required := range dst.ConfigurationSchema.Required {
				// does dataSource.Options have everything in configuration_schema.required[] ?
				_, exists := dataSource.Options[required]
				if !exists {
					return nil, fmt.Errorf("Required field missing: " + required)
				}
			}

			for propName, propVal := range dataSource.Options {
				// does dataSource.Options only have what's in configuration_schema.properties[]?
				_, exists := dst.ConfigurationSchema.Properties[propName]
				if !exists {
					if c.IsStrict() {
						return nil, fmt.Errorf("Invalid field (%s) for type: %s", propName, dataSource.Type)
					}

					log.Warn(fmt.Sprintf("[WARN] Ignoring invalid field (%s) for type: %s", propName, dataSource.Type))
					delete((*dataSource).Options, propName)
					continue
				}

				// is the input value a valid data type?
				switch propVal.(type) {
				case int:
					if dst.ConfigurationSchema.Properties[propName].Type != "number" {
						return nil, fmt.Errorf("Invalid value type for %s", propName)
					}
				case string:
					if dst.ConfigurationSchema.Properties[propName].Type != "string" {
						return nil, fmt.Errorf("Invalid value type for %s", propName)
					}
				case bool:
					if dst.ConfigurationSchema.Properties[propName].Type != "boolean" {
						return nil, fmt.Errorf("Invalid value type for %s", propName)
					}
				default:
					return nil, fmt.Errorf("Invalid value type for %s", propName)
				}
			}
		}
	}

	return dataSource, nil
}

//CreateDataSource creates a new DataSource
func (c *Client) CreateDataSource(dataSourcePayload *DataSource) (*DataSource, error) {
	path := "/api/data_sources"

	dataSourcePayload, err := c.SanitizeDataSourceOptions(dataSourcePayload)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(dataSourcePayload)
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

	dataSource := DataSource{}

	err = json.Unmarshal(body, &dataSource)
	if err != nil {
		return nil, err
	}

	return &dataSource, nil
}

//UpdateDataSource Updates an existing DataSource
func (c *Client) UpdateDataSource(id int, dataSourcePayload *DataSource) (*DataSource, error) {
	path := "/api/data_sources/" + strconv.Itoa(id)

	dataSourcePayload, err := c.SanitizeDataSourceOptions(dataSourcePayload)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(dataSourcePayload)
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

	dataSource := DataSource{}

	err = json.Unmarshal(body, &dataSource)
	if err != nil {
		return nil, err
	}

	return &dataSource, nil
}

//DeleteDataSource deletes a specific DataSource
func (c *Client) DeleteDataSource(id int) error {
	path := "/api/data_sources/" + strconv.Itoa(id)

	query := url.Values{}
	_, err := c.delete(path, query)
	if err != nil {
		return err
	}

	return nil
}
