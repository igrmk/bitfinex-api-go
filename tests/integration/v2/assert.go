package tests

import (
	"runtime/debug"
	"testing"

	"github.com/go-test/deep"
)

func init() {
	deep.CompareUnexportedFields = true
}

func assert(t *testing.T, expected interface{}, actual interface{}) {
	if diff := deep.Equal(expected, actual); diff != nil {
		t.Log(string(debug.Stack()))
		t.Error(diff)
	}
}
