package k8s_test

import (
	"os"

	. "github.com/ZacharyChang/kcui/k8s"
	"github.com/ZacharyChang/kcui/pkg/option"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("K8s", func() {
	var (
		testClient *Client
		opts       *option.Options
	)
	BeforeEach(func() {
		opts = option.NewOptions()
		opts.Kubeconfig = os.Getenv("KUBECONFIG")
		testClient = NewClient(opts)
		testClient.Namespace = "test"
	})

	Describe("Test kubernetes client", func() {
		Context("Set namespace to test", func() {
			It("should not be nil", func() {
				Expect(testClient).ShouldNot(BeNil())
			})
			It("namespace should be test", func() {
				Expect(testClient.Namespace).To(Equal("test"))
			})
		})
	})

})
