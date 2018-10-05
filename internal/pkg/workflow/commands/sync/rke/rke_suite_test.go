/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package rke

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestRKEPackage(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "RKE package suite")
}
