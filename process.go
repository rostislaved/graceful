package graceful

import (
	"context"
)

type process struct {
	starter  starter
	disabled bool
}

func NewProcess(starter starter) process {
	return process{
		starter:  starter,
		disabled: false,
	}
}

type starter interface {
	Start(ctx context.Context) error
}

func (p process) Disable(disable bool) process {
	p.disabled = disable

	return p
}

func (p process) Enable(enable bool) process {
	p.disabled = !enable

	return p
}
