package redash

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetQueries(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	body, err := ioutil.ReadFile("testdata/get-queries.json")
	if err != nil {
		panic(err.Error())
	}
	httpmock.RegisterResponder("GET", "https://com.acme/api/queries",
		httpmock.NewStringResponder(200, string(body)))

	queries, err := c.GetQueries()
	assert.Nil(err)

	assert.Equal(3, queries.Count)
	assert.Equal(1, queries.Page)
	assert.Equal(10, queries.PageSize)
	assert.Equal(3, len(queries.Results))
}

func TestGetQuery(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, _ := NewClient(&Config{RedashURI: "https://com.acme/", APIKey: "ApIkEyApIkEyApIkEyApIkEyApIkEy"})

	body, err := ioutil.ReadFile("testdata/get-query.json")
	if err != nil {
		panic(err.Error())
	}
	httpmock.RegisterResponder("GET", "https://com.acme/api/queries/1",
		httpmock.NewStringResponder(200, string(body)))

	query, err := c.GetQuery(1)
	assert.Nil(err)

	assert.Equal(1, query.ID)
	assert.Equal("Daily Active Users", query.Name)
	assert.Equal("Service X DAU", query.Description)
	assert.Equal("SELECT 1 + 1;", query.Query)
	assert.Equal("ec2fda0cc5a54b38f81744fcad43ce5a", query.QueryHash)
	assert.Equal(1, query.Version)
	assert.Equal(false, query.IsArchived)
	assert.Equal(false, query.IsDraft)
	assert.Equal(true, query.IsSafe)
	assert.Equal(false, query.IsFavorite)
	assert.Equal(false, query.CanEdit)
	assert.Equal(2, query.DataSourceID)
	expectedUpdateAt, _ := time.Parse(time.RFC3339, "2021-11-07T22:22:34.929Z")
	assert.Equal(expectedUpdateAt, query.UpdatedAt)
	expectedCreatedAt, _ := time.Parse(time.RFC3339, "2021-08-13T23:29:12.743Z")
	assert.Equal(expectedCreatedAt, query.CreatedAt)

	assert.Equal(1, query.User.ID)
	assert.Equal("Admin", query.User.Name)
	assert.Equal("admin@example.com", query.User.Email)

	assert.Equal(2, query.LastModifiedBy.ID)
	assert.Equal("Developer", query.LastModifiedBy.Name)
	assert.Equal("developer@example.com", query.LastModifiedBy.Email)

	assert.Equal(2, len(query.Visualizations))
	queryVisualisation1 := query.Visualizations[0]
	assert.Equal(1, queryVisualisation1.ID)
	assert.Equal("TABLE", queryVisualisation1.Type)
	assert.Equal("Table", queryVisualisation1.Name)
	queryVisualisation2 := query.Visualizations[1]
	assert.Equal(2, queryVisualisation2.ID)
	assert.Equal("CHART", queryVisualisation2.Type)
	assert.Equal("DAU", queryVisualisation2.Name)
}
