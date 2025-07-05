package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/crossplane-contrib/provider-http/apis/common"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
)

const (
	authKey = "Authorization"
)

// Client is the interface to interact with Http
type Client interface {
	SendRequest(ctx context.Context, method string, url string, body Data, headers Data, skipTLSVerify bool) (resp HttpDetails, err error)
	SendRequestWithTLS(ctx context.Context, method string, url string, body Data, headers Data, tlsConfig *common.TLSConfig) (resp HttpDetails, err error)
}

type client struct {
	log                logging.Logger
	timeout            time.Duration
	authorizationToken string
}

type HttpResponse struct {
	Body       string              `json:"body"`
	Headers    map[string][]string `json:"headers"`
	StatusCode int                 `json:"statusCode"`
}

type Data struct {
	Encrypted interface{} // Data containing encrypted data -> to be shown at the status
	Decrypted interface{} // Data containing sensitive data -> to be sent
}

type HttpRequest struct {
	Method  string              `json:"method"`
	Body    string              `json:"body,omitempty"`
	URL     string              `json:"url"`
	Headers map[string][]string `json:"headers,omitempty"`
}

type HttpDetails struct {
	HttpResponse HttpResponse
	HttpRequest  HttpRequest
}

// SendRequest sends an HTTP request to the specified URL with the given method, body, headers and skipTLSVerify.
func (hc *client) SendRequest(ctx context.Context, method string, url string, body Data, headers Data, skipTLSVerify bool) (details HttpDetails, err error) {
	requestBody := []byte(body.Decrypted.(string))

	// request contains the HTTP request that will be sent.
	request, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(requestBody))

	// requestDetails contains the request details that will be logged.
	requestDetails := HttpRequest{
		URL:     url,
		Body:    body.Encrypted.(string),
		Headers: headers.Encrypted.(map[string][]string),
		Method:  method,
	}

	if err != nil {
		return HttpDetails{
			HttpRequest: requestDetails,
		}, err
	}

	for key, values := range headers.Decrypted.(map[string][]string) {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	// Add the authorization token to the request if it doesn't already exist.
	if _, exists := request.Header[authKey]; !exists && hc.authorizationToken != "" {
		request.Header[authKey] = []string{hc.authorizationToken}
	}

	client := &http.Client{
		Transport: &http.Transport{
			// #nosec G402
			TLSClientConfig: &tls.Config{InsecureSkipVerify: skipTLSVerify},
			Proxy:           http.ProxyFromEnvironment, // Use proxy settings from environment
		},
		Timeout: hc.timeout,
	}

	response, err := client.Do(request)
	if err != nil {
		return HttpDetails{
			HttpRequest: requestDetails,
		}, err
	}

	responsebody, err := io.ReadAll(response.Body)
	if err != nil {
		return HttpDetails{
			HttpRequest: requestDetails,
		}, err
	}

	beautifiedResponse := HttpResponse{
		Body:       string(responsebody),
		Headers:    response.Header,
		StatusCode: response.StatusCode,
	}

	err = response.Body.Close()
	if err != nil {
		return HttpDetails{
			HttpRequest: requestDetails,
		}, err
	}

	hc.log.Info(fmt.Sprint("http request sent: ", toJSON(requestDetails)))

	return HttpDetails{
		HttpResponse: beautifiedResponse,
		HttpRequest:  requestDetails,
	}, nil
}

// NewClient returns a new Http Client
func NewClient(log logging.Logger, timeout time.Duration, authorizationToken string) (Client, error) {
	return &client{
		log:                log,
		timeout:            timeout,
		authorizationToken: authorizationToken,
	}, nil
}

// toJSON converts the request to a JSON string.
func toJSON(request HttpRequest) string {
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		return ""
	}

	return string(jsonBytes)
}

// SendRequestWithTLS sends an HTTP request with TLS configuration
func (hc *client) SendRequestWithTLS(ctx context.Context, method string, url string, body Data, headers Data, tlsConfig *common.TLSConfig) (details HttpDetails, err error) {
	// Build TLS configuration
	tlsClientConfig, err := hc.buildTLSConfig(tlsConfig)
	if err != nil {
		return HttpDetails{}, fmt.Errorf("failed to build TLS configuration: %w", err)
	}

	requestBody := []byte(body.Decrypted.(string))

	// request contains the HTTP request that will be sent.
	request, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(requestBody))

	// requestDetails contains the request details that will be logged.
	requestDetails := HttpRequest{
		URL:     url,
		Body:    body.Encrypted.(string),
		Headers: headers.Encrypted.(map[string][]string),
		Method:  method,
	}

	if err != nil {
		return HttpDetails{
			HttpRequest: requestDetails,
		}, err
	}

	for key, values := range headers.Decrypted.(map[string][]string) {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	// Add the authorization token to the request if it doesn't already exist.
	if _, exists := request.Header[authKey]; !exists && hc.authorizationToken != "" {
		request.Header[authKey] = []string{hc.authorizationToken}
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsClientConfig,
			Proxy:           http.ProxyFromEnvironment, // Use proxy settings from environment
		},
		Timeout: hc.timeout,
	}

	response, err := client.Do(request)
	if err != nil {
		return HttpDetails{
			HttpRequest: requestDetails,
		}, err
	}

	responsebody, err := io.ReadAll(response.Body)
	if err != nil {
		return HttpDetails{
			HttpRequest: requestDetails,
		}, err
	}

	beautifiedResponse := HttpResponse{
		Body:       string(responsebody),
		Headers:    response.Header,
		StatusCode: response.StatusCode,
	}

	err = response.Body.Close()
	if err != nil {
		return HttpDetails{
			HttpRequest: requestDetails,
		}, err
	}

	hc.log.Info(fmt.Sprint("http request sent: ", toJSON(requestDetails)))

	return HttpDetails{
		HttpResponse: beautifiedResponse,
		HttpRequest:  requestDetails,
	}, nil
}

// buildTLSConfig creates a tls.Config based on the provided common.TLSConfig
func (hc *client) buildTLSConfig(tlsConfig *common.TLSConfig) (*tls.Config, error) {
	if tlsConfig == nil {
		return &tls.Config{}, nil
	}

	config := &tls.Config{
		InsecureSkipVerify: tlsConfig.InsecureSkipVerify,
	}

	// Handle CA certificate
	if tlsConfig.CAData != "" {
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM([]byte(tlsConfig.CAData)) {
			return nil, fmt.Errorf("failed to parse CA certificate data")
		}
		config.RootCAs = caCertPool
	}

	// Handle client certificate and key for mutual TLS
	if tlsConfig.ClientCertData != "" && tlsConfig.ClientKeyData != "" {
		cert, err := tls.X509KeyPair([]byte(tlsConfig.ClientCertData), []byte(tlsConfig.ClientKeyData))
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate and key: %w", err)
		}
		config.Certificates = []tls.Certificate{cert}
	} else if tlsConfig.ClientCertData != "" || tlsConfig.ClientKeyData != "" {
		return nil, fmt.Errorf("both client certificate and key must be provided for mutual TLS")
	}

	return config, nil
}
