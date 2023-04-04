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

// Type of Destination
type Destination struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	Icon string `json:"icon,omitempty"`
}

// DestinationType stucture
type DestinationTypes struct {
	Destination
	ConfigurationSchema interface{} `json:"configuration_schema"`
}

// EmailDestination represents an email alert destination
type EmailDestination struct {
	Destination
	ConfgurationSchema struct {
		Type       string `json:"type,omitempty"`
		Properties struct {
			Addresses struct {
				Type string `json:"type,omitempty"`
			} `json:"addresses"`
			SubjectTemplate struct {
				Type    string `json:"type,omitempty"`
				Default string `json:"default,omitempty"`
				Title   string `json:"title,omitempty"`
			} `json:"subject_template"`
		} `json:"properties"`
		Required     []string `json:"required,omitempty"`
		ExtraOptions []string `json:"extra_options,omitempty"`
	} `json:"configuration_schema,omitempty"`
}

// SlackDestination represents a slack alert destination
type SlackDestination struct {
	Destination
	ConfgurationSchema struct {
		Type       string `json:"type"`
		Properties struct {
			URL struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"url,omitempty"`
			Username struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"username,omitempty"`
			IconEmoji struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"icon_emoji,omitempty"`
			IconURL struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"icon_url,omitempty"`
			Channel struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"channel,omitempty"`
		} `json:"properties,omitempty"`
		Secret []string `json:"secret,omitempty"`
	} `json:"configuration_schema,omitempty"`
}

// WebhookDestination represents a webhook alert destination
type WebhookDestination struct {
	Destination
	ConfgurationSchema struct {
		Type       string `json:"type"`
		Properties struct {
			URL struct {
				Type string `json:"type,omitempty"`
			} `json:"url,omitempty"`
			Username struct {
				Type string `json:"type,omitempty"`
			} `json:"username,omitempty"`
			Password struct {
				Type string `json:"type,omitempty"`
			} `json:"password,omitempty"`
		} `json:"properties,omitempty"`
		Required []string `json:"required,omitempty"`
		Secret   []string `json:"secret,omitempty"`
	} `json:"configuration_schema,omitempty"`
}

// HipChatDestination represents a Hip Chat alert destination
type HipChatDestination struct {
	Destination
	ConfgurationSchema struct {
		Type       string `json:"type"`
		Properties struct {
			URL struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"url,omitempty"`
		} `json:"properties,omitempty"`
		Secret   []string `json:"secret,omitempty"`
		Required []string `json:"required,omitempty"`
	} `json:"configuration_schema,omitempty"`
	Deprecated bool `json:"deprecated,omitempty"`
}

// MattermostDestination represents a Mattermost alert destination
type MattermostDestination struct {
	Destination
	ConfgurationSchema struct {
		Type       string `json:"type"`
		Properties struct {
			URL struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"url,omitempty"`
			Username struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"username,omitempty"`
			IconURL struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"icon_url,omitempty"`
			Channel struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"channel,omitempty"`
		} `json:"properties,omitempty"`
		Secret string `json:"secret,omitempty"`
	} `json:"configuration_schema,omitempty"`
}

// ChatWorkDestination represents a ChatWork alert destination
type ChatWorkDestination struct {
	Destination
	ConfgurationSchema struct {
		Type       string `json:"type"`
		Properties struct {
			APIToken struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"api_token,omitempty"`
			RoomID struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"room_id,omitempty"`
			MessageTemplate struct {
				Type    string `json:"type,omitempty"`
				Default string `json:"default,omitempty"`
				Title   string `json:"title,omitempty"`
			} `json:"message_template,omitempty"`
		} `json:"properties,omitempty"`
		Secret   []string `json:"secret,omitempty"`
		Required []string `json:"required,omitempty"`
	} `json:"configuration_schema,omitempty"`
}

// PagerDutyDestination represents PagerDuty alert destination
type PagerDutyDestination struct {
	Destination
	ConfgurationSchema struct {
		Type       string `json:"type"`
		Properties struct {
			IntegrationKey struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"integration_key,omitempty"`
			Description struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"description,omitempty"`
		} `json:"properties,omitempty"`
		Secret   []string `json:"secret,omitempty"`
		Required []string `json:"required,omitempty"`
	} `json:"configuration_schema,omitempty"`
}

// HangoutsChatDestination represents a Google Handgouts Chat alert destination
type HangoutsChatDestination struct {
	Destination
	ConfgurationSchema struct {
		Type       string `json:"type,omitempty"`
		Properties struct {
			URL struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"url"`
			IconURL struct {
				Type  string `json:"type,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"icon_url"`
		} `json:"properties"`
		Secret   []string `json:"secret,omitempty"`
		Required []string `json:"required,omitempty"`
	} `json:"configuration_schema,omitempty"`
}
