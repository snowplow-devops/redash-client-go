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

//GetUsers returns a paginated list of users
func (c *Client) GetUsers() (*UserList, error) {
	path := "/api/users"

	response, err := c.get(path)

	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(response.Body)

	users := UserList{}
	err = json.Unmarshal(body, &users)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return &users, nil
}

//GetUser gets a specific User
func (c *Client) GetUser(id int) (*User, error) {
	path := "/api/users/" + strconv.Itoa(id)

	response, err := c.get(path)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	user := User{}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

//SearchUsers finds a list of users matching a string (searches `name` and `email` fields)
func (c *Client) SearchUsers(term string) (*UserList, error) {
	path := "/api/users?q=" + term

	response, err := c.get(path)

	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(response.Body)

	users := UserList{}
	err = json.Unmarshal(body, &users)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return &users, nil
}

// GetUserByEmail returns a single  user from their email address
func (c *Client) GetUserByEmail(email string) (*User, error) {

	results, err := c.SearchUsers(email)
	if err != nil {
		return nil, err
	}

	for _, result := range results.Results {
		if result.Email != "" && result.Email == email {
			return c.GetUser(result.ID)
		}
	}

	return nil, fmt.Errorf("No user found with email address: %s", email)
}
