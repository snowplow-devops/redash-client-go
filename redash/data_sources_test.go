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
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

// captureStdout redirects os.Stdout for the duration of fn and returns
// whatever was written to it. Used to prove that library calls only
// communicate errors through their return values, since a Terraform
// plugin's stdout is reserved for the go-plugin RPC handshake - anything
// else written there gets logged by Terraform core as "unexpected data"
// and can break the provider.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %s", err)
	}

	original := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = original }()

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("failed to close pipe writer: %s", err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("failed to read captured stdout: %s", err)
	}

	return buf.String()
}

// A misconfigured RedashURI (wrong host, missing /api prefix, a proxy in
// front of Redash, etc.) can result in GetDataSourceTypes() receiving an
// HTML error page instead of JSON. SanitizeDataSourceOptions must surface
// that failure rather than silently continuing with an empty type list,
// and must never write it directly to stdout.
func TestSanitizeDataSourceOptions_GetDataSourceTypesError(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	httpmock.RegisterResponder("GET", "https://com.acme/api/data_sources/types",
		httpmock.NewStringResponder(200, `<html><body>404 Not Found</body></html>`))

	dataSource := &DataSource{
		Type:    "pg",
		Options: map[string]interface{}{"host": "localhost"},
	}

	var result *DataSource
	var err error

	stdout := captureStdout(t, func() {
		result, err = c.SanitizeDataSourceOptions(dataSource)
	})

	assert.Empty(stdout, "SanitizeDataSourceOptions must not write directly to stdout: it corrupts the Terraform plugin RPC stream")
	assert.Error(err, "a failure to fetch data source types must be surfaced, not swallowed")
	assert.Nil(result)
}
