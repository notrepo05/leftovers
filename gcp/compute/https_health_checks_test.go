package compute_test

import (
	"errors"

	"github.com/genevievelesperance/leftovers/gcp/compute"
	"github.com/genevievelesperance/leftovers/gcp/compute/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gcpcompute "google.golang.org/api/compute/v1"
)

var _ = Describe("HttpsHealthChecks", func() {
	var (
		client *fakes.HttpsHealthChecksClient
		logger *fakes.Logger
		filter string

		httpsHealthChecks compute.HttpsHealthChecks
	)

	BeforeEach(func() {
		client = &fakes.HttpsHealthChecksClient{}
		logger = &fakes.Logger{}
		filter = "banana"

		httpsHealthChecks = compute.NewHttpsHealthChecks(client, logger)
	})

	Describe("Delete", func() {
		BeforeEach(func() {
			logger.PromptCall.Returns.Proceed = true
			client.ListHttpsHealthChecksCall.Returns.Output = &gcpcompute.HttpsHealthCheckList{
				Items: []*gcpcompute.HttpsHealthCheck{{
					Name: "banana-check",
				}},
			}
		})

		It("deletes https health checks", func() {
			err := httpsHealthChecks.Delete(filter)
			Expect(err).NotTo(HaveOccurred())

			Expect(client.ListHttpsHealthChecksCall.CallCount).To(Equal(1))

			Expect(logger.PromptCall.Receives.Message).To(Equal("Are you sure you want to delete https health check banana-check?"))

			Expect(client.DeleteHttpsHealthCheckCall.CallCount).To(Equal(1))
			Expect(client.DeleteHttpsHealthCheckCall.Receives.HttpsHealthCheck).To(Equal("banana-check"))

			Expect(logger.PrintfCall.Messages).To(Equal([]string{"SUCCESS deleting https health check banana-check\n"}))
		})

		Context("when the client fails to list https health checks", func() {
			BeforeEach(func() {
				client.ListHttpsHealthChecksCall.Returns.Error = errors.New("some error")
			})

			It("returns the error", func() {
				err := httpsHealthChecks.Delete(filter)
				Expect(err).To(MatchError("Listing https health checks: some error"))
			})
		})

		Context("when the health check name does not contain the filter", func() {
			It("does not try to delete it", func() {
				err := httpsHealthChecks.Delete("grape")
				Expect(err).NotTo(HaveOccurred())

				Expect(logger.PromptCall.CallCount).To(Equal(0))
				Expect(client.DeleteHttpsHealthCheckCall.CallCount).To(Equal(0))
			})
		})

		Context("when the client fails to delete the https health check", func() {
			BeforeEach(func() {
				client.DeleteHttpsHealthCheckCall.Returns.Error = errors.New("some error")
			})

			It("logs the error", func() {
				err := httpsHealthChecks.Delete(filter)
				Expect(err).NotTo(HaveOccurred())

				Expect(logger.PrintfCall.Messages).To(Equal([]string{"ERROR deleting https health check banana-check: some error\n"}))
			})
		})

		Context("when the user says no to the prompt", func() {
			BeforeEach(func() {
				logger.PromptCall.Returns.Proceed = false
			})

			It("does not delete the https health check", func() {
				err := httpsHealthChecks.Delete(filter)
				Expect(err).NotTo(HaveOccurred())

				Expect(client.DeleteHttpsHealthCheckCall.CallCount).To(Equal(0))
			})
		})
	})
})
