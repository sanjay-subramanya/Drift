package model

import "fmt"

type Finding struct {
	Drift 	 Drift
	Severity Severity
	Message  string
}

func (f Finding) String() string {
	return fmt.Sprintf("[%s], %s", f.Severity.String(), f.Message)
}