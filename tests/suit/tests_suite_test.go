package suit_test

import (
	"flag"
	"fmt"
	"github.com/konveyor/forklift-controller/tests/suit/framework"
	"github.com/konveyor/forklift-controller/tests/suit/reporters"
	"github.com/konveyor/forklift-controller/tests/suit/utils"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

var (
	kubectlPath       = flag.String("kubectl-path", "/home/eslutsky/bin/kubectl", "The path to the kubectl binary")
	ocPath            = flag.String("oc-path", "oc", "The path to the oc binary")
	forkliftInstallNs = flag.String("forklift-namespace", "konveyor-forklift", "The namespace of the CDI controller")
	kubeConfig        = flag.String("kubeconfig2", "/tmp/kubeconfig", "The absolute path to the kubeconfig file")
	kubeURL           = flag.String("kubeurl", "", "kube URL url:port")
	goCLIPath         = flag.String("gocli-path", "cli.sh", "The path to cli script")
	dockerPrefix      = flag.String("docker-prefix", "", "The docker host:port")
	dockerTag         = flag.String("docker-tag", "", "The docker tag")
)

// forkliftFailHandler call ginkgo.Fail with printing the additional information
func forkliftFailHandler(message string, callerSkip ...int) {
	if len(callerSkip) > 0 {
		callerSkip[0]++
	}
	ginkgo.Fail(message, callerSkip...)
}

func TestTests(t *testing.T) {
	defer GinkgoRecover()
	RegisterFailHandler(forkliftFailHandler)
	BuildTestSuite()
	RunSpecsWithDefaultAndCustomReporters(t, "Tests Suite", reporters.NewReporters())
}

// To understand the order in which things are run, read http://onsi.github.io/ginkgo/#understanding-ginkgos-lifecycle
// flag parsing happens AFTER ginkgo has constructed the entire testing tree. So anything that uses information from flags
// cannot work when called during test tree construction.
func BuildTestSuite() {
	BeforeSuite(func() {
		fmt.Fprintf(ginkgo.GinkgoWriter, "Reading parameters\n")
		// Read flags, and configure client instances
		framework.ClientsInstance.KubectlPath = *kubectlPath
		framework.ClientsInstance.OcPath = *ocPath
		framework.ClientsInstance.ForkliftInstallNs = *forkliftInstallNs
		framework.ClientsInstance.KubeConfig = *kubeConfig
		framework.ClientsInstance.KubeURL = *kubeURL
		framework.ClientsInstance.GoCLIPath = *goCLIPath
		framework.ClientsInstance.DockerPrefix = *dockerPrefix
		framework.ClientsInstance.DockerTag = *dockerTag

		fmt.Fprintf(ginkgo.GinkgoWriter, "Kubectl path: %s\n", framework.ClientsInstance.KubectlPath)
		fmt.Fprintf(ginkgo.GinkgoWriter, "OC path: %s\n", framework.ClientsInstance.OcPath)
		fmt.Fprintf(ginkgo.GinkgoWriter, "Forklift install NS: %s\n", framework.ClientsInstance.ForkliftInstallNs)
		fmt.Fprintf(ginkgo.GinkgoWriter, "Kubeconfig: %s\n", framework.ClientsInstance.KubeConfig)
		fmt.Fprintf(ginkgo.GinkgoWriter, "KubeURL: %s\n", framework.ClientsInstance.KubeURL)
		fmt.Fprintf(ginkgo.GinkgoWriter, "GO CLI path: %s\n", framework.ClientsInstance.GoCLIPath)
		fmt.Fprintf(ginkgo.GinkgoWriter, "DockerPrefix: %s\n", framework.ClientsInstance.DockerPrefix)
		fmt.Fprintf(ginkgo.GinkgoWriter, "DockerTag: %s\n", framework.ClientsInstance.DockerTag)

		restConfig, err := framework.ClientsInstance.LoadConfig()
		if err != nil {
			// Can't use Expect here due this being called outside of an It block, and Expect
			// requires any calls to it to be inside an It block.
			ginkgo.Fail("ERROR, unable to load RestConfig")
		}
		framework.ClientsInstance.RestConfig = restConfig
		// clients
		kcs, err := framework.ClientsInstance.GetKubeClient()
		if err != nil {
			ginkgo.Fail(fmt.Sprintf("ERROR, unable to create K8SClient: %v", err))
		}
		framework.ClientsInstance.K8sClient = kcs

		crClient, err := framework.ClientsInstance.GetCrClient()
		if err != nil {
			ginkgo.Fail(fmt.Sprintf("ERROR, unable to create CrClient: %v", err))
		}
		framework.ClientsInstance.CrClient = crClient

		dyn, err := framework.ClientsInstance.GetDynamicClient()
		if err != nil {
			ginkgo.Fail(fmt.Sprintf("ERROR, unable to create DynamicClient: %v", err))
		}
		framework.ClientsInstance.DynamicClient = dyn

		utils.CacheTestsData(framework.ClientsInstance.K8sClient, framework.ClientsInstance.ForkliftInstallNs)
	})

}
