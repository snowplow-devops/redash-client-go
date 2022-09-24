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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	assert := assert.New(t)

	c, err := NewClient(&Config{RedashURI: "", APIKey: ""})
	assert.NotNil(err)
	assert.Nil(c)

	c, err = NewClient(&Config{RedashURI: "invalid.url", APIKey: "RanD0mStr1nG"})
	assert.NotNil(err)
	assert.Nil(c)

	c, err = NewClient(&Config{RedashURI: "s3://invalid.url/", APIKey: "RanD0mStr1nG"})
	assert.NotNil(err)
	assert.Nil(c)

	c, err = NewClient(&Config{RedashURI: "https://valid.url/", APIKey: ""})
	assert.NotNil(err)
	assert.Nil(c)

	c, err = NewClient(&Config{RedashURI: "https://valid.url/", APIKey: "RanD0mStr1nG"})
	assert.Nil(err)
	assert.NotNil(c)

	c, err = NewClient(&Config{RedashURI: "http://valid.url", APIKey: "RanD0mStr1nG"})
	assert.Nil(err)
	assert.NotNil(c)
}
