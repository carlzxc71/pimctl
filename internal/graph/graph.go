package graph

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	kiotajson "github.com/microsoft/kiota-serialization-json-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
)

type Client struct {
	client *msgraphsdk.GraphServiceClient
}

func NewClient(cred azcore.TokenCredential, scopes []string) (*Client, error) {
	client, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, scopes)
	if err != nil {
		return nil, fmt.Errorf("failed to create GraphServiceClient: %w", err)
	}
	return &Client{
		client: client,
	}, nil
}

type User struct {
	ID                string `json:"id"`
	DisplayName       string `json:"displayName"`
	GivenName         string `json:"givenName"`
	Surname           string `json:"surname"`
	Mail              string `json:"mail"`
	UserPrincipalName string `json:"userPrincipalName"`
}

func (c *Client) Me(ctx context.Context) (User, error) {
	res, err := c.client.Me().Get(ctx, &users.UserItemRequestBuilderGetRequestConfiguration{})
	if err != nil {
		return User{}, fmt.Errorf("failed to get me: %w", err)
	}

	user := User{}
	err = unmarshalGraphValue(res, &user)
	if err != nil {
		return User{}, fmt.Errorf("failed to unmarshal graph value: %w", err)
	}

	return user, nil
}

func unmarshalGraphValue[T any](data serialization.Parsable, v *T) error {
	kiotaBytes, err := kiotajson.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON using kiotajson: %w", err)
	}

	jsonString := fmt.Sprintf("{%s}", string(kiotaBytes))
	err = json.Unmarshal([]byte(jsonString), v)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

func prettyPrintGraphValue(data serialization.Parsable) {
	value := map[string]interface{}{}
	err := unmarshalGraphValue(data, &value)
	if err != nil {
		fmt.Printf("failed to unmarshal graph value: %v\n", err)
		return
	}
	json, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Printf("failed to pretty print json: %v\n", err)
		return
	}
	fmt.Println(string(json))
}