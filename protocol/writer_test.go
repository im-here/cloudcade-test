package protocol

import (
	"bytes"
	"testing"
)

func TestWriteCommand(t *testing.T) {
	tests := []struct {
		commands []interface{}
		result   string
	}{
		{
			commands: []interface{}{
				SendCommand{"Hello", ""},
			},
			result: "SEND Hello\n",
		},
	}

	buf := new(bytes.Buffer)
	for _, test := range tests {
		buf.Reset()
		cmdWriter := NewCommandWriter(buf)

		for _, cmd := range test.commands {
			if cmdWriter.Write(cmd) != nil {
				t.Errorf("Unable to write command %v", cmd)
			}
		}

		if buf.String() != test.result {
			t.Errorf("Command output is not the same: %v %v", buf.String(), test.result)
		}
	}
}
