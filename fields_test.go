package suplog

import (
	"testing"
)

func TestFnName(t *testing.T) {
	outputCallerName := FnName()
	if outputCallerName != "TestFnName" {
		t.Fail()
	}
}

func TestCallerName(t *testing.T) {
	instanceCallerName := NewLogger(nil, nil).(LoggerConfigurator).CallerName()
	if instanceCallerName != "TestCallerName" {
		t.Fail()
	}
}
