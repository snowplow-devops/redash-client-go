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
	"time"
)

// UserList struct
type UserList struct {
	Count    int `json:"count"`
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Results  []struct {
		AuthType            string    `json:"auth_type,omitempty"`
		IsDisabled          bool      `json:"is_disabled,omitempty"`
		UpdatedAt           time.Time `json:"updated_at,omitempty"`
		ProfileImageURL     string    `json:"profile_image_url,omitempty"`
		IsInvitationPending bool      `json:"is_invitation_pending,omitempty"`
		Groups              []struct {
			ID   int    `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"groups,omitempty"`
		ID              int         `json:"id,omitempty"`
		Name            string      `json:"name,omitempty"`
		CreatedAt       time.Time   `json:"created_at,omitempty"`
		DisabledAt      interface{} `json:"disabled_at,omitempty"`
		IsEmailVerified bool        `json:"is_email_verified,omitempty"`
		ActiveAt        time.Time   `json:"active_at,omitempty"`
		Email           string      `json:"email,omitempty"`
	} `json:"results,omitempty"`
}

// User representation
type User struct {
	AuthType            string      `json:"auth_type,omitempty"`
	IsDisabled          bool        `json:"is_disabled,omitempty"`
	UpdatedAt           time.Time   `json:"updated_at,omitempty"`
	ProfileImageURL     string      `json:"profile_image_url,omitempty"`
	IsInvitationPending bool        `json:"is_invitation_pending,omitempty"`
	Groups              []int       `json:"groups,omitempty"`
	ID                  int         `json:"id,omitempty"`
	Name                string      `json:"name,omitempty"`
	CreatedAt           time.Time   `json:"created_at,omitempty"`
	DisabledAt          interface{} `json:"disabled_at,omitempty"`
	IsEmailVerified     bool        `json:"is_email_verified,omitempty"`
	ActiveAt            time.Time   `json:"active_at,omitempty"`
	Email               string      `json:"email,omitempty"`
}

// UserCreatePayload struct for mutating users.
type UserCreatePayload struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserUpdatePayload struct for mutating users.
type UserUpdatePayload struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Groups []int  `json:"group_ids"`
}

//GetUsers returns a paginated list of users
func (c *Client) GetUsers() (*UserList, error) {
	path := "/api/users"

	query := url.Values{}
	response, err := c.get(path, query)

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

	user := User{}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateUser creates a new Redash user
func (c *Client) CreateUser(userCreatePayload *UserCreatePayload) (*User, error) {
	path := "/api/users"

	payload, err := json.Marshal(userCreatePayload)
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

	user := User{}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates an existing Redash user
func (c *Client) UpdateUser(id int, userUpdatePayload *UserUpdatePayload) (*User, error) {
	path := "/api/users/" + strconv.Itoa(id)

	payload, err := json.Marshal(userUpdatePayload)
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

	user := User{}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

//DisableUser disables an active user.
func (c *Client) DisableUser(id int) error {
	path := "/api/users/" + strconv.Itoa(id) + "/disable"

	query := url.Values{}
	response, err := c.post(path, "", query)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return nil
}

//SearchUsers finds a list of users matching a string (searches `name` and `email` fields)
func (c *Client) SearchUsers(term string) (*UserList, error) {
	path := "/api/users"

	query := url.Values{}
	query.Add("q", term)
	response, err := c.get(path, query)

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
