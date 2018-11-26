/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-user-manager-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/rs/zerolog/log"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

const ClusterUserSecretName = "cluster-user-credentials"

type AddClusterUser struct {
	Kubernetes
	OrganizationID     string `json:"organization_id"`
	ClusterID          string `json:"cluster_id"`
	UserManagerAddress string `json:"user_manager_address"`
}

func NewAddClusterUser(kubeConfigPath string, organizationID string, clusterID string, userManagerAddress string) *AddClusterUser {
	return &AddClusterUser{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.AddClusterUser),
			KubeConfigPath:     kubeConfigPath,
		},
		OrganizationID:     organizationID,
		ClusterID:          clusterID,
		UserManagerAddress: userManagerAddress,
	}
}

// NewAddClusterUserFromJSON creates an AddClusterUser command from a JSON object.
func NewAddClusterUserFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	lc := &AddClusterUser{}
	if err := json.Unmarshal(raw, &lc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	lc.CommandID = entities.GenerateCommandID(lc.Name())
	var r entities.Command = lc
	return &r, nil
}

func (acu *AddClusterUser) getRoleId(client grpc_user_manager_go.UserManagerClient) (string, derrors.Error) {
	orgID := &grpc_organization_go.OrganizationId{
		OrganizationId: acu.OrganizationID,
	}
	roles, err := client.ListRoles(context.Background(), orgID)
	if err != nil {
		return "", conversions.ToDerror(err)
	}
	// Find the cluster role
	for _, r := range roles.Roles {
		if r.Name == "AppCluster" {
			return r.RoleId, nil
		}
	}
	return "", derrors.NewNotFoundError("cannot find AppCluster role")
}

func (acu *AddClusterUser) createNewUser(roleID string, client grpc_user_manager_go.UserManagerClient) (string, string, derrors.Error) {

	addUserRequest := &grpc_user_manager_go.AddUserRequest{
		OrganizationId: acu.OrganizationID,
		Email:          fmt.Sprintf("%s@nalej.internal", acu.ClusterID),
		Password:       uuid.NewV4().String(),
		Name:           acu.ClusterID,
		RoleId:         roleID,
	}

	added, err := client.AddUser(context.Background(), addUserRequest)

	if err != nil {
		return "", "", conversions.ToDerror(err)
	}
	return added.Email, addUserRequest.Password, nil
}

func (acu *AddClusterUser) storeClusterUserCredentials(email string, password string) derrors.Error {
	secret := &v1.Secret{
		TypeMeta: v12.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: v12.ObjectMeta{
			Name:      ClusterUserSecretName,
			Namespace: "nalej",
		},
		StringData: map[string]string{
			"email":    email,
			"password": password,
		},
		Type: v1.SecretTypeBasicAuth,
	}

	client := acu.Client.CoreV1().Secrets(secret.Namespace)
	created, err := client.Create(secret)
	if err != nil {
		return derrors.AsError(err, "cannot create secret")
	}
	log.Debug().Interface("created", created).Msg("secret has been created")
	return nil
}

func (acu *AddClusterUser) Run(workflowID string) (*entities.CommandResult, derrors.Error) {

	connectErr := acu.Connect()
	if connectErr != nil {
		return nil, connectErr
	}

	umConn, err := grpc.Dial(acu.UserManagerAddress, grpc.WithInsecure())
	if err != nil {
		return nil, derrors.AsError(err, "cannot create connection with the user manager")
	}
	defer umConn.Close()
	client := grpc_user_manager_go.NewUserManagerClient(umConn)

	roleId, dErr := acu.getRoleId(client)
	if dErr != nil {
		return entities.NewCommandResult(
			false, "cannot determine cluster role", dErr), nil
	}

	email, password, dErr := acu.createNewUser(roleId, client)
	if dErr != nil {
		return entities.NewCommandResult(
			false, "cannot add cluster user", dErr), nil
	}

	dErr = acu.storeClusterUserCredentials(email, password)
	if dErr != nil {
		return entities.NewCommandResult(
			false, "cannot determine store cluster credentials", dErr), nil
	}
	return entities.NewSuccessCommand([]byte("cluster credentials has been created")), nil
}

func (acu *AddClusterUser) String() string {
	return fmt.Sprintf("SYNC AddClusterUser")
}

func (acu *AddClusterUser) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + acu.String()
}

func (acu *AddClusterUser) UserString() string {
	return fmt.Sprintf("Creating cluster user")
}
