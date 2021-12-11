package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{ID: "123", Name: "name", Age: 99, Email: "email", Role: "custom", Phones: []string{"+123456"}, meta: nil},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   errors.New("len should be 36"),
				},
				ValidationError{
					Field: "Age",
					Err:   errors.New("should be less than 50"),
				},
				ValidationError{
					Field: "Email",
					Err:   errors.New("should fit ^\\w+@\\w+\\.\\w+$"),
				},
				ValidationError{
					Field: "Role",
					Err:   errors.New("should be in admin,stuff"),
				},
				ValidationError{
					Field: "Phones",
					Err:   errors.New("len should be 11"),
				},
			},
		},
		{
			in: User{
				ID:     "100000000020000000003000000000123456",
				Name:   "name",
				Age:    20,
				Email:  "email@g.com",
				Role:   "stuff",
				Phones: []string{"01234567890"},
			},
			expectedErr: ValidationErrors{},
		},
		{
			in: App{Version: "1.3456"},
			expectedErr: ValidationErrors{ValidationError{
				Field: "Version",
				Err:   errors.New("len should be 5"),
			}},
		},
		{
			in:          App{Version: "12345"},
			expectedErr: ValidationErrors{},
		},
		{
			in:          Token{Header: []byte("head"), Payload: []byte("pay"), Signature: []byte("sign")},
			expectedErr: ValidationErrors{},
		},
		{
			in: Response{Code: 100, Body: ""},
			expectedErr: ValidationErrors{ValidationError{
				Field: "Code",
				Err:   errors.New("should be in 200,404,500"),
			}},
		},
		{
			in:          Response{Code: 404, Body: ""},
			expectedErr: ValidationErrors{},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			var validationErrors ValidationErrors

			ok := errors.As(Validate(tt.in), &validationErrors)
			require.True(t, ok)
			require.Equal(t, tt.expectedErr, validationErrors)
			_ = tt
		})
	}
}
