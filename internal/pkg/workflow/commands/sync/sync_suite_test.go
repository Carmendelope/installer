/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package sync

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestSyncPackage(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Sync package suite")
}

