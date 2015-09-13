package main

import "errors"

var NoViolationError = errors.New("There were no violations found")

type Violation struct{}
type Violations []Violation
