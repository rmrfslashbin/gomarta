package database

// ErrMySqlOpen is returned when there is an error opening the mysql database
type ErrMySqlOpen struct {
	Err error
	Dsn string
	Msg string
}

// Error returns the error message.
func (e *ErrMySqlOpen) Error() string {
	if e.Msg == "" {
		e.Msg = "error opening mysql database"
	}
	if e.Dsn != "" {
		e.Msg += ": " + e.Dsn
	}
	if e.Err != nil {
		e.Msg += ": " + e.Err.Error()
	}
	return e.Msg
}

// ErrNoDatabase is returned when no database is specified
type ErrNoDatabase struct {
	Err error
	Msg string
}

// Error returns the error message.
func (e *ErrNoDatabase) Error() string {
	if e.Msg == "" {
		e.Msg = "no databsae specified"
	}
	if e.Err != nil {
		e.Msg += ": " + e.Err.Error()
	}
	return e.Msg
}

// ErrPGSqlOpen is returned when there is an error opening the pgsql database
type ErrPGSqlOpen struct {
	Err error
	Dsn string
	Msg string
}

// Error returns the error message.
func (e *ErrPGSqlOpen) Error() string {
	if e.Msg == "" {
		e.Msg = "error opening pgsql database"
	}
	if e.Dsn != "" {
		e.Msg += ": " + e.Dsn
	}
	if e.Err != nil {
		e.Msg += ": " + e.Err.Error()
	}
	return e.Msg
}

// ErrSqliteOpen is returned when there is an error opening the sqlite database
type ErrSqliteOpen struct {
	Err      error
	Filename string
	Msg      string
}

// Error returns the error message.
func (e *ErrSqliteOpen) Error() string {
	if e.Msg == "" {
		e.Msg = "error opening sqlite database"
	}
	if e.Filename != "" {
		e.Msg += ": " + e.Filename
	}
	if e.Err != nil {
		e.Msg += ": " + e.Err.Error()
	}
	return e.Msg
}
