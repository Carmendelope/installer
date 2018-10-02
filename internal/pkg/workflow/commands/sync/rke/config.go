/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

// Management structures for the RKE config file

package rke

// ClusterConfig defines the options required to generate an RKE config file.
type ClusterConfig struct {
	ClusterName    string   `json:"clusterName"`
	TargetNodes    []string `json:"targetNodes"`
	NodeUsername   string   `json:"nodeUsername"`
	PrivateKeyPath string   `json:"privateKeyPath"`
}

// NewClusterConfig creates a new set of config parameters for creating the RKE definition file
func NewClusterConfig(
	clusterName string,
	targetNodes []string,
	nodeUsername string,
	privateKeyPath string) *ClusterConfig {
	return &ClusterConfig{
		clusterName,
		targetNodes,
		nodeUsername,
		privateKeyPath}
}
