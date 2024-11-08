package errs_test

import (
	"strings"
	"testing"

	"github.com/sulton0011/errs"
)

func TestJoinMsg_VeryLongInput(t *testing.T) {
	const sep = ","
	const numMessages = 100000
	veryLongString := strings.Repeat("very_long_string", 1000)

	messages := make([]any, numMessages)
	for i := 0; i < numMessages; i++ {
		messages[i] = veryLongString
	}

	result := errs.JoinMsg(sep, messages...)
	expectedLength := (numMessages-1)*len(sep) + numMessages*len(veryLongString)

	if len(result) != expectedLength {
		t.Errorf("Expected length %d, but got %d", expectedLength, len(result))
	}
}
func TestJoinMsg_ErrorsWithNumbersAndLetters(t *testing.T) {
	const sep = ","
	errors := []error{
		errs.New("error1"),
		errs.New("2error"),
		errs.New("error34"),
		errs.New("error_5"),
	}

	result := errs.Join(sep, errors...)
	expectedOrigErr := "error1,2error,error34,error_5"
	expectedMessage := "error1,2error,error34,error_5"

	if result.Error() != expectedOrigErr {
		t.Errorf("Expected original error '%s', but got '%s'", expectedOrigErr, result.Error())
	}

	if errs.Unwrap(result) != expectedMessage {
		t.Errorf("Expected message '%s', but got '%s'", expectedMessage, errs.Unwrap(result))
	}
}
