package pii_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"mmm-osint/internal/pkg/pii"
)

var _ = Describe("Interface", func() {
	Describe("Result", func() {
		var result *pii.Result

		BeforeEach(func() {
			result = &pii.Result{
				Total: 3,
				Entities: []pii.Entity{
					{
						Type:     pii.PIITypeEmail,
						Value:    "test@example.com",
						Count:    1,
						Contexts: []string{"Contact test@example.com"},
					},
					{
						Type:     pii.PIITypePhone,
						Value:    "(555) 123-4567",
						Count:    1,
						Contexts: []string{"Call (555) 123-4567"},
					},
					{
						Type:     pii.PIITypeEmail,
						Value:    "admin@test.org",
						Count:    1,
						Contexts: []string{"Email admin@test.org"},
					},
				},
				Stats: map[string]int{
					"email": 2,
					"phone": 1,
				},
			}
		})

		It("should report correct emptiness", func() {
			Expect(result.IsEmpty()).To(BeFalse())

			emptyResult := &pii.Result{Total: 0}
			Expect(emptyResult.IsEmpty()).To(BeTrue())
		})

		It("should check type existence correctly", func() {
			Expect(result.HasType(pii.PIITypeEmail)).To(BeTrue())
			Expect(result.HasType(pii.PIITypePhone)).To(BeTrue())
			Expect(result.HasType(pii.PIITypeCreditCard)).To(BeFalse())
		})

		It("should filter entities by type", func() {
			emails := result.GetByType(pii.PIITypeEmail)
			Expect(emails).To(HaveLen(2))
			
			phones := result.GetByType(pii.PIITypePhone)
			Expect(phones).To(HaveLen(1))
			
			creditCards := result.GetByType(pii.PIITypeCreditCard)
			Expect(creditCards).To(BeEmpty())
		})

		It("should provide convenience methods", func() {
			emails := result.GetEmails()
			Expect(emails).To(HaveLen(2))
			
			phones := result.GetPhones()
			Expect(phones).To(HaveLen(1))
			
			creditCards := result.GetCreditCards()
			Expect(creditCards).To(BeEmpty())
			
			ips := result.GetIPAddresses()
			Expect(ips).To(BeEmpty())
		})
	})

	Describe("Entity", func() {
		It("should have proper structure", func() {
			entity := pii.Entity{
				Type:     pii.PIITypeEmail,
				Value:    "test@example.com",
				Count:    2,
				Contexts: []string{"Email: test@example.com", "Contact test@example.com"},
			}

			Expect(entity.Type).To(Equal(pii.PIITypeEmail))
			Expect(entity.Value).To(Equal("test@example.com"))
			Expect(entity.Count).To(Equal(2))
			Expect(entity.Contexts).To(HaveLen(2))
		})
	})

	Describe("PIIType constants", func() {
		It("should have all expected types", func() {
			Expect(pii.PIITypeEmail).To(Equal(pii.PIIType("email")))
			Expect(pii.PIITypePhone).To(Equal(pii.PIIType("phone")))
			Expect(pii.PIITypeCreditCard).To(Equal(pii.PIIType("credit_card")))
			Expect(pii.PIITypeSSN).To(Equal(pii.PIIType("ssn")))
			Expect(pii.PIITypeIPAddress).To(Equal(pii.PIIType("ip_address")))
			Expect(pii.PIITypeAddress).To(Equal(pii.PIIType("address")))
			Expect(pii.PIITypeName).To(Equal(pii.PIIType("name")))
			Expect(pii.PIITypeBitcoin).To(Equal(pii.PIIType("bitcoin")))
			Expect(pii.PIITypeIBAN).To(Equal(pii.PIIType("iban")))
			Expect(pii.PIITypeOther).To(Equal(pii.PIIType("other")))
		})
	})
})