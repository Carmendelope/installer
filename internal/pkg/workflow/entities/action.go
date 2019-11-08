/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
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
