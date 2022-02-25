package errors

type ErrHelpWithMessage struct {
	Message string
	ErrHelp error
}

func (e ErrHelpWithMessage) Error() string {
	return e.Message
}

func (e ErrHelpWithMessage) Unwrap() error {
	return e.ErrHelp
}
