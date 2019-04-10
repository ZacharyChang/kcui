package k8s_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/ZacharyChang/kcui/k8s"
)

var _ = Describe("K8s", func() {
	var (
		testClient *Client
	)
	BeforeEach(func() {
		testClient = NewClient(filepath.Join(os.Getenv("HOME"), ".kube", "config"))
		testClient.SetNamespace("test")
	})

	Describe("Test kubernetes client", func() {
		Context("Set namespace to test", func() {
			It("should not be nil", func() {
				Expect(testClient).ShouldNot(BeNil())
			})
			It("namespace should be test", func() {
				Expect(testClient.GetNamespace()).To(Equal("test"))
			})
		})
	})

})
