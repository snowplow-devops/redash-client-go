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

package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/snowplow-devops/redash-client-go/redash"
)

func main() {

	apiKey := os.Getenv("REDASH_API_KEY")
	hostname := os.Getenv("REDASH_URL")

	log.SetLevel(log.DebugLevel)
	c, err := redash.NewClient(&redash.Config{RedashURI: hostname, APIKey: apiKey})
	if err != nil {
		log.Fatal(fmt.Errorf("Error loading client: %q", err))
		return
	}

	// --- Data source interactions

	// Get existing Data source
	dataSource, err := c.GetDataSource(1)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(fmt.Sprintf("GetDataSource - %#v", dataSource))

	// DataSource creation
	postPayload := redash.DataSource{
		Name: "My new Redshift data source",
		Type: "redshift",
		Options: map[string]interface{}{
			"host":     "localhost",
			"port":     5439,
			"dbname":   "my_database",
			"user":     "user_name",
			"password": "S3cuR3PaSsW0rD",
		},
	}

	newDataSource, err := c.CreateDataSource(&postPayload)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(fmt.Sprintf("CreateDataSource - %#v", newDataSource))

	postPayload = redash.DataSource{
		Name: "My new Redshift data source v2",
		Type: "redshift",
		Options: map[string]interface{}{
			"host":     "localhost",
			"port":     5439,
			"dbname":   "my_database",
			"user":     "user_name",
			"password": "S3cuR3PaSsW0rD",
		},
	}

	newDataSource, err = c.UpdateDataSource(newDataSource.ID, &postPayload)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(fmt.Sprintf("UpdateDataSource - %#v", newDataSource))

	// --- Group interactions
	group, err := c.GetGroup(1)
	if err != nil {
		fmt.Println(fmt.Errorf("Error retreiving group: %q", err))
		return
	}
	fmt.Println(fmt.Sprintf("GetGroup - %#v", group))

	// Create a new group
	groupPayload := redash.GroupCreatePayload{
		Name: "com.acme group",
	}

	newGroup, err := c.CreateGroup(&groupPayload)
	if err != nil {
		fmt.Println(fmt.Errorf("Error creating group: %q", err))
		return
	}
	fmt.Println(fmt.Sprintf("CreateGroup - %#v", newGroup))

	// Add a user to new group
	err = c.GroupAddUser(newGroup.ID, 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	dataSource, err = c.GetDataSource(newDataSource.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(fmt.Sprintf("GroupAddUser>GetDataSource - %#v", dataSource))

	// Add a data source to new group
	err = c.GroupAddDataSource(newGroup.ID, newDataSource.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	dataSource, err = c.GetDataSource(newDataSource.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(fmt.Sprintf("GroupAddDataSource>GetDataSource - %#v", dataSource))

	// Remove user from new group
	err = c.GroupRemoveUser(newGroup.ID, 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	newGroup, err = c.GetGroup(newGroup.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(fmt.Sprintf("GroupRemoveUser>GetGroup - %#v", newGroup))

	// --- Cleanup
	err = c.GroupRemoveDataSource(newGroup.ID, newDataSource.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	dataSource, err = c.GetDataSource(newDataSource.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(fmt.Sprintf("GroupRemoveDataSource>GetDataSource - %#v", dataSource))

	// Delete data source
	err = c.DeleteDataSource(newDataSource.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	// --- Alert destination interactions
	// Create a new alert destination
	destinationPayload := []byte(`{"name": "Slack", "type": "slack", "options": {}}`)
	newDestination, err := c.CreateDestination(destinationPayload)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Desitnation %#v created", newDestination.Name)

	// Get the destination
	getDestinationResponse, err := c.GetDestination(newDestination.ID)
	if err != nil {
		log.Fatal(err)
		return
	}
	destination, ok := getDestinationResponse.(*redash.SlackDestination)
	if !ok {
		log.Fatal(err)
		return
	}
	fmt.Printf("The desitnation is %#v", destination.Name)

	// Update the destination
	updateDestinationPayload := []byte(`{"name": "Slack", "type": "slack", "options": {"url": "https://test.slack.com/hook"}}`)
	updateDestionationResponse, err := c.UpdateDestination(newDestination.ID, updateDestinationPayload)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("The desitnation %#v updated", updateDestionationResponse.Name)

	// Delete the destionation
	err = c.DeleteDestination(newDestination.ID)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("The desitnation %#v deleted", destination.Name)
}
