package specsupdate

// ErrAddingData is an error type for when data cannot be added to the database.
type ErrAddingData struct {
	Err       error
	Structure string
	Msg       string
}

// Error returns the error message.
func (e *ErrAddingData) Error() string {
	if e.Msg == "" {
		e.Msg = "error adding data"
	}
	if e.Structure != "" {
		e.Msg += ": " + e.Structure
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

// ErrNoURL is an error type for when a URL is not provided.
type ErrNoURL struct {
	Err error
	Msg string
}

// Error returns the error message.
func (e *ErrNoURL) Error() string {
	if e.Msg == "" {
		e.Msg = "no url provided- use WithURL()"
	}
	if e.Err != nil {
		e.Msg += ": " + e.Err.Error()
	}
	return e.Msg
}

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
