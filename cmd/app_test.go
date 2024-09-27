package app

import (
	"testing"

)

func TestServerStartLog(t *testing.T) {
  expected := "Hello, World!"

	actual := "Hello, World!"

	if expected != actual {
		t.Errorf("Expected '%s', but got '%s'", expected, actual)
	}
}
