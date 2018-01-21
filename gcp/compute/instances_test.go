package compute_test

import (
	"errors"

	"github.com/genevievelesperance/leftovers/gcp/compute"
	"github.com/genevievelesperance/leftovers/gcp/compute/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gcpcompute "google.golang.org/api/compute/v1"
)

var _ = Describe("Instances", func() {
	var (
		client *fakes.InstancesClient
		logger *fakes.Logger
		zones  map[string]string
		filter string

		instances compute.Instances
	)

	BeforeEach(func() {
		client = &fakes.InstancesClient{}
		logger = &fakes.Logger{}
		zones = map[string]string{
			"https://zone-1": "zone-1",
		}
		filter = "banana"

		instances = compute.NewInstances(client, logger, zones)
	})

	Describe("Delete", func() {
		BeforeEach(func() {
			logger.PromptCall.Returns.Proceed = true
			client.ListInstancesCall.Returns.Output = &gcpcompute.InstanceList{
				Items: []*gcpcompute.Instance{{
					Name: "banana-instance",
					Zone: "https://zone-1",
				}},
			}
		})

		It("deletes instances", func() {
			err := instances.Delete(filter)
			Expect(err).NotTo(HaveOccurred())

			Expect(client.ListInstancesCall.CallCount).To(Equal(1))
			Expect(client.ListInstancesCall.Receives.Zone).To(Equal("zone-1"))

			Expect(logger.PromptCall.Receives.Message).To(Equal("Are you sure you want to delete instance banana-instance?"))

			Expect(client.DeleteInstanceCall.CallCount).To(Equal(1))
			Expect(client.DeleteInstanceCall.Receives.Zone).To(Equal("zone-1"))
			Expect(client.DeleteInstanceCall.Receives.Instance).To(Equal("banana-instance"))

			Expect(logger.PrintfCall.Messages).To(Equal([]string{"SUCCESS deleting instance banana-instance\n"}))
		})

		Context("when the instance has tags", func() {
			BeforeEach(func() {
				client.ListInstancesCall.Returns.Output = &gcpcompute.InstanceList{
					Items: []*gcpcompute.Instance{{
						Name: "banana-instance",
						Tags: &gcpcompute.Tags{Items: []string{"banana-director"}},
					}},
				}
			})

			It("uses them in the prompt", func() {
				err := instances.Delete(filter)
				Expect(err).NotTo(HaveOccurred())

				Expect(logger.PromptCall.Receives.Message).To(Equal("Are you sure you want to delete instance banana-instance (banana-director)?"))
			})
		})

		Context("when the client fails to list instances", func() {
			BeforeEach(func() {
				client.ListInstancesCall.Returns.Error = errors.New("some error")
			})

			It("returns the error", func() {
				err := instances.Delete(filter)
				Expect(err).To(MatchError("Listing instances for zone zone-1: some error"))
			})
		})

		Context("when the clearer name for the instance group does not contain the filter", func() {
			It("does not delete the instance", func() {
				err := instances.Delete("grape")
				Expect(err).NotTo(HaveOccurred())

				Expect(logger.PromptCall.CallCount).To(Equal(0))
				Expect(client.DeleteInstanceCall.CallCount).To(Equal(0))
			})
		})

		Context("when the client fails to delete the instance", func() {
			BeforeEach(func() {
				client.DeleteInstanceCall.Returns.Error = errors.New("some error")
			})

			It("logs the error", func() {
				err := instances.Delete(filter)
				Expect(err).NotTo(HaveOccurred())

				Expect(logger.PrintfCall.Messages).To(Equal([]string{"ERROR deleting instance banana-instance: some error\n"}))
			})
		})

		Context("when the user says no to the prompt", func() {
			BeforeEach(func() {
				logger.PromptCall.Returns.Proceed = false
			})

			It("does not delete the instance", func() {
				err := instances.Delete(filter)
				Expect(err).NotTo(HaveOccurred())

				Expect(client.DeleteInstanceCall.CallCount).To(Equal(0))
			})
		})
	})
})
