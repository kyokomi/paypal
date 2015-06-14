package paypal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	identityUserInfoURL = "/v1/identity/openidconnect/userinfo/?schema=openid"
)

type IdentityService struct {
	client *PayPalClient
}

/*
{
  "user_id": "https://www.paypal.com/webapps/auth/server/64ghr894040044",
  "name": "Peter Pepper",
  "given_name": "Peter",
  "family_name": "Pepper",
  "email": "ppuser@example.com"
}
*/
type UserInfoResponse struct {
	UserIDURL  string `json:"user_id"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Email      string `json:"email"`
}

func (s IdentityService) UserInfo() (UserInfoResponse, error) {
	req, err := http.NewRequest("GET", s.client.URL(identityUserInfoURL), nil)
	if err != nil {
		return UserInfoResponse{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", s.client.Authorization())

	res, err := s.client.Do(req)
	if err != nil {
		return UserInfoResponse{}, err
	}
	defer res.Body.Close()

	outData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return UserInfoResponse{}, err
	}

	if res.StatusCode >= 400 {
		return UserInfoResponse{}, fmt.Errorf("response error %d %s", res.StatusCode, string(outData))
	}

	var r UserInfoResponse
	return r, json.Unmarshal(outData, &r)

}
