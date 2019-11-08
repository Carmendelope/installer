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

package utils

import (
	"github.com/nalej/derrors"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"strings"
)

func GetPath(path string) string {
	if strings.HasPrefix(path, "~") {
		usr, _ := user.Current()
		return strings.Replace(path, "~", usr.HomeDir, 1)
	}
	if strings.HasPrefix(path, ".") {
		abs, _ := filepath.Abs("./")
		return strings.Replace(path, ".", abs, 1)
	}
	return path
}

func ExtendComponentsPath(path string, appClusterInstall bool) string {
	if appClusterInstall {
		return filepath.Join(path, "appcluster")
	}
	return filepath.Join(path, "mngtcluster")
}

func GetKubeConfigContent(kubeConfigPath string) (string, derrors.Error) {
	if kubeConfigPath == "" {
		return "", nil
	}

	content, err := ioutil.ReadFile(GetPath(kubeConfigPath))
	if err != nil {
		return "", derrors.AsError(err, "cannot read kubeconfig file")
	}
	// TODO Check the contents of the kubeconfig file to make sure the keys are embedded
	return string(content), nil
}

func GetPrivateKeyContent(privateKeyPath string) (string, derrors.Error) {
	if privateKeyPath == "" {
		return "", nil
	}
	content, err := ioutil.ReadFile(GetPath(privateKeyPath))
	if err != nil {
		return "", derrors.AsError(err, "cannot read kubeconfig file")
	}
	return string(content), nil
}
