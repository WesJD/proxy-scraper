package chrome

import (
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
	"context"
			)

var (
	instances = make(map[string]*Instance)
	options = chromedp.WithRunnerOptions(
		runner.Flag("headless", true),
		runner.Flag("no-default-browser-check", true),
		runner.Flag("disable-gpu", true),
		runner.Flag("no-sandbox", true))
)

type Instance struct {
	Chrome *chromedp.CDP
	Context context.Context
	ContextCancel context.CancelFunc
}

func DpInstance(holder string) (data *Instance, err error) {
	data = instances[holder]
	if data != nil {
		return
	}

	ctxt, cancel := context.WithCancel(context.Background())

	chrome, err := chromedp.New(ctxt, options)
	if err != nil {
		return
	}

	data = &Instance{
		Chrome: chrome,
		Context: ctxt,
		ContextCancel: cancel,
	}
	instances[holder] = data
	return
}

func CloseInstances() {
	for _, instance := range instances {
		instance.Chrome.Shutdown(instance.Context)
		instance.Chrome.Wait()
		instance.ContextCancel()
	}
}