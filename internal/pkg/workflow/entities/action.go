/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

// This file contains the action definition

package entities

import (
	"fmt"
	"strings"
)

// Action structure that defines remote commands.
type Action struct {
	// Name of the action.
	Name string `json:"name"`
	// Args is the set of arguments.
	Args []string `json:"args"`
}

// NewAction creates a new action with all parameters.
func NewAction(name string, args []string) *Action {
	return &Action{name, args}
}

// String returns the string representation fo the given action.
func (a *Action) String() string {
	return fmt.Sprintf("%s - %s", a.Name, strings.Join(a.Args, " "))
}

// PrettyPrint returns a simple space indexed string.
func (a *Action) PrettyPrint(indentation int) string {
	return fmt.Sprintf("%s%s", strings.Repeat(" ", indentation), a.String())
}
