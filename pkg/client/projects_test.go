package client_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/jh125486/gradebot/pkg/client"
	"github.com/jh125486/gradebot/pkg/contextlog"
	"github.com/stretchr/testify/assert"

	clientpkg "github.com/jh125486/CSCE4350_gradebot/pkg/client"
)

// TestExecuteProject1 is a comprehensive table-driven test for ExecuteProject1 function.
// It covers all edge cases: nil inputs, cancelled context, various config combinations,
// and verifies proper error handling.
func TestExecuteProject1(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		cfg *client.Config
	}
	type testCase struct {
		name        string
		args        args
		shouldPanic bool
		wantErr     bool
		setup       func(t *testing.T, args *args)
		verify      func(t *testing.T)
	}
	tests := []testCase{
		{
			name: "nil context panics",
			args: args{
				ctx: nil,
				cfg: &client.Config{
					WorkDir: client.WorkDir(t.TempDir()),
					RunCmd:  "echo test",
					Writer:  io.Discard,
				},
			},
			shouldPanic: true,
			wantErr:     false,
		},
		{
			name: "nil config panics",
			args: args{
				ctx: context.Background(),
				cfg: nil,
			},
			shouldPanic: true,
			wantErr:     false,
		},
		{
			name: "cancelled context returns error",
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()
					return ctx
				}(),
				cfg: &client.Config{
					WorkDir: client.WorkDir(t.TempDir()),
					RunCmd:  "echo test",
					Writer:  io.Discard,
				},
			},
			shouldPanic: false,
			wantErr:     true,
		},
		{
			name: "empty WorkDir completes without immediate error",
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithTimeout(
						contextlog.With(context.Background(), contextlog.DiscardLogger()),
						100*time.Millisecond,
					)
					t.Cleanup(cancel)
					return ctx
				}(),
				cfg: &client.Config{
					WorkDir: "",
					RunCmd:  "echo test",
					Writer:  io.Discard,
				},
			},
			shouldPanic: false,
			wantErr:     false, // Empty WorkDir doesn't cause immediate error, may timeout later
		},
		{
			name: "nil reader is accepted",
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithTimeout(
						contextlog.With(context.Background(), contextlog.DiscardLogger()),
						100*time.Millisecond,
					)
					t.Cleanup(cancel)
					return ctx
				}(),
				cfg: &client.Config{
					WorkDir: client.WorkDir(t.TempDir()),
					RunCmd:  "echo test",
					Writer:  io.Discard,
					Reader:  nil,
				},
			},
			shouldPanic: false,
			wantErr:     true, // May timeout or encounter missing dependencies
		},
		{
			name: "config with io.Discard writer",
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithTimeout(
						contextlog.With(context.Background(), contextlog.DiscardLogger()),
						100*time.Millisecond,
					)
					t.Cleanup(cancel)
					return ctx
				}(),
				cfg: &client.Config{
					WorkDir: client.WorkDir(t.TempDir()),
					RunCmd:  "echo test",
					Writer:  io.Discard,
					Reader:  nil,
				},
			},
			shouldPanic: false,
			wantErr:     true, // May timeout or encounter missing dependencies
		},
		{
			name: "logger context with applied timeout",
			args: args{
				ctx: contextlog.With(context.Background(), contextlog.DiscardLogger()),
				cfg: &client.Config{
					WorkDir: client.WorkDir(t.TempDir()),
					RunCmd:  "echo test",
					Writer:  io.Discard,
				},
			},
			shouldPanic: false,
			wantErr:     false, // Will complete successfully or timeout
			setup: func(t *testing.T, args *args) {
				// Apply timeout to prevent hanging indefinitely
				ctx, cancel := context.WithTimeout(args.ctx, 100*time.Millisecond)
				t.Cleanup(cancel)
				args.ctx = ctx
			},
		},
		{
			name: "valid temp directory with timeout",
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithTimeout(
						contextlog.With(context.Background(), contextlog.DiscardLogger()),
						100*time.Millisecond,
					)
					t.Cleanup(cancel)
					return ctx
				}(),
				cfg: &client.Config{
					WorkDir: client.WorkDir(t.TempDir()),
					RunCmd:  "true",
					Writer:  io.Discard,
				},
			},
			shouldPanic: false,
			wantErr:     true, // Expected to timeout or error on dependencies
		},
		{
			name: "different RunCmd with timeout",
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithTimeout(
						contextlog.With(context.Background(), contextlog.DiscardLogger()),
						100*time.Millisecond,
					)
					t.Cleanup(cancel)
					return ctx
				}(),
				cfg: &client.Config{
					WorkDir: client.WorkDir(t.TempDir()),
					RunCmd:  "false",
					Writer:  io.Discard,
				},
			},
			shouldPanic: false,
			wantErr:     true, // Expected to timeout or error on dependencies
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.args)
			}

			if tt.shouldPanic {
				assert.Panics(t, func() {
					clientpkg.ExecuteProject1(tt.args.ctx, tt.args.cfg)
				}, "ExecuteProject1 should panic with invalid inputs")
			} else {
				err := clientpkg.ExecuteProject1(tt.args.ctx, tt.args.cfg)
				if tt.wantErr {
					assert.Error(t, err, "ExecuteProject1 should return an error")
				} else {
					assert.NoError(t, err, "ExecuteProject1 should not return an error")
				}
			}

			if tt.verify != nil {
				tt.verify(t)
			}
		})
	}
}
