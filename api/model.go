package api

import "time"

type VM struct {
	ID         string `json:"ID"`
	OwnerID    string `json:"ownerID"`
	Name       string `json:"name"`
	InstanceID string `json:"instanceID"`
	PublicDNS  string `json:"publicDNS"`
	PublicIP   string `json:"publicIP"`
	KeyName    string `json:"keyName"`
	Zone       string `json:"zone"`
}

type PrivateKey struct {
	KeyName        string `json:"keyName"`
	KeyMaterial    string `json:"keyMaterial"`
	KeyFingerprint string `json:"keyFingerprint"`
}

type TokenDetails struct {
	AccessToken     string
	RefreshToken    string
	AccessTokenExp  time.Time
	RefreshTokenExp time.Time
}
