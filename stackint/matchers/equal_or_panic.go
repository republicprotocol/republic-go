package matchers

import (
	"fmt"
	"reflect"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/matchers"
	"github.com/onsi/gomega/types"
)

func EqualOrPanic(expected interface{}) types.GomegaMatcher {
	return &EqualOrPanicMatcher{
		Expected: expected,
	}
}

type EqualOrPanicMatcher struct {
	Expected interface{}
}

func (matcher *EqualOrPanicMatcher) Match(actual interface{}) (success bool, err error) {
	switch matcher.Expected.(type) {
	case *matchers.PanicMatcher:
		panic := PanicMatcher{}
		return panic.Match(actual)
	default:
		x := make([]reflect.Value, 0)
		actualValue := reflect.ValueOf(actual).Call(x)[0].Interface()
		return gomega.Equal(matcher.Expected).Match(actualValue)
	}
}

func (matcher *EqualOrPanicMatcher) FailureMessage(actual interface{}) (message string) {
	actualString, actualOK := actual.(string)
	expectedString, expectedOK := matcher.Expected.(string)
	if actualOK && expectedOK {
		return format.MessageWithDiff(actualString, "to equal", expectedString)
	}

	return format.Message(actual, "to equal", matcher.Expected)
}

func (matcher *EqualOrPanicMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to equal", matcher.Expected)
}

type PanicMatcher struct {
	object interface{}
}

func (matcher *PanicMatcher) Match(actual interface{}) (success bool, err error) {
	if actual == nil {
		return false, fmt.Errorf("PanicMatcher expects a non-nil actual.")
	}

	actualType := reflect.TypeOf(actual)
	if actualType.Kind() != reflect.Func {
		return false, fmt.Errorf("PanicMatcher expects a function.  Got:\n%s", format.Object(actual, 1))
	}

	success = false
	defer func() {
		if e := recover(); e != nil {
			matcher.object = e
			success = true
		}
	}()

	reflect.ValueOf(actual).Call([]reflect.Value{})

	return
}

func (matcher *PanicMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to panic")
}

func (matcher *PanicMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, fmt.Sprintf("not to panic, but panicked with\n%s", format.Object(matcher.object, 1)))
}
