package zerotier

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"os"
	"testing"
)

var itKubeConfigFile = os.Getenv("IT_K8S_KUBECONFIG")


func TestZerotierPackage(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "ZeroTier package suite")
}

