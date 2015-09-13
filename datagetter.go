package main

type Getter interface {
	GetCitationByNumber(number uint64) ([]Citation, error)
	GetViolationsForCitation(citation Citation) (Violations, error)
}
