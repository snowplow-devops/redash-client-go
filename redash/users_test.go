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
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("GET", "https://com.acme/api/users/1",
		httpmock.NewStringResponder(200, `{"id": 1, "name": "Existing User"}`))

	user, err := c.GetUser(1)
	if err != nil {
		panic(err.Error())
	}

	assert.Equal(1, user.ID)
	assert.Equal("Existing User", user.Name)
}

func TestCreateUser(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("POST", "https://com.acme/api/users",
		httpmock.NewStringResponder(200, `{"id": 2, "name": "New User", "email": "test@email.com"}`))

	userPayload := UserCreatePayload{
		Name:  "New User",
		Email: "test@email.com",
	}

	user, err := c.CreateUser(&userPayload)
	if err != nil {
		panic(err.Error())
	}

	assert.Equal(2, user.ID)
	assert.Equal("New User", user.Name)
}

func TestUpdateUser(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("POST", "https://com.acme/api/users/2",
		httpmock.NewStringResponder(200, `{"id": 2, "name": "New User Updated", "email": "test-update@email.com"}`))

	userPayload := UserUpdatePayload{
		Name:   "New User Updated",
		Email:  "test-update@email.com",
		Groups: []int{2, 3, 4},
	}

	user, err := c.UpdateUser(2, &userPayload)
	if err != nil {
		panic(err.Error())
	}

	assert.Equal(2, user.ID)
	assert.Equal("New User Updated", user.Name)
	assert.Equal("test-update@email.com", user.Email)
}

func TestGetUserByEmail(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("GET", "https://com.acme/api/users?q=tëst@email.com",
		httpmock.NewStringResponder(200, `{"count": 1, "page": 1, "page_size": 25, "results": [ {"id": 1, "name": "Existing User", "email": "tëst@email.com"} ]}`))

	httpmock.RegisterResponder("GET", "https://com.acme/api/users/1",
		httpmock.NewStringResponder(200, `{"id": 1, "name": "Existing User", "email": "tëst@email.com"}`))

	user, err := c.GetUserByEmail("tëst@email.com")
	if err != nil {
		panic(err.Error())
	}

	assert.Equal(1, user.ID)
	assert.Equal("tëst@email.com", user.Email)

	httpmock.RegisterResponder("GET", "https://com.acme/api/users?q=tëst-not-found@email.com",
		httpmock.NewStringResponder(200, `{"count": 0, "page": 1, "page_size": 25, "results": []}`))

	user, err = c.GetUserByEmail("tëst-not-found@email.com")
	if err != nil {
		assert.Equal("No user found with email address: tëst-not-found@email.com", err.Error())
	}

}
