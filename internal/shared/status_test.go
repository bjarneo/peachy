package shared

import "testing"

func TestNewStatus(t *testing.T) {
	s := NewStatus(StatusSuccess, "test message")
	if s.Type != StatusSuccess {
		t.Errorf("NewStatus Type = %v, want %v", s.Type, StatusSuccess)
	}
	if s.Message != "test message" {
		t.Errorf("NewStatus Message = %q, want %q", s.Message, "test message")
	}
}

func TestStatusConstructors(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(string) Status
		wantType StatusType
	}{
		{"Info", Info, StatusInfo},
		{"Success", Success, StatusSuccess},
		{"Error", Error, StatusError},
		{"Warning", Warning, StatusWarning},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.fn("test")
			if s.Type != tt.wantType {
				t.Errorf("%s() Type = %v, want %v", tt.name, s.Type, tt.wantType)
			}
		})
	}
}

func TestStatusIsError(t *testing.T) {
	if !Error("err").IsError() {
		t.Error("Error().IsError() = false, want true")
	}
	if Info("info").IsError() {
		t.Error("Info().IsError() = true, want false")
	}
}

func TestStatusIsSuccess(t *testing.T) {
	if !Success("ok").IsSuccess() {
		t.Error("Success().IsSuccess() = false, want true")
	}
	if Info("info").IsSuccess() {
		t.Error("Info().IsSuccess() = true, want false")
	}
}

func TestStatusString(t *testing.T) {
	s := Info("hello world")
	if s.String() != "hello world" {
		t.Errorf("String() = %q, want %q", s.String(), "hello world")
	}
}
