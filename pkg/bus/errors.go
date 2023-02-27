package bus

// ErrNoDatabase is an error type for when a database is not provided.
type ErrNoDatabase struct {
	Err error
	Msg string
}

// Error returns the error message.
func (e *ErrNoDatabase) Error() string {
	if e.Msg == "" {
		e.Msg = "no database provided- use WithDatabase()"
	}
	if e.Err != nil {
		e.Msg += ": " + e.Err.Error()
	}
	return e.Msg
}

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
