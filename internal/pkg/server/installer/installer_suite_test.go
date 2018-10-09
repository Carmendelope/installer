/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package installer

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestApplicationPackage(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Installer package suite")
}