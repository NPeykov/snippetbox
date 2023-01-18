package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
    t.Helper()

    if actual != expected {
        t.Errorf("got: %v but expected: %v", actual, expected)
    }
}

func StringContains(t *testing.T, actual, substring string) {
    t.Helper()

    if !strings.Contains(actual, substring) {
        t.Errorf("got: %q but expected to contain %q", actual, substring)
    }
}

func NilError(t *testing.T, actual error) {
    t.Helper()

    if actual != nil {
        t.Errorf("got: %v but expected: nil", actual)
    }
}

