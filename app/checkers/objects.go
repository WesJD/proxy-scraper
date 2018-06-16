package checkers

import "time"

type Checker interface {
	Check(url string, trueResponse string) (*CheckResult, error)
	WaitTime() time.Duration
}

type CheckResult struct {
	Passing int
	Failing int
	WorkingProxies []string
}
