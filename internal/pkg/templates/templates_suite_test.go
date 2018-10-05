/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package templates

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestTemplatesPackage(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Templates package suite")
}
