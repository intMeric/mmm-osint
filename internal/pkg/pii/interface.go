package pii

import "context"

type Extractor interface {
	ExtractPII(ctx context.Context, text string) (*Result, error)
	Close() error
}

type Result struct {
	Total    int            `json:"total"`
	Entities []Entity       `json:"entities"`
	Stats    map[string]int `json:"stats"`
}

type Entity struct {
	Type     PIIType `json:"type"`
	Value    string  `json:"value"`
	Count    int     `json:"count"`
	Contexts []string `json:"contexts"`
}

type PIIType string

const (
	PIITypeEmail       PIIType = "email"
	PIITypePhone       PIIType = "phone"
	PIITypeCreditCard  PIIType = "credit_card"
	PIITypeSSN         PIIType = "ssn"
	PIITypeIPAddress   PIIType = "ip_address"
	PIITypeAddress     PIIType = "address"
	PIITypeName        PIIType = "name"
	PIITypeBitcoin     PIIType = "bitcoin"
	PIITypeIBAN        PIIType = "iban"
	PIITypeOther       PIIType = "other"
)

func (r *Result) IsEmpty() bool {
	return r.Total == 0
}

func (r *Result) HasType(piiType PIIType) bool {
	for _, entity := range r.Entities {
		if entity.Type == piiType {
			return true
		}
	}
	return false
}

func (r *Result) GetByType(piiType PIIType) []Entity {
	var entities []Entity
	for _, entity := range r.Entities {
		if entity.Type == piiType {
			entities = append(entities, entity)
		}
	}
	return entities
}

func (r *Result) GetEmails() []Entity {
	return r.GetByType(PIITypeEmail)
}

func (r *Result) GetPhones() []Entity {
	return r.GetByType(PIITypePhone)
}

func (r *Result) GetCreditCards() []Entity {
	return r.GetByType(PIITypeCreditCard)
}

func (r *Result) GetIPAddresses() []Entity {
	return r.GetByType(PIITypeIPAddress)
}