package dag

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func toMessages(args []any) string {
	if len(args) == 0 {
		return ""
	}
	var msgs = make([]string, len(args))
	for i, arg := range args {
		msgs[i] = fmt.Sprintf("%#v", arg)
	}
	return strings.Join(msgs, ", ")
}

func assertEqual(t *testing.T, a, b interface{}, args ...any) {
	t.Helper()
	if !reflect.DeepEqual(a, b) {
		t.Errorf("assertEqual faild: %#v != %#v with messages %s",
			a, b, toMessages(args))
	}
}

func assertNotEqual(t *testing.T, a, b interface{}, args ...any) {
	t.Helper()
	if reflect.DeepEqual(a, b) {
		t.Errorf("assertNotEqual faild: %#v equal to %#v with messages %s", a, b,
			toMessages(args))
	}
}

func assertNil(t *testing.T, a interface{}, args ...any) {
	t.Helper()
	assertEqual(t, a, nil, args...)
}

func assertNotNil(t *testing.T, a interface{}, args ...any) {
	t.Helper()
	assertNotEqual(t, a, nil, args...)
}
