package errors

import "fmt"

type NotFound struct {
	Aggregate string
	ID        any
}

func (e *NotFound) Error() string {
	return fmt.Sprintf("%s with id %v not found", e.Aggregate, e.ID)
}

type InvalidArgument struct {
	Field  string
	Reason string
}

func (e *InvalidArgument) Error() string {
	return fmt.Sprintf("invalid argument %s: %s", e.Field, e.Reason)
}

type InvalidTransition struct {
	Current any
	Target  any
}

func (e *InvalidTransition) Error() string {
	return fmt.Sprintf("invalid transition from %v to %v", e.Current, e.Target)
}
