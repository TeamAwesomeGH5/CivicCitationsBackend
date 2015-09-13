package main

import "errors"

var NoViolationError = errors.New("There were no violations found")

type Violation struct {
	Id                    int
	Citation_number       uint64
	Violation_number      string
	Violation_description string
	Warrant_status        string
	Warrant_number        string
	Status                string
	Status_date           string
	Fine_amount           string
	Court_cost            string
}
type Violations []Violation

func NewViolation() Violation {
	return Violation{
		Id:                    0,
		Citation_number:       0,
		Violation_number:      "",
		Violation_description: "",
		Warrant_status:        "",
		Warrant_number:        "",
		Status:                "",
		Status_date:           "",
		Fine_amount:           "",
		Court_cost:            "",
	}
}
