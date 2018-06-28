package sites

import (
	"time"
)

type Site interface {
	Check(url string) (map[string]bool, error)
	WaitTime() time.Duration
}
