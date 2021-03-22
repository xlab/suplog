package bugsnag

import "fmt"

type reportPublisher interface {
	publishReport(*payload) error
}

type defaultReportPublisher struct{}

func (*defaultReportPublisher) publishReport(p *payload) error {
	if len(p.Stacktrace) > 0 {
		p.logf("notifying bugsnag %s: %s (src: %s:%d)", p.Severity.String, p.Message, p.Stacktrace[0].File, p.Stacktrace[0].LineNumber)
	} else {
		p.logf("notifying bugsnag %s: %s", p.Severity.String, p.Message)
	}

	if !p.notifyInReleaseStage() {
		return fmt.Errorf("not notifying in %s", p.ReleaseStage)
	}
	if p.Synchronous {
		return p.deliver()
	}

	go func(p *payload) {
		if err := p.deliver(); err != nil {
			// Ensure that any errors are logged if they occur in a goroutine.
			p.logf("bugsnag/defaultReportPublisher.publishReport: %v", err)
		}
	}(p)
	return nil
}
