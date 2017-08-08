package uaa

import (
	"strings"

	"code.cloudfoundry.org/cli/api/uaa/internal"
)

func (client *Client) GetSSHPasscode(config CCConfig) (string, error) {
	// tokens, err := client.RefreshAccessToken(config.RefreshToken())
	// if err != nil {
	// 	return "", err
	// }
	// config.SetTokenInformation(tokens.AccessToken, tokens.RefreshToken, config.SSHOauthClient())

	request, err := client.newRequest(requestOptions{
		RequestName: internal.GetOAuthAuthorizeRequest,
		Query: map[string][]string{
			"client_id":     []string{"ssh-proxy"},
			"response_type": []string{"code"},
		},
	})
	if err != nil {
		return "", err
	}

	response := Response{}
	err = client.connection.Make(request, &response)
	if err != nil {
		return "", err
	}

	return strings.Split(response.ResourceLocationURL, "=")[1], nil
}
