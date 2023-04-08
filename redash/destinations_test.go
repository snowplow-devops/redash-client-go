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
	"github.com/jarcoal/httpmock"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func isType(a, b interface{}) bool {
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}

func loadFixture(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

func TestGetEmailDestination(t *testing.T) {
	tests := struct {
		name    string
		payload string
	}{
		"Test Email",
		loadFixture("../testdata/destinations/email.json"),
	}

	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("GET", "https://com.acme/api/destinations/1", httpmock.NewStringResponder(200, tests.payload))

	resp, err := c.GetDestination(1)
	if err != nil {
		panic(err.Error())
	}
	destination := resp.(*EmailDestination)
	assert.Equal(1, destination.ID)
	assert.Equal(tests.name, destination.Name)
}

func TestGetSlackDestination(t *testing.T) {
	tests := struct {
		name    string
		payload string
	}{
		"Test Slack",
		loadFixture("../testdata/destinations/slack.json"),
	}

	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("GET", "https://com.acme/api/destinations/1", httpmock.NewStringResponder(200, tests.payload))

	resp, err := c.GetDestination(1)
	if err != nil {
		panic(err.Error())
	}
	destination := resp.(*SlackDestination)
	assert.Equal(1, destination.ID)
	assert.Equal(tests.name, destination.Name)
}

func TestGetChatWorkDestination(t *testing.T) {
	tests := struct {
		name    string
		payload string
	}{
		"Test ChatWork",
		loadFixture("../testdata/destinations/chatwork.json"),
	}

	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("GET", "https://com.acme/api/destinations/1", httpmock.NewStringResponder(200, tests.payload))

	resp, err := c.GetDestination(1)
	if err != nil {
		panic(err.Error())
	}
	destination := resp.(*ChatWorkDestination)
	assert.Equal(1, destination.ID)
	assert.Equal(tests.name, destination.Name)
}

func TestGetHangoutsChatDestination(t *testing.T) {
	tests := struct {
		name    string
		payload string
	}{
		"Test Google Hangouts Chat",
		loadFixture("../testdata/destinations/hangoutschat.json"),
	}

	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("GET", "https://com.acme/api/destinations/1", httpmock.NewStringResponder(200, tests.payload))

	resp, err := c.GetDestination(1)
	if err != nil {
		panic(err.Error())
	}
	destination := resp.(*HangoutsChatDestination)
	assert.Equal(1, destination.ID)
	assert.Equal(tests.name, destination.Name)
}

func TestGetHipChatDestination(t *testing.T) {
	tests := struct {
		name    string
		payload string
	}{
		"Test HipChat",
		loadFixture("../testdata/destinations/hipchat.json"),
	}

	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("GET", "https://com.acme/api/destinations/1", httpmock.NewStringResponder(200, tests.payload))

	resp, err := c.GetDestination(1)
	if err != nil {
		panic(err.Error())
	}
	destination := resp.(*HipChatDestination)
	assert.Equal(1, destination.ID)
	assert.Equal(tests.name, destination.Name)
}

func TestGetMattermostDestination(t *testing.T) {
	tests := struct {
		name    string
		payload string
	}{
		"Test Mattermost",
		loadFixture("../testdata/destinations/mattermost.json"),
	}

	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("GET", "https://com.acme/api/destinations/1", httpmock.NewStringResponder(200, tests.payload))

	resp, err := c.GetDestination(1)
	if err != nil {
		panic(err.Error())
	}
	destination := resp.(*MattermostDestination)
	assert.Equal(1, destination.ID)
	assert.Equal(tests.name, destination.Name)
}

func TestGetWebhookDestination(t *testing.T) {
	tests := struct {
		name    string
		payload string
	}{
		"Test Webhook",
		loadFixture("../testdata/destinations/webhook.json"),
	}

	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("GET", "https://com.acme/api/destinations/1", httpmock.NewStringResponder(200, tests.payload))

	resp, err := c.GetDestination(1)
	if err != nil {
		panic(err.Error())
	}
	destination := resp.(*WebhookDestination)
	assert.Equal(1, destination.ID)
	assert.Equal(tests.name, destination.Name)
}

func TestGetPagerDutyDestination(t *testing.T) {
	tests := struct {
		name    string
		payload string
	}{
		"Test PagerDuty",
		loadFixture("../testdata/destinations/pagerduty.json"),
	}

	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("GET", "https://com.acme/api/destinations/1", httpmock.NewStringResponder(200, tests.payload))

	resp, err := c.GetDestination(1)
	if err != nil {
		panic(err.Error())
	}
	destination := resp.(*PagerDutyDestination)
	assert.Equal(1, destination.ID)
	assert.Equal(tests.name, destination.Name)
}