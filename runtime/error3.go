package runtime

import (
	"context"
	"fmt"
	"log"
)

// ErrorHandleFn3 - function type for error handling
type ErrorHandleFn3 func(location string, errs ...error) *Status

// ErrorHandleWithContextFn3 - function type for error handling with context
type ErrorHandleWithContextFn3 func(ctx context.Context, location string, errs ...error) *Status

// ErrorStatusHandleFn3 - function type for error status handling
type ErrorStatusHandleFn3 func(s *Status) *Status

// ErrorHandler3 - template parameter error handler interface
type ErrorHandler3 interface {
	Handle(location string, errs ...error) *Status
	HandleWithContext(ctx context.Context, location string, errs ...error) *Status
	HandleStatus(s *Status) *Status
}

// NoOpError3 - no operation on errors
type NoOpError3 struct{}

func (NoOpError3) Handle(location string, errs ...error) *Status {
	if !IsErrors(errs) {
		return NewStatusOK()
	}
	return NewStatusError(location, errs...)
}

func (NoOpError3) HandleWithContext(ctx context.Context, location string, errs ...error) *Status {
	if !IsErrors(errs) {
		return NewStatusOK()
	}
	return NewStatusError(location, errs...).SetContext(ctx)
}

func (NoOpError3) HandleStatus(s *Status) *Status {
	return s
}

// DebugError3 - debug error handler
type DebugError3 struct{}

func (h DebugError3) Handle(location string, errs ...error) *Status {
	if !IsErrors(errs) {
		return NewStatusOK()
	}
	return h.HandleStatus(NewStatus(StatusInternal, location, errs...))
}

func (h DebugError3) HandleWithContext(ctx context.Context, location string, errs ...error) *Status {
	if !IsErrors(errs) {
		return NewStatusOK()
	}
	return h.HandleStatus(NewStatus(StatusInternal, location, errs...).SetContext(ctx))
}

func (h DebugError3) HandleStatus(s *Status) *Status {
	if s != nil && s.IsErrors() {
		loc := IfElse[string](s.Location() == "", "[]", s.Location())
		req := IfElse[string](s.RequestId() == "", "[]", s.RequestId())
		fmt.Printf("[%v %v %v]\n", req, loc, s.Errors())
		s.RemoveErrors()
	}
	return s
}

// LogError3 - logging error handler
type LogError3 struct{}

func (h LogError3) Handle(location string, errs ...error) *Status {
	if !IsErrors(errs) {
		return NewStatusOK()
	}
	return h.HandleStatus(NewStatus(StatusInternal, location, errs...))
}

func (h LogError3) HandleWithContext(ctx context.Context, location string, errs ...error) *Status {
	if !IsErrors(errs) {
		return NewStatusOK()
	}
	return h.HandleStatus(NewStatus(StatusInternal, location, errs...).SetContext(ctx))
}

func (h LogError3) HandleStatus(s *Status) *Status {
	if s != nil && s.IsErrors() {
		loc := IfElse[string](s.Location() == "", "[]", s.Location())
		req := IfElse[string](s.RequestId() == "", "[]", s.RequestId())
		log.Println(req, loc, s.Errors())
		s.RemoveErrors()
	}
	return s
}

// Handle3 - templated function providing an error handle function via a closure
func Handle3[E ErrorHandler3]() ErrorHandleFn3 {
	var e E
	return func(location string, errs ...error) *Status {
		return e.Handle(location, errs...)
	}
}

// HandleWithContext3 - templated function providing an error handle function with context via a closure
func HandleWithContext3[E ErrorHandler3]() ErrorHandleWithContextFn3 {
	var e E
	return func(ctx context.Context, location string, errs ...error) *Status {
		return e.HandleWithContext(ctx, location, errs...)
	}
}

// StatusHandle3 - templated function providing an error status handle function via a closure
func StatusHandle3[E ErrorHandler3]() ErrorStatusHandleFn3 {
	var e E
	return func(s *Status) *Status {
		return e.HandleStatus(s)
	}
}