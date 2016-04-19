package log_test

import (
	"testing"
	"util/log"
)

func TestLog(t *testing.T) {
	log.Init()

	log.AddLog("this is log")
	log.AddAssert("this is assert")
	log.AddError("this is error")
	log.AddException("this is exception")
	log.AddWarning("this is warning")
}