// Copyright 2013 Canonical Ltd.  This software is licensed under the
// GNU Lesser General Public License version 3 (see the file COPYING).

package gomaasapi

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	Signer OAuthSigner
}

const (
	operationParamName = "op"
)

func (client Client) dispatchRequest(request *http.Request) ([]byte, error) {
	client.Signer.OAuthSign(request)
	httpClient := http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode/100 != 2 {
		return body, errors.New("Error requesting the MAAS server: " + response.Status + ".")
	}
	return body, nil
}

func (client Client) Get(URL string, operation string, parameters url.Values) ([]byte, error) {
	opParameter := parameters.Get(operationParamName)
	if opParameter != "" {
		errString := fmt.Sprintf("The parameters contain a value for '%s' which is reserved parameter.")
		return nil, errors.New(errString)
	}
	if operation != "" {
		parameters.Set(operationParamName, operation)
	}
	queryUrl := URL + "?" + parameters.Encode()
	request, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		return nil, err
	}
	return client.dispatchRequest(request)
}

// nonIdempotentRequest is a utility method to issue a PUT or a POST request.
func (client Client) nonIdempotentRequest(method string, URL string, parameters url.Values) ([]byte, error) {
	request, err := http.NewRequest(method, URL, strings.NewReader(string(parameters.Encode())))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}
	return client.dispatchRequest(request)
}

func (client Client) Post(URL string, operation string, parameters url.Values) ([]byte, error) {
	queryParams := url.Values{operationParamName: {operation}}
	queryURL := URL + "?" + queryParams.Encode()
	return client.nonIdempotentRequest("POST", queryURL, parameters)
}

func (client Client) Put(URL string, parameters url.Values) ([]byte, error) {
	return client.nonIdempotentRequest("PUT", URL, parameters)
}

func (client Client) Delete(URL string) error {
	request, err := http.NewRequest("DELETE", URL, strings.NewReader(""))
	if err != nil {
		return err
	}
	_, err2 := client.dispatchRequest(request)
	if err2 != nil {
		return err2
	}
	return nil
}

type anonSigner struct{}

func (signer anonSigner) OAuthSign(request *http.Request) error {
	return nil
}

// Trick to ensure *anonSigner implements the OAuthSigner interface.
var _ OAuthSigner = (*anonSigner)(nil)

// NewAnonymousClient creates a client that issues anonymous requests.
func NewAnonymousClient() (*Client, error) {
	return &Client{Signer: &anonSigner{}}, nil
}

// NewAuthenticatedClient parses the given MAAS API key into the individual
// OAuth tokens and creates an Client that will use these tokens to sign the
// requests it issues.
func NewAuthenticatedClient(apiKey string) (*Client, error) {
	elements := strings.Split(apiKey, ":")
	if len(elements) != 3 {
		errString := "Invalid API key. The format of the key must be \"<consumer secret>:<token key>:<token secret>\"."
		err := errors.New(errString)
		return nil, err
	}
	// The consumer secret is the empty string in MAAS' authentication.
	token := &OAuthToken{
		ConsumerKey:    elements[0],
		ConsumerSecret: "",
		TokenKey:       elements[1],
		TokenSecret:    elements[2],
	}
	signer, err := NewPLAINTEXTOAuthSigner(token, "MAAS API")
	if err != nil {
		return nil, err
	}
	return &Client{Signer: signer}, nil
}