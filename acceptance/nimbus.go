package acceptance

import (
	"os"

	"github.com/genevieve/leftovers/app"
	"github.com/genevieve/leftovers/nimbus"
	. "github.com/onsi/gomega"
)

type NimbusAcceptance struct {
	User       string
	Client     nimbus.Client
	ResultsDir string
	Logger     *app.Logger
}

func NewNimbusAcceptance() NimbusAcceptance {
	user := os.Getenv("BBL_NIMBUS_USER")
	Expect(user).NotTo(Equal(""), "Missing $BBL_NIMBUS_USER.")

	resultsDir := os.Getenv("LEFTOVERS_NIMBUS_RESULTS_DIR")
	Expect(resultsDir).NotTo(Equal(""), "Missing $LEFTOVERS_NIMBUS_RESULTS_DIR.")

	return NimbusAcceptance{
		User:       user,
		Client:     nimbus.NewClient(user),
		ResultsDir: resultsDir,
		Logger:     app.NewLogger(os.Stdin, os.Stdout, true, false),
	}
}

func (n *NimbusAcceptance) CreateTestbed(name string) {
	// Simple VC-ESX example from https://confluence.eng.vmware.com/display/DevToolsDocKB/4.5+Testbed+Spec+Examples
	simpleTestbedSpec := `str:{"name":"1vc1esx","version":3,"vcs":[{"name":"vc.0","type":"vcva","customBuild":"ob-7363249"}],"esx":[{"name":"esx.0","customBuild":"ob-5969303","style":"fullInstall","iScsi":true,"guestOSlist":[{"vmName":"centos","ovfuri":"NIMBUS_OVF_BASE/CentOS6_x64_2GB/CentOS6_x64_2GB.ovf"}]}],"iscsi":[{"name":"iscsi.0","luns":[100]}]}`
	err := n.Client.CreateTestbed(name, simpleTestbedSpec, n.ResultsDir)
	Expect(err).NotTo(HaveOccurred())
}
