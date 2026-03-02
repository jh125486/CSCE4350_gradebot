package cli_test

import (
	"bytes"
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/jh125486/CSCE4350_gradebot/pkg/cli"
	basecli "github.com/jh125486/gradebot/pkg/cli"
	baseclient "github.com/jh125486/gradebot/pkg/client"
	"github.com/jh125486/gradebot/pkg/contextlog"
)

const (
	testServerURL     = "http://example.invalid"
	testRunCmd        = "echo test"
	testStdinNegative = "n\n"
	testTimeout       = 3 * time.Second
)

func TestWorkDirValidate(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()
	testCases := []struct {
		name    string
		dir     baseclient.WorkDir
		wantErr bool
	}{
		{
			name:    "valid directory",
			dir:     baseclient.WorkDir(tempDir),
			wantErr: false,
		},
		{
			name:    "nonexistent directory",
			dir:     baseclient.WorkDir("./no-such-dir"),
			wantErr: true,
		},
		{
			name:    "empty directory",
			dir:     baseclient.WorkDir(""),
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.dir.Validate()
			if tc.wantErr && err == nil {
				t.Fatalf("expected error for dir %q, got nil", tc.dir)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected no error for dir %q, got %v", tc.dir, err)
			}
		})
	}
}

func TestProject1CmdRun(t *testing.T) {
	t.Parallel()
	type args struct {
		serverURL string
		dir       string
		runCmd    string
		client    *http.Client
		ctx       context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "returns error for cancelled context",
			args: args{
				serverURL: testServerURL,
				dir:       t.TempDir(),
				runCmd:    testRunCmd,
				client: &http.Client{
					Timeout: testTimeout,
				},
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(contextlog.With(context.Background(), contextlog.DiscardLogger()))
					cancel()
					return ctx
				}(),
			},
			wantErr: true,
		},
		{
			name: "returns error for empty server URL",
			args: args{
				serverURL: "",
				dir:       t.TempDir(),
				runCmd:    testRunCmd,
				client: &http.Client{
					Timeout: testTimeout,
				},
				ctx: contextlog.With(context.Background(), contextlog.DiscardLogger()),
			},
			wantErr: true,
		},
		{
			name: "returns error for invalid work directory",
			args: args{
				serverURL: testServerURL,
				dir:       "./nonexistent-path-that-does-not-exist",
				runCmd:    testRunCmd,
				client: &http.Client{
					Timeout: testTimeout,
				},
				ctx: contextlog.With(context.Background(), contextlog.DiscardLogger()),
			},
			wantErr: true,
		},
		{
			name: "returns error for empty work directory",
			args: args{
				serverURL: testServerURL,
				dir:       "",
				runCmd:    testRunCmd,
				client: &http.Client{
					Timeout: testTimeout,
				},
				ctx: contextlog.With(context.Background(), contextlog.DiscardLogger()),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := cli.Project1Cmd{
				CommonArgs: basecli.CommonArgs{
					ServerURL: tt.args.serverURL,
					WorkDir:   baseclient.WorkDir(tt.args.dir),
					RunCmd:    tt.args.runCmd,
				},
			}

			svc := &basecli.Service{
				Client: tt.args.client,
				Stdin:  nil,
				Stdout: new(bytes.Buffer),
			}

			runCtx := tt.args.ctx
			if runCtx == nil {
				runCtx = contextlog.With(t.Context(), contextlog.DiscardLogger())
			}

			err := p.Run(basecli.Context{Context: runCtx}, svc)

			if (err != nil) != tt.wantErr {
				t.Errorf("Project1Cmd.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
