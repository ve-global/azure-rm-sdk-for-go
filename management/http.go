package management

import (
	"bytes"
	"fmt"
	"net/http"
	"io/ioutil"
)

const (
	msVersionHeader           = "x-ms-version"
	requestIDHeader           = "x-ms-request-id"
	uaHeader                  = "User-Agent"
	contentHeader             = "Content-Type"
	defaultContentHeaderValue = "application/json"
)



func (client client) SendAzureGetRequest(url string) ([]byte, error) {

	resp, err := client.sendAzureRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return getResponseBody(resp)
}


func (client client) SendAzurePutRequest(url string, data []byte) (error) {
	_,err := client.doAzureOperation("PUT", url, data)
	return err
}

func (client client) SendAzurePutRequestWithReturnedResponse(url string, data []byte) ([]byte, error) {
	resp, err := client.doAzureOperation("PUT", url, data)
	if err != nil{
		return nil,err
	}
	return getResponseBody(resp)
}


func (client client) SendAzureDeleteRequest(url string) (error) {
	_,err := client.doAzureOperation("DELETE", url, nil)
	return err
}

func (client client) doAzureOperation(method, url string, data []byte) (*http.Response, error) {
	resp, err := client.sendAzureRequest(method, url, data)
	if err != nil {
		return resp, err
	}
	return resp,err
}

// sendAzureRequest constructs an HTTP client for the request, sends it to the
// management API and returns the response or an error.
func (client client) sendAzureRequest(method , url string, data []byte) (*http.Response, error) {
	if method == "" {
		return nil, fmt.Errorf(errParamNotSpecified, "method")
	}
	if url == "" {
		return nil, fmt.Errorf(errParamNotSpecified, "url")
	}

	httpClient := client.createHTTPClient()

	response, err := client.sendRequest(httpClient, url, method, data, 5)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// createHTTPClient creates an HTTP Client configured with the key pair for
// the subscription for this client.
func (client client) createHTTPClient() *http.Client {
	return nil
}

// sendRequest sends a request to the Azure management API using the given
// HTTP client and parameters. It returns the response from the call or an
// error.
func (client client) sendRequest(httpClient *http.Client, url, requestType string, data []byte, numberOfRetries int) (*http.Response, error) {
	request, reqErr := client.createAzureRequest(url, requestType , data)
	if reqErr != nil {
		return nil, reqErr
	}

	response, err := httpClient.Do(request)
	if err != nil {
		if numberOfRetries == 0 {
			return nil, err
		}

		return client.sendRequest(httpClient, url, requestType, data, numberOfRetries-1)
	}

	if response.StatusCode >= http.StatusBadRequest {
		body, err := getResponseBody(response)
		if err != nil {
			// Failed to read the response body
			return nil, err
		}
		azureErr := getAzureError(body)
		if azureErr != nil {
			if numberOfRetries == 0 {
				return nil, azureErr
			}

			return client.sendRequest(httpClient, url, requestType, data, numberOfRetries-1)
		}
	}

	return response, nil
}

// createAzureRequest packages up the request with the correct set of headers and returns
// the request object or an error.
func (client client) createAzureRequest(url string, requestType string, data []byte) (*http.Request, error) {
	var request *http.Request
	var err error

	if data != nil {
		body := bytes.NewBuffer(data)
		request, err = http.NewRequest(requestType, url, body)
	} else {
		request, err = http.NewRequest(requestType, url, nil)
	}

	if err != nil {
		return nil, err
	}

	//request.Header.Set(msVersionHeader, client.config.APIVersion)

	request.Header.Set(uaHeader, client.config.UserAgent)
	request.Header.Set(contentHeader, defaultContentHeaderValue)

	return request, nil
}

func getResponseBody(response *http.Response) ([]byte, error) {
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

func getAzureError(body []byte) (error){
	return nil
}
