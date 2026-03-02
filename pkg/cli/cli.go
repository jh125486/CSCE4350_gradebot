package cli

import (
	"fmt"

	"github.com/jh125486/CSCE4350_gradebot/pkg/client"
	basecli "github.com/jh125486/gradebot/pkg/cli"
	baseclient "github.com/jh125486/gradebot/pkg/client"
	"github.com/jh125486/gradebot/pkg/proto/protoconnect"
)

type (
	// CLI defines the command-line interface structure for the gradebot application.
	CLI struct {
		basecli.BaseCLI `embed:""`
		Project1        Project1Cmd `cmd:""   help:"Execute project1 grading client"`
	}
	// Project1Cmd defines the command structure for running Project 1 grading.
	Project1Cmd struct {
		basecli.CommonArgs
	}
)

// Run executes the Project 1 grading client.
// The buildID is injected by Kong from the bound value.
func (cmd *Project1Cmd) Run(ctx basecli.Context, svc *basecli.Service) error {
	// Validate required inputs upfront to fail fast
	if cmd.ServerURL == "" {
		return fmt.Errorf("server URL is required")
	}
	if err := cmd.WorkDir.Validate(); err != nil {
		return fmt.Errorf("invalid work directory: %w", err)
	}

	cfg := &baseclient.Config{
		ServerURL:     cmd.ServerURL,
		WorkDir:       cmd.WorkDir,
		RunCmd:        cmd.RunCmd,
		Env:           cmd.Env,
		QualityClient: protoconnect.NewQualityServiceClient(svc.Client, cmd.ServerURL),
		RubricClient:  protoconnect.NewRubricServiceClient(svc.Client, cmd.ServerURL),
		Reader:        svc.Stdin,
		Writer:        svc.Stdout,
	}

	return client.ExecuteProject1(ctx, cfg)
}
