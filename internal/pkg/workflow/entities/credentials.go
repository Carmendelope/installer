/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

// Definition of the credentials entity.

package entities

// Credentials structure with the information required to connect to remote hosts.
type Credentials struct {
	// Username for the SSH credentials.
	Username string `json:"username"`
	// Password for the SSH credentials.
	Password string `json:"password"`
	// PrivateKey alternative for the credentials.
	PrivateKey string `json:"privateKey"`
}

// NewCredentials creates a new Credentials structure.
//   params:
//     username The SSH username.
//     password The SSH password.
//   returns:
//     A credentials instance.
func NewCredentials(username string, password string) *Credentials {
	return &Credentials{username, password, ""}
}

// NewPKICredentials creates a credentials entity using ssh public key auth.
func NewPKICredentials(username string, sshKey string) *Credentials {
	return &Credentials{username, "", sshKey}
}

// UsePKI determines is the credentials should be used as username/password or with public key.
func (c *Credentials) UsePKI() bool {
	return c.PrivateKey != ""
}

type InstallCredentials struct {
	// Username for the SSH credentials.
	Username string `json:"username"`
	// PrivateKeyPath with the path of the private key.
	PrivateKeyPath string `json:"privateKeyPath"`
}

// NewInstallCredentials creates a InstallCredentials structure with all parameters.
func NewInstallCredentials(username string, privateKeyPath string) *InstallCredentials {
	return &InstallCredentials{username, privateKeyPath}
}
