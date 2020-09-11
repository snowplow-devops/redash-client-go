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
	"io/ioutil"
	"strconv"
)

// GetGroups returns a list of Redash groups
func (c *Client) GetGroups() ([]Group, error) {
	path := "/api/groups"

	response, err := c.get(path)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	groups := []Group{}
	err = json.Unmarshal(body, &groups)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

// GetGroup returns an individual Redash group
func (c *Client) GetGroup(id int) (*Group, error) {
	path := "/api/groups/" + strconv.Itoa(id)

	response, err := c.get(path)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	group := Group{}

	err = json.Unmarshal(body, &group)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

// CreateGroup creates a new Redash group
func (c *Client) CreateGroup(group Group) (*Group, error) {
	path := "/api/groups"

	payload, err := json.Marshal(group)
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

	err = json.Unmarshal(body, &group)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

// UpdateGroup updates an existing Redash group
func (c *Client) UpdateGroup(id int, group Group) (*Group, error) {
	path := "/api/groups/" + strconv.Itoa(id)

	payload, err := json.Marshal(group)
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

	err = json.Unmarshal(body, &group)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

// DeleteGroup deletes a Redash group
func (c *Client) DeleteGroup(id int) error {
	path := "/api/groups/" + strconv.Itoa(id)

	_, err := c.delete(path)
	if err != nil {
		return err
	}

	return nil
}

// GroupAddUser adds a user to a Redash group
func (c *Client) GroupAddUser(groupID int, userID int) error {
	path := "/api/groups/" + strconv.Itoa(groupID) + "/members"

	user := GroupUser{userID}
	payload, err := json.Marshal(user)
	if err != nil {
		return err
	}

	response, err := c.post(path, string(payload))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

// GroupRemoveUser removes a user from a Redash group
func (c *Client) GroupRemoveUser(groupID int, userID int) error {
	path := "/api/groups/" + strconv.Itoa(groupID) + "/members/" + strconv.Itoa(userID)

	response, err := c.delete(path)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

// GroupAddDataSource adds a Data Source to a Redash group
func (c *Client) GroupAddDataSource(groupID int, dataSourceID int) error {
	path := "/api/groups/" + strconv.Itoa(groupID) + "/data_sources"

	dataSource := GroupDataSource{dataSourceID}
	payload, err := json.Marshal(dataSource)
	if err != nil {
		return err
	}

	response, err := c.post(path, string(payload))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

// GroupRemoveDataSource removes a Data Source from a Redash group
func (c *Client) GroupRemoveDataSource(groupID int, dataSourceID int) error {
	path := "/api/groups/" + strconv.Itoa(groupID) + "/data_sources/" + strconv.Itoa(dataSourceID)

	response, err := c.delete(path)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}
