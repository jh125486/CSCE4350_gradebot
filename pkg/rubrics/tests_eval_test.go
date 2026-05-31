package rubrics_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jh125486/CSCE4350_gradebot/pkg/rubrics"
	baserubrics "github.com/jh125486/gradebot/pkg/rubrics"
)

const testFilesPresentMsg = "Project contains test files/directories"

func TestEvaluateTestsPresent(t *testing.T) {
	tests := []struct {
		name       string
		setupMock  func(t *testing.T, dir string)
		wantPoints float64
		wantNote   string
	}{
		{
			name: "no tests",
			setupMock: func(t *testing.T, dir string) {
				os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0o644)
			},
			wantPoints: 0,
			wantNote:   "No test files detected",
		},
		{
			name: "go tests",
			setupMock: func(t *testing.T, dir string) {
				os.WriteFile(filepath.Join(dir, "main_test.go"), []byte("package main"), 0o644)
			},
			wantPoints: 5,
			wantNote:   testFilesPresentMsg,
		},
		{
			name: "python tests",
			setupMock: func(t *testing.T, dir string) {
				os.WriteFile(filepath.Join(dir, "test_main.py"), []byte(""), 0o644)
			},
			wantPoints: 5,
			wantNote:   testFilesPresentMsg,
		},
		{
			name: "tests directory",
			setupMock: func(t *testing.T, dir string) {
				os.MkdirAll(filepath.Join(dir, "tests"), 0o755)
				os.WriteFile(filepath.Join(dir, "tests", "dummy.txt"), []byte(""), 0o644)
			},
			wantPoints: 5,
			wantNote:   testFilesPresentMsg,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newKVStoreMock(t)
			tt.setupMock(t, mock.tempDir)

			result := rubrics.EvaluateTestsPresent(context.Background(), mock, make(baserubrics.RunBag))

			assert.Equal(t, tt.wantPoints, result.Awarded)
			assert.Contains(t, result.Note, tt.wantNote)
		})
	}
}
