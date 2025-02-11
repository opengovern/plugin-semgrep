package models

type IntegrationCredentials struct {
	Token        string `json:"token"`
	Organization string `json:"organization"`
}
