package strategy

import (
	"testing"

	"obs/internal/task"
)

func TestMakeTarget_ReturnsStrategy(t *testing.T) {
	s := MakeTarget("build")
	mr, ok := s.(*makeRun)
	if !ok {
		t.Fatalf("expected *makeRun, got %T", s)
	}
	if mr.Target != "build" {
		t.Fatalf("expected target 'build', got %q", mr.Target)
	}
}

func TestNPMRun_ReturnsStrategy(t *testing.T) {
	s := NPMRun("install", "--no-save")
	n, ok := s.(*npm)
	if !ok {
		t.Fatalf("expected *npm, got %T", s)
	}
	if n.Cmd != "install" {
		t.Fatalf("expected cmd 'install', got %q", n.Cmd)
	}
}

func TestCompose_ReturnsStrategy(t *testing.T) {
	s := Compose("docker-compose.yml")
	pc, ok := s.(*podmanCompose)
	if !ok {
		t.Fatalf("expected *podmanCompose, got %T", s)
	}
	if pc.File != "docker-compose.yml" {
		t.Fatalf("expected file 'docker-compose.yml', got %q", pc.File)
	}
}

func TestDockerBuild_ReturnsStrategy(t *testing.T) {
	s := DockerBuild("Dockerfile.dev")
	cb, ok := s.(*containerBuild)
	if !ok {
		t.Fatalf("expected *containerBuild, got %T", s)
	}
	if cb.Dockerfile != "Dockerfile.dev" {
		t.Fatalf("expected dockerfile 'Dockerfile.dev', got %q", cb.Dockerfile)
	}
}

func TestMakeRun_ImplementsStrategy(t *testing.T) {
	var _ task.Strategy = &makeRun{}
}

func TestNPM_ImplementsStrategy(t *testing.T) {
	var _ task.Strategy = &npm{}
}

func TestPodmanCompose_ImplementsStrategy(t *testing.T) {
	var _ task.Strategy = &podmanCompose{}
}

func TestContainerBuild_ImplementsStrategy(t *testing.T) {
	var _ task.Strategy = &containerBuild{}
}
