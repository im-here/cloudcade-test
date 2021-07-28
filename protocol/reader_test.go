package protocol

import (
	"reflect"
	"strings"
	"testing"
)

func TestCommandReader(t *testing.T) {
	tests := []struct {
		input   string
		results []interface{}
	}{
		{
			input: "SEND test\n",
			results: []interface{}{
				SendCommand{"test", ""},
			},
		},
		{
			input: "MESSAGE user1 hello\nMESSAGE user2 world\n",
			results: []interface{}{
				MessageCommand{"user1", "hello"},
				MessageCommand{"user2", "world"},
			},
		},
	}

	for _, test := range tests {
		reader := NewCommandReader(strings.NewReader(test.input))
		results, err := reader.ReadAll()

		t.Log(results)

		if err != nil {
			t.Errorf("Unable to read command, error %v", err)
		} else if !reflect.DeepEqual(results, test.results) {
			t.Errorf("Command output is not the same: %v %v", results, test.results)
		}
	}
}
