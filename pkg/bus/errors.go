package bus

// ErrSpecsNotSet is an error type for when Specs are not set.
type ErrSpecsNotSet struct {
	Err error
	Msg string
}

// Error returns the error message.
func (e *ErrSpecsNotSet) Error() string {
	if e.Msg == "" {
		e.Msg = "specs not set"
	}
	if e.Err != nil {
		e.Msg += ": " + e.Err.Error()
	}
	return e.Msg
}
