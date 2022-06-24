package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

const server = "http://localhost:8000/api/sources/v3.1"

// x-hr-identity {"identity":{"account_number":"12cdcb0e-e7e1-11ec-97d8-0242ac110003"}}
// for tenant id 3 in my db
const xHrIdentity = "eyJpZGVudGl0eSI6eyJhY2NvdW50X251bWJlciI6IjEyY2RjYjBlLWU3ZTEtMTFlYy05N2Q4LTAyNDJhYzExMDAwMyJ9fQ=="

type Collection struct {
	Data  []interface{} `json:"data"`
	Meta  Metadata      `json:"meta"`
	Links Links         `json:"links"`
}

type Metadata struct {
	Count  int `json:"count"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type Links struct {
	First string `json:"first"`
	Last  string `json:"last"`
}

func main() {
	client, err := NewClientWithResponses(server)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// List Sources without filter
	response, err := client.ListSourcesWithResponse(ctx, &ListSourcesParams{}, AddIdentityHeader)
	if err != nil {
		panic(err)
	}

	var out Collection
	err = json.Unmarshal(response.Body, &out)
	if err != nil {
		panic(err)
	}

	// In my db the correct result are 14 sources
	if out.Meta.Count != 14 {
		log.Print("count not equal to 14")
	}

	// Add filter
	name := "ibm"
	filter := ListSourcesParams_Filter{map[string]string{"[source_type][name]": name}}

	// List Sources with filter
	response, err = client.ListSourcesWithResponse(ctx, &ListSourcesParams{Filter: &filter}, AddIdentityHeader)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(response.Body, &out)
	if err != nil {
		panic(err)
	}

	// In my db the correct result are 2 sources
	if out.Meta.Count != 2 {
		log.Print("count not equal to 2")
	}
}

func AddIdentityHeader(ctx context.Context, req *http.Request) error {
	req.Header.Set("x-rh-identity", xHrIdentity)
	return nil
}

