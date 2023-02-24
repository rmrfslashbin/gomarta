package gtfspec

// ErrCreatingFile is an error type for when a file cannot be created.
type ErrCreatingFile struct {
	Err      error
	Msg      string
	Filename string
}

// Error returns the error message.
func (e *ErrCreatingFile) Error() string {
	if e.Msg == "" {
		e.Msg = "error creating file"
	}
	if e.Filename != "" {
		e.Msg += ": " + e.Filename
	}
	if e.Err != nil {
		e.Msg += ": " + e.Err.Error()
	}
	return e.Msg
}

// ErrCSVReader is an error type for when a CSV file cannot be read.
type ErrCSVReader struct {
	Err error
	Msg string
}

// Error returns the error message.
func (e *ErrCSVReader) Error() string {
	if e.Msg == "" {
		e.Msg = "error reading csv"
	}
	if e.Err != nil {
		e.Msg += ": " + e.Err.Error()
	}
	return e.Msg
}

// ErrFetchingURL is an error type for when a URL cannot be fetched.
type ErrFetchingURL struct {
	Err error
	Msg string
}

// Error returns the error message.
func (e *ErrFetchingURL) Error() string {
	if e.Msg == "" {
		e.Msg = "error fetching url"
	}
	if e.Err != nil {
		e.Msg += ": " + e.Err.Error()
	}
	return e.Msg
}

// ErrMarshallingJSON is an error type for when JSON cannot be marshalled.
type ErrMarshallingJSON struct {
	Err error
	Msg string
}

// Error returns the error message.
func (e *ErrMarshallingJSON) Error() string {
	if e.Msg == "" {
		e.Msg = "error marshalling json"
	}
	if e.Err != nil {
		e.Msg += ": " + e.Err.Error()
	}
	return e.Msg
}

// ErrMarshallingGOB is an error type for when GOB cannot be marshalled.
type ErrMarshallingGOB struct {
	Err error
	Msg string
}

// Error returns the error message.
func (e *ErrMarshallingGOB) Error() string {
	if e.Msg == "" {
		e.Msg = "error marshalling gob"
	}
	if e.Err != nil {
		e.Msg += ": " + e.Err.Error()
	}
	return e.Msg
}

// ErrOpeningFile is an error type for when a file cannot be opened.
type ErrOpeningFile struct {
	Err      error
	Msg      string
	Filename string
}

// Error returns the error message.
func (e *ErrOpeningFile) Error() string {
	if e.Msg == "" {
		e.Msg = "error opening file"
	}
	if e.Filename != "" {
		e.Msg += ": " + e.Filename
	}
	if e.Err != nil {
		e.Msg += ": " + e.Err.Error()
	}
	return e.Msg
}

// ErrParsingFile is an error type for when a file cannot be parsed.
type ErrParsingFile struct {
	Err  error
	File string
	Msg  string
}

// Error returns the error message.
func (e *ErrParsingFile) Error() string {
	if e.Msg == "" {
		e.Msg = "error parsing file"
	}
	if e.File != "" {
		e.Msg += ": " + e.File
	}
	if e.Err != nil {
		e.Msg += ": " + e.Err.Error()
	}
	return e.Msg
}

// ErrReadingUrlBody is an error type for when a URL body cannot be read.
type ErrReadingUrlBody struct {
	Err error
	Msg string
}

// Error returns the error message.
func (e *ErrReadingUrlBody) Error() string {
	if e.Msg == "" {
		e.Msg = "error reading url body"
	}
	if e.Err != nil {
		e.Msg += ": " + e.Err.Error()
	}
	return e.Msg
}

// ErrUnmarshallingGOB is an error type for when GOB cannot be unmarshalled.
type ErrUnmarshallingGOB struct {
	Err error
	Msg string
}

// Error returns the error message.
func (e *ErrUnmarshallingGOB) Error() string {
	if e.Msg == "" {
		e.Msg = "error unmarshalling gob"
	}
	if e.Err != nil {
		e.Msg += ": " + e.Err.Error()
	}
	return e.Msg
}

// ErrWritingFile is an error type for when a file cannot be read.
type ErrWritingFile struct {
	Err      error
	Msg      string
	Filename string
}

// Error returns the error message.
func (e *ErrWritingFile) Error() string {
	if e.Msg == "" {
		e.Msg = "error writing file"
	}
	if e.Filename != "" {
		e.Msg += ": " + e.Filename
	}
	if e.Err != nil {
		e.Msg += ": " + e.Err.Error()
	}
	return e.Msg
}

// ErrZipFileReader is an error type for when a zip file cannot be read.
type ErrZipFileReader struct {
	Err error
	Msg string
}

// Error returns the error message.
func (e *ErrZipFileReader) Error() string {
	if e.Msg == "" {
		e.Msg = "error reading zip file"
	}
	if e.Err != nil {
		e.Msg += ": " + e.Err.Error()
	}
	return e.Msg
}

// ErrZipReader is an error type for when a zip file cannot be read.
type ErrZipReader struct {
	Err error
	Msg string
}

// Error returns the error message.
func (e *ErrZipReader) Error() string {
	if e.Msg == "" {
		e.Msg = "error reading zip file"
	}
	if e.Err != nil {
		e.Msg += ": " + e.Err.Error()
	}
	return e.Msg
}
