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
	"time"
)

// DataSource struct
type DataSource struct {
	ID                 int                    `json:"id,omitempty"`
	Name               string                 `json:"name,omitempty"`
	ScheduledQueueName string                 `json:"scheduled_queue_name,omitempty"`
	QueueName          string                 `json:"queue_name,omitempty"`
	Options            map[string]interface{} `json:"options,omitempty"`
	Paused             int                    `json:"paused,omitempty"`
	PauseReason        string                 `json:"pause_reason,omitempty"`
	Type               string                 `json:"type,omitempty"`
	Syntax             string                 `json:"syntax,omitempty"`
	Groups             map[int]bool           `json:"groups,omitempty"`
}

// DataSourceType struct
type DataSourceType struct {
	Type                string `json:"type"`
	Name                string `json:"name,omitempty"`
	ConfigurationSchema struct {
		Secret     []string                               `json:"secret,omitempty"`
		Required   []string                               `json:"required,omitempty"`
		Type       string                                 `json:"type,omitempty"`
		Order      []string                               `json:"order,omitempty"`
		Properties map[string]DataSourceTypePropertyField `json:"properties,omitempty"`
	} `json:"configuration_schema,omitempty"`
}

// DataSourceTypePropertyField struct
type DataSourceTypePropertyField struct {
	Type    string
	Title   string
	Default interface{}
}

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

// User struct
type User struct {
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
}

// Group struct
type Group struct {
	CreatedAt   time.Time `json:"created_at,omitempty"`
	Permissions []string  `json:"permissions,omitempty"`
	Type        string    `json:"type,omitempty"`
	ID          int       `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
}

// GroupUser struct
type GroupUser struct {
	MemberID int `json:"user_id"`
}

// GroupDataSource struct
type GroupDataSource struct {
	DataSourceID int `json:"data_source_id"`
}
