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

// Destination represents Base structure of alert destination
type Destination struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	Icon string `json:"icon,omitempty"`
}

// DestinationCommon represents common alert destination stucture
type DestinationCommon struct {
	Destination
	ConfigurationSchema interface{} `json:"configuration_schema"`
}

// EmailDestination represents an email alert destination
type EmailDestination struct {
	Destination
	Options struct {
		Addresses       string `json:"addresses"`
		SubjectTemplate string `json:"subject_template,omitempty"`
	} `json:"options"`
}

// SlackDestination represents a slack alert destination
type SlackDestination struct {
	Destination
	Options struct {
		URL       string `json:"url,omitempty"`
		Username  string `json:"username,omitempty"`
		IconEmoji string `json:"icon_emoji,omitempty"`
		IconURL   string `json:"icon_url,omitempty"`
		Channel   string `json:"channel,omitempty"`
	} `json:"options,omitempty"`
}

// WebhookDestination represents a webhook alert destination
type WebhookDestination struct {
	Destination
	Options struct {
		URL      string `json:"url"`
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	} `json:"options,omitempty"`
}

// MattermostDestination represents a Mattermost alert destination
type MattermostDestination struct {
	Destination
	Options struct {
		URL      string `json:"url,omitempty"`
		Username string `json:"username,omitempty"`
		IconURL  string `json:"icon_url,omitempty"`
		Channel  string `json:"channel,omitempty"`
	} `json:"options,omitempty"`
}

// ChatWorkDestination represents a ChatWork alert destination
type ChatWorkDestination struct {
	Destination
	Options struct {
		APIToken        string `json:"api_token"`
		RoomID          string `json:"room_id"`
		MessageTemplate string `json:"message_template"`
	} `json:"options"`
}

// PagerDutyDestination represents PagerDuty alert destination
type PagerDutyDestination struct {
	Destination
	Options struct {
		IntegrationKey string `json:"integration_key"`
		Description    string `json:"description,omitempty"`
	} `json:"options"`
}

// HangoutsChatDestination represents a Google Handgouts Chat alert destination
type HangoutsChatDestination struct {
	Destination
	Options struct {
		URL     string `json:"url"`
		IconURL string `json:"icon_url,omitempty"`
	} `json:"options"`
}
