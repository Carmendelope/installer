/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package sync

import (
	"encoding/json"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"os"
	"strings"

	"github.com/nalej/derrors"
)

// CheckAsset structure with the command parameters.
type CheckAsset struct {
	entities.GenericSyncCommand
	Path string `json:"path"`
}

// NewCheckAsset creates a new CheckAsset command.
func NewCheckAsset(path string) *CheckAsset {
	return &CheckAsset{*entities.NewSyncCommand(entities.CheckAsset), path}
}

// NewCheckAssetFromJSON creates a new CheckAsset command using a raw JSON payload.
func NewCheckAssetFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	checkAsset := &CheckAsset{}
	if err := json.Unmarshal(raw, &checkAsset); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	checkAsset.CommandID = entities.GenerateCommandID(checkAsset.Name())
	var r entities.Command = checkAsset
	return &r, nil
}

// Run the current command.
//   returns:
//     The CommandResult
//     An error if the command execution fails
func (ca *CheckAsset) Run(_ string) (*entities.CommandResult, derrors.Error) {
	// Check if the file exists.
	fileInfo, err := os.Stat(ca.Path)

	if os.IsNotExist(err) {
		return nil, derrors.NewNotFoundError(errors.FileDoesNotExist, err)
	}
	if err != nil {
		return nil, derrors.NewInternalError(errors.IOError, err)
	}
	if fileInfo.IsDir() {
		return nil, derrors.NewInvalidArgumentError(errors.ExpectingFile, err)
	}
	return entities.NewSuccessCommand([]byte("OK")), nil
}

// String obtains a string representation
func (ca *CheckAsset) String() string {
	return "SYNC CheckAsset " + ca.Path
}

// PrettyPrint returns a simple space indexed string.
func (ca *CheckAsset) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + ca.String()
}

// UserString returns a simple string representation of the command for the user.
func (ca *CheckAsset) UserString() string {
	return "Check asset " + ca.Path
}
