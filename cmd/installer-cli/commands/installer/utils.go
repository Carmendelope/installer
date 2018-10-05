/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package installer

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
	if strings.HasPrefix(path, "."){
		abs, _ := filepath.Abs("./")
		return strings.Replace(path, ".", abs, 1)
	}
	return path
}

func GetKubeConfigContent(kubeConfigPath string) (string, derrors.Error) {
	if kubeConfigPath == "" {
		return "", nil
	}

	content, err := ioutil.ReadFile(GetPath(kubeConfigPath))
	if err != nil{
		return "", derrors.AsError(err, "cannot read kubeconfig file")
	}
	// TODO Check the contents of the kubeconfig file to make sure the keys are embedded
	return string(content), nil
}

func GetPrivateKeyContent(privateKeyPath string) (string, derrors.Error){
	if privateKeyPath == "" {
		return "", nil
	}
	content, err := ioutil.ReadFile(GetPath(privateKeyPath))
	if err != nil{
		return "", derrors.AsError(err, "cannot read kubeconfig file")
	}
	return string(content), nil
}