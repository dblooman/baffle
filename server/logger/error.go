package logger

import "runtime"

type errorWrapper struct {
	Message    string
	Context    Fields
	Frame      *errorFrame
	InnerError error
}

type errorFrame struct {
	Function string
	File     string
	Line     int
}

func (err *errorWrapper) Error() string {
	return err.Message
}

var _ Fielder = (*errorWrapper)(nil)

// ToFields builds a Fields map containing all the information in the error
// If a piece of information is missing the field won't be present
func (err *errorWrapper) ToFields() Fields {
	fields := Fields{"message": err.Message}

	if err.Context != nil {
		fields["context"] = err.Context
	}

	if err.InnerError != nil {
		fields["innerError"] = errorToField(err.InnerError)
	}

	if frame := err.Frame; frame != nil {
		if frame.Function != "" {
			fields["function"] = frame.Function
		}
		fields["file"] = frame.File
		fields["line"] = frame.Line
	}

	return fields
}

func WrapError(inner error, message string) error {
	return &errorWrapper{
		Message:    message,
		Frame:      getFrame(2),
		InnerError: inner,
	}
}

func WrapErrorWithContext(inner error, message string, context Fields) error {
	return &errorWrapper{
		Message:    message,
		Context:    context,
		Frame:      getFrame(2),
		InnerError: inner,
	}
}

func WrapErrorWithContextAndStack(inner error, message string, context Fields, stackSkip int) error {
	return &errorWrapper{
		Message:    message,
		Context:    context,
		Frame:      getFrame(2 + stackSkip),
		InnerError: inner,
	}
}

func getFrame(skip int) *errorFrame {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return nil
	}

	var funcName string
	if f := runtime.FuncForPC(pc); f != nil {
		funcName = f.Name()
	}

	return &errorFrame{
		Function: funcName,
		File:     file,
		Line:     line,
	}
}
