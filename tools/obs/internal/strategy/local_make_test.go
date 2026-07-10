package strategy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveTarget_SingleCandidate(t *testing.T) {
	target, err := resolveTarget("", "build")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if target != "build" {
		t.Fatalf("expected 'build', got %q", target)
	}
}

func TestResolveTarget_SingleCandidateTrimmed(t *testing.T) {
	target, err := resolveTarget("", "  build  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if target != "build" {
		t.Fatalf("expected 'build', got %q", target)
	}
}

func TestResolveTarget_MultipleCandidatesPicksFirst(t *testing.T) {
	dir := t.TempDir()
	writeMakefile(t, dir, "start-frontend")

	target, err := resolveTarget(dir, "start-frontend,start-console")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if target != "start-frontend" {
		t.Fatalf("expected 'start-frontend', got %q", target)
	}
}

func TestResolveTarget_MultipleCandidatesFallsBack(t *testing.T) {
	dir := t.TempDir()
	writeMakefile(t, dir, "start-console")

	target, err := resolveTarget(dir, "start-frontend,start-console")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if target != "start-console" {
		t.Fatalf("expected 'start-console', got %q", target)
	}
}

func TestResolveTarget_MultipleCandidatesNoneFound(t *testing.T) {
	dir := t.TempDir()
	writeMakefile(t, dir, "other-target")

	_, err := resolveTarget(dir, "start-frontend,start-console")
	if err == nil {
		t.Fatal("expected error when no candidates match, got nil")
	}
}

func TestMakeTargetExists(t *testing.T) {
	dir := t.TempDir()
	writeMakefile(t, dir, "build")

	if !makeTargetExists(dir, "build") {
		t.Fatal("expected 'build' target to exist")
	}
	if makeTargetExists(dir, "nonexistent") {
		t.Fatal("expected 'nonexistent' target to not exist")
	}
}

func writeMakefile(t *testing.T, dir, target string) {
	t.Helper()
	content := target + ":\n\t@echo ok\n"
	if err := os.WriteFile(filepath.Join(dir, "Makefile"), []byte(content), 0644); err != nil {
		t.Fatalf("failed to write Makefile: %v", err)
	}
}
