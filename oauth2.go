package paypal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	authTokenURL = "/v1/oauth2/token"
)

/*
{
  "scope": "https://api.paypal.com/v1/payments/.* https://api.paypal.com/v1/vault/credit-card https://api.paypal.com/v1/vault/credit-card/.*",
  "access_token": "<Access-Token>",
  "token_type": "Bearer",
  "app_id": "APP-6XR95014SS315863X",
  "expires_in": 28800
}
*/
type AdminAuthToken struct {
	Scope       string `json:"scope"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	AppID       string `json:"app_id"`
	ExpiresIn   int    `json:"expires_in"`
}

func (a AdminAuthToken) ScopeSlice() []string {
	return strings.Fields(a.Scope)
}

func (a AdminAuthToken) Authorization() string {
	return fmt.Sprintf("%s %s", a.TokenType, a.AccessToken)
}

type OAuth2Service struct {
	client *PayPalClient
}

func (s OAuth2Service) GetToken() (AdminAuthToken, error) {
	opts := s.client.Options

	val := url.Values{}
	val.Add("grant_type", "client_credentials")
	req, err := http.NewRequest("POST", s.client.URL(authTokenURL), bytes.NewBufferString(val.Encode()))
	if err != nil {
		return AdminAuthToken{}, err
	}

	req.SetBasicAuth(opts.ClientID, opts.Secret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Language", opts.AcceptLanguage)

	res, err := s.client.Do(req)
	if err != nil {
		return AdminAuthToken{}, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return AdminAuthToken{}, err
	}

	var token AdminAuthToken
	return token, json.Unmarshal(data, &token)
}
