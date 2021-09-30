package nimbus_test

import (
	"errors"

	"github.com/genevieve/leftovers/nimbus"
	"github.com/genevieve/leftovers/nimbus/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Testbeds", func() {
	var (
		testbeds nimbus.Testbeds
		logger   *fakes.Logger
		filter   string

		client *fakes.TestbedsClient
	)

	BeforeEach(func() {
		logger = &fakes.Logger{}
		client = &fakes.TestbedsClient{}
		filter = ""

		testbeds = nimbus.NewTestbeds(client, logger)
	})

	Describe("List", func() {
		BeforeEach(func() {
			logger.PromptWithDetailsCall.Returns.Bool = true
			client.ListTestbedsCall.Returns.TestbedSlice = []nimbus.Testbed{
				nimbus.NewTestbed(nil, "some-name"),
				nimbus.NewTestbed(nil, "other-name"),
			}
		})

		It("returns all the deletables", func() {
			list, err := testbeds.List(filter, false)
			Expect(err).NotTo(HaveOccurred())

			Expect(list).To(HaveLen(2))
			Expect(list[0].Name()).To(Equal("some-name"))
			Expect(list[1].Name()).To(Equal("other-name"))
		})

		Context("when the resource does not contain the filter", func() {
			BeforeEach(func() {
				client.ListTestbedsCall.Returns.TestbedSlice = []nimbus.Testbed{
					nimbus.NewTestbed(nil, "other-name"),
				}
			})

			It("returns all the deletables", func() {
				list, err := testbeds.List("kiwi", false)
				Expect(err).NotTo(HaveOccurred())

				Expect(list).To(HaveLen(0))
			})
		})

		Context("when listing fails", func() {
			BeforeEach(func() {
				client.ListTestbedsCall.Returns.Error = errors.New("error-description")
			})

			It("returns an error", func() {
				_, err := testbeds.List(filter, false)
				Expect(err).To(MatchError("List Nimbus Testbeds: error-description"))
			})
		})
	})

	Describe("Type", func() {
		It("returns the type of the resource", func() {
			Expect(testbeds.Type()).To(Equal("Testbed"))
		})
	})
})
