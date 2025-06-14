package graceful

import (
	"context"
	"errors"
	"testing"
)

// dummyStarter реализация starter, возвращающая заранее заданную ошибку.
type dummyStarter struct {
	err error
}

func (d dummyStarter) Start(_ context.Context) error {
	return d.err
}

func TestProcessDisableEnable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		enable           bool
		wantAfterEnable  bool
		disable          bool
		wantAfterDisable bool
	}{
		{
			name:             "EnableTrue",
			enable:           true,
			wantAfterEnable:  false,
			disable:          false,
			wantAfterDisable: false,
		},
		{
			name:             "EnableFalse",
			enable:           false,
			wantAfterEnable:  true,
			disable:          false,
			wantAfterDisable: false,
		},
		{
			name:             "DisableTrue",
			enable:           false,
			wantAfterEnable:  true,
			disable:          true,
			wantAfterDisable: true,
		},
		{
			name:             "DisableFalse",
			enable:           false,
			wantAfterEnable:  true,
			disable:          false,
			wantAfterDisable: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			p := NewProcess(dummyStarter{err: nil})

			pE := p.Enable(tc.enable)
			if pE.disabled != tc.wantAfterEnable {
				t.Errorf("Enable(%v): disabled = %v; want %v", tc.enable, pE.disabled, tc.wantAfterEnable)
			}

			pD := p.Disable(tc.disable)
			if pD.disabled != tc.wantAfterDisable {
				t.Errorf("Disable(%v): disabled = %v; want %v", tc.disable, pD.disabled, tc.wantAfterDisable)
			}
		})
	}
}

func TestStart(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	tests := []struct {
		name    string
		procs   []process
		wantErr bool
	}{
		{
			name:    "NoProcess",
			procs:   nil,
			wantErr: false,
		},
		{
			name: "SingleSuccess",
			procs: []process{
				NewProcess(dummyStarter{err: nil}),
			},
			wantErr: false,
		},
		{
			name: "SingleError",
			procs: []process{
				NewProcess(dummyStarter{err: errFoo}),
			},
			wantErr: true,
		},
		{
			name: "DisabledError",
			procs: []process{
				NewProcess(dummyStarter{err: errFoo}).Disable(true),
			},
			wantErr: false,
		},
		{
			name: "MultiMixed",
			procs: []process{
				NewProcess(dummyStarter{err: nil}),
				NewProcess(dummyStarter{err: errFoo}),
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gr := New(tc.procs...)

			err := gr.Start(context.Background())
			if (err != nil) != tc.wantErr {
				t.Errorf("Start() error = %v; wantErr %v", err, tc.wantErr)
			}
		})
	}
}
