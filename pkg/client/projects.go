package client

import (
	"context"
	_ "embed"

	"github.com/go-git/go-billy/v5/osfs"

	"github.com/jh125486/CSCE4350_gradebot/pkg/rubrics"
	"github.com/jh125486/gradebot/pkg/client"
	baserubrics "github.com/jh125486/gradebot/pkg/rubrics"
)

var (
	//go:embed instructions/project1.txt
	project1Instructions string
)

// ExecuteProject1 executes the project1 grading flow using a runtime config.
func ExecuteProject1(ctx context.Context, cfg *client.Config) error {
	return client.ExecuteProject(ctx, cfg, "CSCE4350:Project1", project1Instructions,
		baserubrics.RunBag{},
		baserubrics.EvaluateGit(osfs.New(cfg.WorkDir.String())),
		rubrics.EvaluateDataFileCreated,
		rubrics.EvaluateSetGet,
		rubrics.EvaluateOverwriteKey,
		rubrics.EvaluateNonexistentGet,
		rubrics.EvaluatePersistenceAfterRestart,
	)
}
