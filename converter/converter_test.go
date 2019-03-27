package converter

import (
	"bytes"
	"testing"
)

func TestConvert(t *testing.T) {
	dest := bytes.Buffer{}
	c := NewConverter(&dest)

	src := bytes.Buffer{}
	err := c.Convert(&src)
	if err != nil {
		t.Log(err)
	}
}
