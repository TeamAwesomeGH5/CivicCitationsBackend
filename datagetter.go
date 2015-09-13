package main

type Getter interface {
	GetCitationByNumber(number uint64) ([]Citation, error)
	GetCitationsByUser(lastname, licensenumber, dob string) ([]Citation, error)
	GetViolationsForCitation(citation Citation) (Violations, error)
	String() string
}
