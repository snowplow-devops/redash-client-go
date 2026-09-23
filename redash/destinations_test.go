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
	"errors"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/stretchr/testify/assert"
)

func loadFixture(t *testing.T, filePath string) string {
	t.Helper()
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to load fixture %s: %s", filePath, err)
	}

	return string(content)
}

func TestGetDestinations(t *testing.T) {
	payload := loadFixture(t, "../testdata/destinations/destinations.json")

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

	seenIDs := map[int]bool{}
	for _, destination := range *resp {
		assert.False(seenIDs[destination.ID], "duplicate destination ID %d in fixture", destination.ID)
		seenIDs[destination.ID] = true
	}
}

func TestGetDestinationTypes(t *testing.T) {
	payload := loadFixture(t, "../testdata/destinations/types.json")

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
		loadFixture(t, "../testdata/destinations/request_email.json"),
		loadFixture(t, "../testdata/destinations/email.json"),
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
		destination, ok := resp.(*EmailDestination)
		if !ok {
			t.Errorf("Expected EmailDestination, but parsing produced %T", destination)
		}
		assert.Equal(1, destination.ID)
		assert.Equal(test.name, destination.Name)
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
		destination, ok := resp.(*EmailDestination)
		if !ok {
			t.Errorf("Expected EmailDestination, but parsing produced %T", destination)
		}
		assert.Equal(1, destination.ID)
		assert.Equal(test.name, destination.Name)
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
			t.Errorf("Expected EmailDestination, but parsing produced %T", destination)
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
	payload := loadFixture(t, "../testdata/destinations/email.json")

	result, err := ParseDestinationType([]byte(payload))
	if err != nil {
		panic(err.Error())
	}
	destination, ok := result.(*EmailDestination)
	if !ok {
		t.Errorf("Expected EmailDestination, but parsing produced %T", destination)
	}
}

func TestParseSlackDestination(t *testing.T) {
	assert := assert.New(t)
	payload := loadFixture(t, "../testdata/destinations/slack.json")

	result, err := ParseDestinationType([]byte(payload))
	if err != nil {
		panic(err.Error())
	}
	destination, ok := result.(*SlackDestination)
	if !ok {
		t.Errorf("Expected SlackDestination, but parsing produced %T", destination)
	}
	assert.Equal("#redash", destination.Options.Channel, "SlackDestination.Options should be populated like every other destination type")
}

func TestParseChatWorkDestination(t *testing.T) {
	payload := loadFixture(t, "../testdata/destinations/chatwork.json")

	result, err := ParseDestinationType([]byte(payload))
	if err != nil {
		panic(err.Error())
	}
	destination, ok := result.(*ChatWorkDestination)
	if !ok {
		t.Errorf("Expected ChatWorkDestination, but parsing produced %T", destination)
	}
}

func TestParseHangoutsChatDestination(t *testing.T) {
	payload := loadFixture(t, "../testdata/destinations/hangoutschat.json")

	result, err := ParseDestinationType([]byte(payload))
	if err != nil {
		panic(err.Error())
	}
	destination, ok := result.(*HangoutsChatDestination)
	if !ok {
		t.Errorf("Expected HangoutsChatDestination, but parsing produced %T", destination)
	}
}

func TestParseMattermostDestination(t *testing.T) {
	payload := loadFixture(t, "../testdata/destinations/mattermost.json")

	result, err := ParseDestinationType([]byte(payload))
	if err != nil {
		panic(err.Error())
	}
	destination, ok := result.(*MattermostDestination)
	if !ok {
		t.Errorf("Expected MattermostDestination, but parsing produced %T", destination)
	}
}

func TestParseWebhookDestination(t *testing.T) {
	payload := loadFixture(t, "../testdata/destinations/webhook.json")

	result, err := ParseDestinationType([]byte(payload))
	if err != nil {
		panic(err.Error())
	}
	destination, ok := result.(*WebhookDestination)
	if !ok {
		t.Errorf("Expected WebhookDestination, but parsing produced %T", destination)
	}
}

func TestParsePagerDutyDestination(t *testing.T) {
	payload := loadFixture(t, "../testdata/destinations/pagerduty.json")

	result, err := ParseDestinationType([]byte(payload))
	if err != nil {
		panic(err.Error())
	}
	destination, ok := result.(*PagerDutyDestination)
	if !ok {
		t.Errorf("Expected PagerDutyDestination, but parsing produced %T", destination)
	}
}

func TestParseDestinationType_InvalidJSON(t *testing.T) {
	assert := assert.New(t)

	result, err := ParseDestinationType([]byte("not json"))
	assert.Nil(result)
	assert.Error(err)
}

func TestParseDestinationType_MissingOptions(t *testing.T) {
	assert := assert.New(t)

	result, err := ParseDestinationType([]byte(`{"id": 1, "name": "x", "type": "slack"}`))
	assert.Nil(result)
	assert.EqualError(err, "empty options field")
}

func TestParseDestinationType_UnknownType(t *testing.T) {
	assert := assert.New(t)

	result, err := ParseDestinationType([]byte(`{"id": 1, "name": "x", "type": "hipchat", "options": {"room": "general"}}`))
	assert.True(errors.Is(err, ErrUnsupportedDestinationType), "expected error to wrap ErrUnsupportedDestinationType, got %v", err)

	destination, ok := result.(*Destination)
	if !ok {
		t.Fatalf("Expected generic *Destination fallback, but parsing produced %T", result)
	}
	assert.Equal(1, destination.ID)
	assert.Equal("x", destination.Name)
	assert.Equal("hipchat", destination.Type)
	assert.Equal(map[string]interface{}{"room": "general"}, destination.Options)
}

func TestParseDestinationType_TypeMismatch(t *testing.T) {
	assert := assert.New(t)

	result, err := ParseDestinationType([]byte(`{"id": 1, "name": "x", "type": "slack", "options": {"url": 12345}}`))
	assert.Nil(result)
	assert.Error(err)
}

func TestGetDestinations_Error(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("GET", "https://com.acme/api/destinations", httpmock.NewStringResponder(500, "internal error"))

	resp, err := c.GetDestinations()
	assert.Nil(resp)
	assert.Error(err)
}

func TestGetDestinations_InvalidJSON(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("GET", "https://com.acme/api/destinations", httpmock.NewStringResponder(200, "not json"))

	resp, err := c.GetDestinations()
	assert.Nil(resp)
	assert.Error(err)
}

func TestGetDestination_Error(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("GET", "https://com.acme/api/destinations/1", httpmock.NewStringResponder(404, "not found"))

	resp, err := c.GetDestination(1)
	assert.Nil(resp)
	assert.Error(err)
}

func TestGetDestination_ParseError(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("GET", "https://com.acme/api/destinations/1",
		httpmock.NewStringResponder(200, `{"id": 1, "name": "x", "type": "slack"}`))

	resp, err := c.GetDestination(1)
	assert.Nil(resp)
	assert.EqualError(err, "empty options field")
}

func TestGetDestination_UnsupportedType(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("GET", "https://com.acme/api/destinations/1",
		httpmock.NewStringResponder(200, `{"id": 1, "name": "x", "type": "hipchat", "options": {"room": "general"}}`))

	resp, err := c.GetDestination(1)
	assert.True(errors.Is(err, ErrUnsupportedDestinationType), "expected error to wrap ErrUnsupportedDestinationType, got %v", err)

	destination, ok := resp.(*Destination)
	if !ok {
		t.Fatalf("Expected generic *Destination fallback, but got %T", resp)
	}
	assert.Equal("hipchat", destination.Type)
}

func TestGetDestinationTypes_Error(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("GET", "https://com.acme/api/destinations/types", httpmock.NewStringResponder(500, "internal error"))

	resp, err := c.GetDestinationTypes()
	assert.Nil(resp)
	assert.Error(err)
}

func TestGetDestinationTypes_InvalidJSON(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("GET", "https://com.acme/api/destinations/types", httpmock.NewStringResponder(200, "not json"))

	resp, err := c.GetDestinationTypes()
	assert.Nil(resp)
	assert.Error(err)
}

func TestCreateDestination_InvalidPayload(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	resp, err := c.CreateDestination([]byte(`{"name": "x", "type": "slack"}`))
	assert.Nil(resp)
	assert.EqualError(err, "empty options field")
}

func TestCreateDestination_Error(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("POST", "https://com.acme/api/destinations", httpmock.NewStringResponder(500, "internal error"))

	resp, err := c.CreateDestination([]byte(`{"name": "x", "type": "slack", "options": {}}`))
	assert.Nil(resp)
	assert.Error(err)
}

func TestCreateDestination_ParseResponseError(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("POST", "https://com.acme/api/destinations", httpmock.NewStringResponder(200, "not json"))

	resp, err := c.CreateDestination([]byte(`{"name": "x", "type": "slack", "options": {}}`))
	assert.Nil(resp)
	assert.Error(err)
}

func TestCreateDestination_UnsupportedType(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("POST", "https://com.acme/api/destinations",
		httpmock.NewStringResponder(200, `{"id": 1, "name": "x", "type": "hipchat", "options": {"room": "general"}}`))

	resp, err := c.CreateDestination([]byte(`{"name": "x", "type": "hipchat", "options": {"room": "general"}}`))
	assert.True(errors.Is(err, ErrUnsupportedDestinationType), "an unrecognized type should not block creation; got %v", err)

	destination, ok := resp.(*Destination)
	if !ok {
		t.Fatalf("Expected generic *Destination fallback, but got %T", resp)
	}
	assert.Equal("hipchat", destination.Type)
}

func TestUpdateDestination_InvalidPayload(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	resp, err := c.UpdateDestination(1, []byte(`{"name": "x", "type": "slack"}`))
	assert.Nil(resp)
	assert.EqualError(err, "empty options field")
}

func TestUpdateDestination_Error(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("POST", "https://com.acme/api/destinations/1", httpmock.NewStringResponder(500, "internal error"))

	resp, err := c.UpdateDestination(1, []byte(`{"name": "x", "type": "slack", "options": {}}`))
	assert.Nil(resp)
	assert.Error(err)
}

func TestUpdateDestination_ParseResponseError(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("POST", "https://com.acme/api/destinations/1", httpmock.NewStringResponder(200, "not json"))

	resp, err := c.UpdateDestination(1, []byte(`{"name": "x", "type": "slack", "options": {}}`))
	assert.Nil(resp)
	assert.Error(err)
}

func TestUpdateDestination_UnsupportedType(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("POST", "https://com.acme/api/destinations/1",
		httpmock.NewStringResponder(200, `{"id": 1, "name": "x", "type": "hipchat", "options": {"room": "general"}}`))

	resp, err := c.UpdateDestination(1, []byte(`{"name": "x", "type": "hipchat", "options": {"room": "general"}}`))
	assert.True(errors.Is(err, ErrUnsupportedDestinationType), "an unrecognized type should not block updating; got %v", err)

	destination, ok := resp.(*Destination)
	if !ok {
		t.Fatalf("Expected generic *Destination fallback, but got %T", resp)
	}
	assert.Equal("hipchat", destination.Type)
}

func TestDeleteDestination_Error(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("DELETE", "https://com.acme/api/destinations/1", httpmock.NewStringResponder(500, "internal error"))

	err := c.DeleteDestination(1)
	assert.Error(err)
}
