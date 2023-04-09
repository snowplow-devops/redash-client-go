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
	"testing"

	"github.com/stretchr/testify/assert"
)

func loadFixture(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

func TestGetDestinations(t *testing.T) {
	payload := loadFixture("../testdata/destinations/destinations.json")

	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("GET", "https://com.acme/api/destinations", httpmock.NewStringResponder(200, payload))

	resp, err := c.GetDestinations()
	if err != nil {
		panic(err.Error())
	}

	assert.Len(*resp, 8)
}

func TestGetDestinationTypes(t *testing.T) {
	payload := loadFixture("../testdata/destinations/types.json")

	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("GET", "https://com.acme/api/destinations/types", httpmock.NewStringResponder(200, payload))

	resp, err := c.GetDestinationTypes()
	if err != nil {
		panic(err.Error())
	}

	assert.Len(*resp, 8)
}

func TestDestination(t *testing.T) {
	test := struct {
		name            string
		requestPayload  string
		responsePayload string
	}{
		"Test Email",
		loadFixture("../testdata/destinations/request_email.json"),
		loadFixture("../testdata/destinations/email.json"),
	}
	t.Run("Create", func(t *testing.T) {
		assert := assert.New(t)
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

		httpmock.RegisterResponder("POST", "https://com.acme/api/destinations", httpmock.NewStringResponder(200, test.responsePayload))

		resp, err := c.CreateDestination([]byte(test.requestPayload))
		if err != nil {
			panic(err.Error())
		}
		assert.Equal(1, resp.ID)
		assert.Equal(test.name, resp.Name)
	})
	t.Run("Update", func(t *testing.T) {
		assert := assert.New(t)
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

		httpmock.RegisterResponder("POST", "https://com.acme/api/destinations/1", httpmock.NewStringResponder(200, test.responsePayload))

		resp, err := c.UpdateDestination(1, []byte(test.requestPayload))
		if err != nil {
			panic(err.Error())
		}
		log.Errorf("Payload: %T", resp)
		assert.Equal(1, resp.ID)
		assert.Equal(test.name, resp.Name)
	})
	t.Run("Get", func(t *testing.T) {
		assert := assert.New(t)
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

		httpmock.RegisterResponder("GET", "https://com.acme/api/destinations/1", httpmock.NewStringResponder(200, test.responsePayload))

		resp, err := c.GetDestination(1)
		if err != nil {
			panic(err.Error())
		}
		destination, ok := resp.(*EmailDestination)
		if !ok {
			t.Errorf("Expected ProjectSystemHookEvent, but parsing produced %T", destination)
		}
		assert.Equal(1, destination.ID)
		assert.Equal(test.name, destination.Name)
	})
	t.Run("Delete", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

		httpmock.RegisterResponder("DELETE", "https://com.acme/api/destinations/1", httpmock.NewStringResponder(204, "NO CONTENT"))

		err := c.DeleteDestination(1)
		if err != nil {
			panic(err.Error())
		}
	})
}

func TestGetEmailDestination(t *testing.T) {
	payload := loadFixture("../testdata/destinations/email.json")

	result, err := ParseDestinationType([]byte(payload))
	if err != nil {
		panic(err.Error())
	}
	destination, ok := result.(*EmailDestination)
	if !ok {
		t.Errorf("Expected ProjectSystemHookEvent, but parsing produced %T", destination)
	}
}

func TestParseSlackDestination(t *testing.T) {
	payload := loadFixture("../testdata/destinations/slack.json")

	result, err := ParseDestinationType([]byte(payload))
	if err != nil {
		panic(err.Error())
	}
	destination, ok := result.(*SlackDestination)
	if !ok {
		t.Errorf("Expected ProjectSystemHookEvent, but parsing produced %T", destination)
	}
}

func TestParseChatWorkDestination(t *testing.T) {
	payload := loadFixture("../testdata/destinations/chatwork.json")

	result, err := ParseDestinationType([]byte(payload))
	if err != nil {
		panic(err.Error())
	}
	destination, ok := result.(*ChatWorkDestination)
	if !ok {
		t.Errorf("Expected ProjectSystemHookEvent, but parsing produced %T", destination)
	}
}

func TestParseHangoutsChatDestination(t *testing.T) {
	payload := loadFixture("../testdata/destinations/hangoutschat.json")

	result, err := ParseDestinationType([]byte(payload))
	if err != nil {
		panic(err.Error())
	}
	destination, ok := result.(*HangoutsChatDestination)
	if !ok {
		t.Errorf("Expected ProjectSystemHookEvent, but parsing produced %T", destination)
	}
}

func TestParseMattermostDestination(t *testing.T) {
	payload := loadFixture("../testdata/destinations/mattermost.json")

	result, err := ParseDestinationType([]byte(payload))
	if err != nil {
		panic(err.Error())
	}
	destination, ok := result.(*MattermostDestination)
	if !ok {
		t.Errorf("Expected ProjectSystemHookEvent, but parsing produced %T", destination)
	}
}

func TestParseWebhookDestination(t *testing.T) {
	payload := loadFixture("../testdata/destinations/webhook.json")

	result, err := ParseDestinationType([]byte(payload))
	if err != nil {
		panic(err.Error())
	}
	destination, ok := result.(*WebhookDestination)
	if !ok {
		t.Errorf("Expected ProjectSystemHookEvent, but parsing produced %T", destination)
	}
}

func TestParsePagerDutyDestination(t *testing.T) {
	payload := loadFixture("../testdata/destinations/pagerduty.json")

	result, err := ParseDestinationType([]byte(payload))
	if err != nil {
		panic(err.Error())
	}
	destination, ok := result.(*PagerDutyDestination)
	if !ok {
		t.Errorf("Expected ProjectSystemHookEvent, but parsing produced %T", destination)
	}
}
