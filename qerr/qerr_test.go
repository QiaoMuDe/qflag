// Package qerr 错误处理功能测试
// 本文件包含了自定义错误类型的单元测试，测试错误创建、格式化、
// 错误链等功能，确保错误处理机制的正确性和完整性。
package qerr

import (
	"errors"
	"testing"
)

// TestNewValidationError 测试NewValidationError函数
func TestNewValidationError(t *testing.T) {
	msg := "test error message"
	err := NewValidationError(msg)
	if err == nil {
		t.Error("NewValidationError returned nil error")
	}
	expected := ErrValidationFailed.Error() + ": " + msg
	if err.Error() != expected {
		t.Errorf("NewValidationError returned unexpected error message: got %q, want %q", err.Error(), expected)
	}
}

// TestNewValidationErrorf 测试NewValidationErrorf函数
func TestNewValidationErrorf(t *testing.T) {
	format := "test error %d"
	param := 123
	err := NewValidationErrorf(format, param)
	if err == nil {
		t.Error("NewValidationErrorf returned nil error")
	}
	expectedMsg := "test error 123"
	expected := ErrValidationFailed.Error() + ": " + expectedMsg

	if err.Error() != expected {
		t.Errorf("NewValidationErrorf returned unexpected error message: got %q, want %q", err.Error(), expected)
	}
}

// TestJoinErrors 测试JoinErrors函数
func TestJoinErrors(t *testing.T) {
	tests := []struct {
		name     string
		errors   []error
		expected string
	}{{
		name:     "no errors",
		errors:   []error{},
		expected: "",
	}, {
		name:     "single error",
		errors:   []error{NewValidationError("error1")},
		expected: "validation failed: error1",
	}, {
		name: "multiple unique errors",
		errors: []error{
			NewValidationError("error1"),
			NewValidationError("error2"),
		},
		expected: "merged error message:\nA total of 2 unique errors:\n  1. validation failed: error1\n  2. validation failed: error2\n",
	}, {
		name: "multiple duplicate errors",
		errors: []error{
			NewValidationError("error1"),
			NewValidationError("error1"),
		},
		expected: "merged error message:\nA total of 1 unique errors:\n  1. validation failed: error1\n",
	}, {
		name: "with nil errors",
		errors: []error{
			nil,
			NewValidationError("error1"),
			nil,
		},
		expected: "validation failed: error1",
	}, {
		name: "empty error messages",
		errors: []error{
			errors.New(""),
			errors.New(""),
		},
		expected: "merged error message:\nA total of 1 unique errors:\n  1. \n",
	}, {
		name: "mixed error types",
		errors: []error{
			ErrFlagParseFailed,
			NewValidationError("invalid value"),
			ErrEnvLoadFailed,
		},
		expected: "merged error message:\nA total of 3 unique errors:\n  1. parameter parsing error\n  2. validation failed: invalid value\n  3. environment variable loading failed\n",
	}, {
		name: "multiple duplicate errors",
		errors: []error{
			NewValidationError("error1"),
			NewValidationError("error1"),
		},
		expected: "merged error message:\nA total of 1 unique errors:\n  1. validation failed: error1\n",
	}, {
		name: "with nil errors",
		errors: []error{
			nil,
			NewValidationError("error1"),
			nil,
		},
		expected: "validation failed: error1",
	}, {
		name: "empty error messages",
		errors: []error{
			errors.New(""),
			errors.New(""),
		},
		expected: "merged error message:\nA total of 1 unique errors:\n  1. \n",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := JoinErrors(tt.errors)
			if tt.expected == "" {
				if err != nil {
					t.Errorf("JoinErrors(%v) = %v, want nil", tt.errors, err)
				}
				return
			}

			if err == nil {
				t.Fatal("JoinErrors returned nil error when expected non-nil")
			}

			if err.Error() != tt.expected {
				t.Errorf("JoinErrors() error message = %q, want %q", err.Error(), tt.expected)
			}
		})
	}
}
