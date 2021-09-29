package acceptance

import (
	"bytes"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/genevieve/leftovers/app"
	"github.com/genevieve/leftovers/nimbus"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("nimbus", func() {
	var (
		acc NimbusAcceptance

		stdout  *bytes.Buffer
		filter  string
		deleter nimbus.Leftovers
	)

	BeforeEach(func() {
		iaas := os.Getenv("LEFTOVERS_ACCEPTANCE")
		if strings.ToLower(iaas) != "nimbus" {
			Skip("Skipping nimbus acceptance tests.")
		}
		filter = os.Getenv("LEFTOVERS_NIMBUS_FILTER")
		if filter == "" {
			filter = "leftovers-acceptance"
		}

		acc = NewNimbusAcceptance()

		noConfirm := true
		debug := true
		stdout = bytes.NewBuffer([]byte{})
		logger := app.NewLogger(stdout, os.Stdin, noConfirm, debug)

		var err error
		deleter, err = nimbus.NewLeftovers(logger, acc.User)
		Expect(err).NotTo(HaveOccurred())

		color.NoColor = true
	})

	Describe("leftovers", func() {
		BeforeEach(func() {
			acc.CreateTestbed("leftovers-acceptance")
		})

		It("can list and delete resources with the filter", func() {
			By("listing resources first", func() {
				deleter.List(filter, false)

				Expect(stdout.String()).To(ContainSubstring("[Nimbus Testbed: leftovers-acceptance]"))
				Expect(stdout.String()).NotTo(ContainSubstring("[Nimbus Testbed: leftovers-acceptance] Deleting..."))
			})

			By("successfully deleting Testbeds", func() {
				err := deleter.Delete(filter, false)
				Expect(err).NotTo(HaveOccurred())

				Expect(stdout.String()).To(ContainSubstring("[Nimbus Testbed: leftovers-acceptance] Deleting..."))
				Expect(stdout.String()).To(ContainSubstring("[Nimbus Testbed: leftovers-acceptance] Deleted!"))
			})
		})
	})
})
