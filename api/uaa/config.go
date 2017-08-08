package uaa

type CCConfig interface {
	SetTokenInformation(accessToken string, refreshToken string, sshOAuthClient string)
	SSHOauthClient() string
	RefreshToken() string
}
