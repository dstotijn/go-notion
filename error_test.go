package notion

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestAPIError(t *testing.T) {
	t.Parallel()

	t.Run("error formatting", func(t *testing.T) {
		t.Parallel()

		err := APIError{
			Status:  429,
			Code:    "rate_limited",
			Message: "notion: this request exceeds the number of requests allowed",
			prefix:  "blabla",
		}

		exp := "blabla: notion: this request exceeds the number of requests allowed (code: rate_limited, status: 429)"
		got := err.Error()

		if exp != got {
			t.Fatalf("wrong output for Error() method (expected: %v, got: %v)", exp, got)
		}
	})

	t.Run("error parsing", func(t *testing.T) {
		t.Parallel()

		response := &http.Response{
			StatusCode: http.StatusBadRequest,
			Status:     http.StatusText(http.StatusBadRequest),
			Body: ioutil.NopCloser(strings.NewReader(
				`{
					"object": "error",
					"status": 400,
					"code": "validation_error",
					"message": "notion: request body does not match the schema for the expected parameters"
				}`,
			)),
		}

		exp := APIError{
			Status:  400,
			Code:    "validation_error",
			Message: "notion: request body does not match the schema for the expected parameters",
			prefix:  "blabla",
		}
		got := parseErrorResponse(response, "blabla")

		if _got, ok := got.(*APIError); !ok {
			t.Fatalf("parseErrorResponse must return an APIError error")
		} else if _got.Code != exp.Code {
			t.Fatalf("parseErrorResponse did not parsed code correctly (expected: %v, got: %v)", exp.Code, _got.Code)
		}

		if exp.Error() != got.Error() {
			t.Fatalf("parseErrorResponse did not parse body correctly (expected: %v, got: %v)", exp.Error(), got.Error())
		}
	})
}
