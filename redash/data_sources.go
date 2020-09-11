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
	"strconv"
)

//GetDataSources gets an array of all DataSources available
func (c *Client) GetDataSources() ([]DataSource, error) {
	path := "/api/data_sources"

	response, err := c.get(path)

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

	return dataSources, nil
}

//GetDataSource gets a specific DataSource
func (c *Client) GetDataSource(id int) (*DataSource, error) {
	path := "/api/data_sources/" + strconv.Itoa(id)

	response, err := c.get(path)
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

	response, err := c.get(path)

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

// ValidateDataSourceOptions checks the validity of the options field in a DataSource.Option against Redash's API
func (c *Client) ValidateDataSourceOptions(dataSource DataSource) error {
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
					return fmt.Errorf("Required field missing: " + required)
				}
			}

			for propName, propVal := range dataSource.Options {
				// is the input value a valid data type?
				switch propVal.(type) {
				case int:
					if dst.ConfigurationSchema.Properties[propName].Type != "number" {
						return fmt.Errorf("Invalid value for %s", propName)
					}
				case string:
					if dst.ConfigurationSchema.Properties[propName].Type != "string" {
						return fmt.Errorf("Invalid value for %s", propName)
					}
				case bool:
					if dst.ConfigurationSchema.Properties[propName].Type != "boolean" {
						return fmt.Errorf("Invalid value for %s", propName)
					}
				default:
					return fmt.Errorf("Invalid value for %s", propName)
				}

				// does dataSource.Options only have what's in configuration_schema.properties[]?
				_, exists := dst.ConfigurationSchema.Properties[propName]
				if !exists {
					return fmt.Errorf("Invalid field for type: " + propName)
				}
			}
		}
	}

	return nil
}

//CreateDataSource creaes a new DataSource
func (c *Client) CreateDataSource(dataSource DataSource) (*DataSource, error) {
	path := "/api/data_sources"

	err := c.ValidateDataSourceOptions(dataSource)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(dataSource)
	if err != nil {
		return nil, err
	}

	response, err := c.post(path, string(payload))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &dataSource)
	if err != nil {
		return nil, err
	}

	return &dataSource, nil
}

//UpdateDataSource Updates an existing DataSource
func (c *Client) UpdateDataSource(id int, dataSource DataSource) (*DataSource, error) {
	path := "/api/data_sources/" + strconv.Itoa(id)

	err := c.ValidateDataSourceOptions(dataSource)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(dataSource)
	if err != nil {
		return nil, err
	}

	response, err := c.post(path, string(payload))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &dataSource)
	if err != nil {
		return nil, err
	}

	return &dataSource, nil
}

//DeleteDataSource deletes a specific DataSource
func (c *Client) DeleteDataSource(id int) error {
	path := "/api/data_sources/" + strconv.Itoa(id)

	_, err := c.delete(path)
	if err != nil {
		return err
	}

	return nil
}
