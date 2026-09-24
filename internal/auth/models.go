package auth

type Credentials struct {
	Installed struct {
		ClientID     string   `json:"client_id"`
		ProjectID    string   `json:"project_id"`
		AuthURI      string   `json:"auth_uri"`
		AuthProvider string   `json:"auth_provider_x509_cert_url"`
		ClientSecret string   `json:"client_secret"`
		RedirectURIs []string `json:"redirect_uris"`
	} `json:"installed"`
}
