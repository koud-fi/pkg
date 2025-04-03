package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/fetch"
)

// DeviceCodeResponse represents the response from the device code endpoint.
type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"` // seconds until expiration
	Interval        int    `json:"interval"`   // polling interval in seconds
	Message         string `json:"message"`
}

// RequestDeviceCode initiates the device authorization flow.
// It sends a POST request to the endpoint specified in conf.AuthBaseURL.
func RequestDeviceCode(conf Config) (DeviceCodeResponse, error) {
	// Build the scope string.
	if conf.ScopeFormat == Default {
		conf.ScopeFormat = Space
	}
	var scopeBuilder strings.Builder
	for i, s := range conf.Scopes {
		if i > 0 {
			scopeBuilder.WriteByte(byte(conf.ScopeFormat))
		}
		scopeBuilder.WriteString(string(s))
	}

	data := url.Values{
		"client_id": {conf.ClientID},
		"scope":     {scopeBuilder.String()},
	}

	req := fetch.Post(conf.AuthBaseURL).
		Body(blob.FromString(data.Encode()), "application/x-www-form-urlencoded")

	deviceResp, err := blob.UnmarshalJSONValue[DeviceCodeResponse](req)
	if err != nil {
		return DeviceCodeResponse{}, err
	}
	return deviceResp, nil
}

// PollForToken continuously polls the token endpoint (conf.TokenBaseURL) until an access token is issued.
// It uses the device code provided in the DeviceCodeResponse.
func PollForToken(conf Config, deviceCode string, interval int, expiresIn int) (Token, error) {
	deadline := time.Now().Add(time.Duration(expiresIn) * time.Second)
	pollInterval := time.Duration(interval) * time.Second
	client := &http.Client{}

	for {
		if time.Now().After(deadline) {
			return Token{}, errors.New("device code expired")
		}

		data := url.Values{}
		data.Set("client_id", conf.ClientID)
		data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
		data.Set("device_code", deviceCode)
		if conf.ClientSecret != "" {
			data.Set("client_secret", conf.ClientSecret)
		}

		req, err := http.NewRequest("POST", conf.TokenBaseURL, strings.NewReader(data.Encode()))
		if err != nil {
			return Token{}, err
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err != nil {
			return Token{}, err
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return Token{}, err
		}

		if resp.StatusCode == http.StatusOK {
			var token Token
			if err := json.Unmarshal(body, &token); err != nil {
				return Token{}, err
			}
			return token, nil
		}

		// For non-200 responses, decode the error response.
		var errResp CallbackError
		if err := json.Unmarshal(body, &errResp); err != nil {
			return Token{}, errors.New("failed to decode error response: " + string(body))
		}

		//dump.AsJSON(errResp)
		fmt.Println(string(body))

		switch errResp.Code {
		case "authorization_pending":
			// User hasn't completed the flow yet; wait and continue polling.
			time.Sleep(pollInterval)
			continue

		case "slow_down":
			// Increase polling interval if requested.
			pollInterval += time.Second
			time.Sleep(pollInterval)
			continue

		default:
			return Token{}, errors.New("failed to obtain access token: " + string(body))
		}
	}
}
