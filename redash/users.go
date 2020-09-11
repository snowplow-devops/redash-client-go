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

// UserListOptions struct
type UserListOptions struct {
	Page       int
	PageSize   int
	SearchTerm string
}

//GetUsers gets an array of all Users available
func (c *Client) GetUsers(options *UserListOptions) (*UserList, error) {
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

//SearchUser finds a list of users matching the query
func (c *Client) SearchUser(term string) (*UserList, error) {
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
