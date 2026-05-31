package rubrics

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	baserubrics "github.com/jh125486/gradebot/pkg/rubrics"
)

// hasTestFileSuffix checks if a file name indicates it's a test file.
func hasTestFileSuffix(name string) bool {
	suffixes := []string{
		"_test.go", "_test.py", ".test.js", ".spec.js",
		".test.ts", ".spec.ts", "test.java", "tests.java",
		"test.rs", "_test.cpp", "_test.c",
	}
	for _, suffix := range suffixes {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}
	return strings.HasPrefix(name, "test_") && strings.HasSuffix(name, ".py")
}

// isTestDir checks if the directory is a test directory.
func isTestDir(info os.FileInfo) bool {
	name := strings.ToLower(info.Name())
	return name == "tests" || name == "test" || name == "spec" || name == "__tests__"
}

// checkDir contains logic to evaluate directories.
func checkDir(path string, info os.FileInfo, hasTests *bool) error {
	if isTestDir(info) {
		if entries, _ := os.ReadDir(path); len(entries) > 0 {
			*hasTests = true
			return filepath.SkipDir
		}
	}
	if info.Name() == "__pycache__" || info.Name() == ".git" {
		return filepath.SkipDir
	}
	return nil
}

// EvaluateTestsPresent checks if the project has any test files.
func EvaluateTestsPresent(_ context.Context, program baserubrics.ProgramRunner, _ baserubrics.RunBag) baserubrics.RubricItem {
	rubricItem := newRubricItem("TestsPresent", 5)

	hasTests := false
	err := filepath.Walk(program.Path(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return checkDir(path, info, &hasTests)
		}

		if hasTestFileSuffix(strings.ToLower(info.Name())) {
			hasTests = true
			return filepath.SkipDir
		}
		return nil
	})

	if err == nil && hasTests {
		return rubricItem(testFilesPresentMsg, 5)
	}

	return rubricItem("No test files detected in the project", 0)
}

const testFilesPresentMsg = "Project contains test files/directories"
