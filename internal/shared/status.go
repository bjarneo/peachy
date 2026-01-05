package shared

// StatusType represents the type of status message
type StatusType int

const (
	StatusInfo StatusType = iota
	StatusSuccess
	StatusError
	StatusWarning
)

// Status represents a status message with a type
type Status struct {
	Type    StatusType
	Message string
}

// NewStatus creates a new status message
func NewStatus(t StatusType, msg string) Status {
	return Status{Type: t, Message: msg}
}

// Info creates an info status
func Info(msg string) Status {
	return NewStatus(StatusInfo, msg)
}

// Success creates a success status
func Success(msg string) Status {
	return NewStatus(StatusSuccess, msg)
}

// Error creates an error status
func Error(msg string) Status {
	return NewStatus(StatusError, msg)
}

// Warning creates a warning status
func Warning(msg string) Status {
	return NewStatus(StatusWarning, msg)
}

// IsError returns true if this is an error status
func (s Status) IsError() bool {
	return s.Type == StatusError
}

// IsSuccess returns true if this is a success status
func (s Status) IsSuccess() bool {
	return s.Type == StatusSuccess
}

// String returns the message
func (s Status) String() string {
	return s.Message
}
