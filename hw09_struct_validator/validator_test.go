package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	SubjectScores struct {
		Name   string
		Scores []int `validate:"in:1,2,3,4,5"`
	}

	User struct {
		ID     string `json:"id" validate:"length:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regExp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"length:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"length:5"`
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

	WrongRegExp struct {
		AnyField string `validate:"regExp:\\w+@\\w+\\.("`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		testName        string
		in              interface{}
		expectedErr     error
		isValidationErr bool
	}{
		{
			"SubjectScores struct, no error",
			SubjectScores{
				Name:   "Fedor Ivanov",
				Scores: []int{3, 5, 4, 4, 5},
			},
			nil, false,
		},
		{
			"User struct, no error",
			User{
				ID:     uuid.New().String(),
				Name:   "Petr",
				Age:    18,
				Email:  "petrpervyi@mail.ru",
				Role:   "admin",
				Phones: getRandomPhones(),
				meta:   nil,
			},
			nil, false,
		},
		{
			"Token struct, no erorr",
			Token{
				Header:    []byte("abc"),
				Payload:   []byte("def"),
				Signature: []byte("xyz"),
			},
			nil, false,
		},
		{
			"App struct, validation (string len) error",
			App{
				Version: "0.0.15",
			},
			fmt.Errorf("validation failed: field 'Version': length is not equal to 5. "),
			true,
		},
		{
			"Response struct, validation (intSet) error",
			Response{
				Code: 201,
				Body: "{name:\"Igor\"}",
			},
			fmt.Errorf("validation failed: field 'Code': value (201) is not in set ([200 404 500]). "),
			true,
		},
		{
			"User struct, multiple validation errors",
			User{
				ID:     uuid.New().String(),
				Name:   "Anna",
				Age:    25,
				Email:  "anya-thebest@google.com",
				Role:   "staff",
				Phones: getRandomPhones(),
				meta:   nil,
			},
			fmt.Errorf("validation failed: " +
				"field 'Email': value (anya-thebest@google.com) is not matched with regexp ^\\w+@\\w+\\.\\w+$. " +
				"field 'Role': value (staff) is not in set ([admin stuff]). "),
			true,
		},
		{
			"WrongRegExp struct, non-validation error",
			WrongRegExp{
				AnyField: "any string",
			},
			fmt.Errorf("invalid regexp provided: \\w+@\\w+\\.("),
			false,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			result := Validate(tt.in)
			switch {
			case tt.expectedErr == nil:
				require.NoError(t, result)
			case tt.isValidationErr:
				require.Error(t, result)
				var vErrors ValidationErrors
				if errors.As(result, &vErrors) {
					require.EqualError(t, result, tt.expectedErr.Error())
				} else {
					t.Errorf("error is not of type ValidationErrors")
				}
			default:
				require.Error(t, result)
				require.EqualError(t, result, tt.expectedErr.Error())
			}
			_ = tt
		})
	}
}

func getRandomPhones() []string {
	numsNumber := rand.Intn(5) + 1
	phones := make([]string, 0)

	for i := 0; i < numsNumber; i++ {
		phone := "7"
		for i := 0; i < 10; i++ {
			phone += fmt.Sprintf("%d", rand.Intn(10))
		}
		phones = append(phones, phone)
	}

	return phones
}
