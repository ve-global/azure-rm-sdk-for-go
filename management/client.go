// Package management provides the main API client to construct other clients
// and make requests to the Microsoft Azure Service Management REST API.
package management

import (
	"errors"
)

const (
	DefaultAzureTokenUrl		 = "https://login.microsoftonline.com/%s.onmicrosoft.com/oauth2/token"
	DefaultAzureManagementURL    = "https://management.azure.com/subscriptions/%s"
	DefaultUserAgent             = "azure-rm-sdk-for-go"

	errPublishSettingsConfiguration       = "PublishSettingsFilePath is set. Consequently ManagementCertificatePath and SubscriptionId must not be set."
	errManagementCertificateConfiguration = "Both ManagementCertificatePath and SubscriptionId should be set, and PublishSettingsFilePath must not be set."
	errParamNotSpecified                  = "Parameter %s is not specified."
)

type client struct {
	azureConfig     azureConfig
	config 			clientConfig
}

// Client is the base Azure Service Management API client instance that
// can be used to construct client instances for various services.
type Client interface {

	// SendAzureGetRequest sends a request to the management API using the HTTP GET method
	// and returns the response body or an error.
	SendAzureGetRequest(url string) ([]byte, error)

	// SendAzurePutRequest sends a request to the management API using the HTTP PUT method
	// and returns the request ID or an error.
	SendAzurePutRequest(url, data []byte) (error)

	// SendAzurePutRequest sends a request to the management API using the HTTP PUT method
	// and returns the request ID or an error.
	SendAzurePutRequestWithReturnedResponse(url, data []byte) ([]byte, error)

	// SendAzureDeleteRequest sends a request to the management API using the HTTP DELETE method
	// and returns the an error.
	SendAzureDeleteRequest(url string) (error)

}

// ClientConfig provides a configuration for use by a Client.
type clientConfig struct {
	AzureTokenUrl      	  string
	AzureManagementURL 	  string
	UserAgent 	      	  string
	ManagementURL         string
}




func makeClient(subscriptionID string, applicationID string, applicationSecret string,resourceGroup string, config clientConfig) (client, error) {
	var c client

	if subscriptionID == "" {
		return c, errors.New("azure: subscription ID required")
	}

	if applicationSecret == "" {
		return c, errors.New("azure: appliication secret required")
	}


	azureConfig := azureConfig{
		SubscriptionID:   subscriptionID,
		ApplicationId: applicationID,
		ApplicationSecret:  applicationSecret,
		ResourceGroup: resourceGroup,
	}

	currentClient := client{
		azureConfig      : azureConfig,
		config			 : config,
	}

	return currentClient, nil
}
