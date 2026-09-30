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
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
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

// A non-2xx response's body is often the only clue to what actually went
// wrong (an HTML error page from a misrouted RedashURI, a JSON error from
// Redash itself, etc.), so doRequest must include it in the returned error
// instead of just the bare status code.
func TestDoRequestUsesCustomHTTPClient(t *testing.T) {
	assert := assert.New(t)

	var extraHeader, authorization string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		extraHeader = r.Header.Get("X-Request-Id")
		authorization = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	httpClient := srv.Client()
	base := httpClient.Transport
	httpClient.Transport = roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		r.Header.Set("X-Request-Id", "abc123")
		return base.RoundTrip(r)
	})

	c, err := NewClient(&Config{
		RedashURI:  srv.URL,
		APIKey:     "ApIkEyApIkEyApIkEyApIkEyApIkEy",
		HTTPClient: httpClient,
	})
	assert.Nil(err)

	_, err = c.get("/api/groups", url.Values{})
	assert.Nil(err)
	assert.Equal("abc123", extraHeader)
	assert.Equal("Key ApIkEyApIkEyApIkEyApIkEyApIkEy", authorization)
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestDoRequestNon2xxIncludesResponseBody(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("GET", "https://com.acme/api/groups",
		httpmock.NewStringResponder(404, `<html><body>404 Not Found</body></html>`))

	_, err := c.get("/api/groups", url.Values{})

	assert.Error(err)
	assert.Contains(err.Error(), "404")
	assert.Contains(err.Error(), "404 Not Found")
}

func TestDoRequestNon2xxWithEmptyBody(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("GET", "https://com.acme/api/groups",
		httpmock.NewStringResponder(404, ""))

	_, err := c.get("/api/groups", url.Values{})

	assert.EqualError(err, "HTTP Response: 404")
}

func TestDoRequestNon2xxTruncatesLongBody(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	longBody := strings.Repeat("x", maxErrorBodySnippet+100)
	httpmock.RegisterResponder("GET", "https://com.acme/api/groups",
		httpmock.NewStringResponder(500, longBody))

	_, err := c.get("/api/groups", url.Values{})

	assert.Error(err)
	assert.LessOrEqual(len(err.Error()), len("HTTP Response: 500: ")+maxErrorBodySnippet+len("..."))
	assert.Contains(err.Error(), "...")
}
