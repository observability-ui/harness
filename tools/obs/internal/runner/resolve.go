package runner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"obs/internal/component"
)

func ResolveSpec(spec component.ProcessSpec) (component.ProcessSpec, func(), error) {
	resolved := spec

	args, tempDir, err := resolveFiles(spec.Args, spec.Files)
	if err != nil {
		return spec, nil, err
	}
	resolved.Args = args

	if spec.StdinFile != "" && len(spec.Files) > 0 {
		_, content, err := readFileRef(spec.StdinFile, spec.Files)
		if err != nil {
			if tempDir != "" {
				os.RemoveAll(tempDir)
			}
			return spec, nil, err
		}
		resolved.Stdin = string(content)
		resolved.StdinFile = ""
	}

	resolved.Files = nil

	cleanup := func() {
		if tempDir != "" {
			os.RemoveAll(tempDir)
		}
	}
	return resolved, cleanup, nil
}

func readFileRef(name string, files map[string]component.FileRef) (string, []byte, error) {
	ref, ok := files[name]
	if !ok {
		return "", nil, fmt.Errorf("file %q not found in Files map", name)
	}
	content, err := fs.ReadFile(ref.FS, ref.Path)
	if err != nil {
		return "", nil, fmt.Errorf("read file %q: %w", ref.Path, err)
	}
	return ref.Path, content, nil
}

func resolveFiles(args []string, files map[string]component.FileRef) ([]string, string, error) {
	if len(files) == 0 {
		return args, "", nil
	}

	resolved := make([]string, len(args))
	var tempDir string

	for i, arg := range args {
		if strings.HasPrefix(arg, "{{content:") && strings.HasSuffix(arg, "}}") {
			name := arg[len("{{content:") : len(arg)-2]
			_, content, err := readFileRef(name, files)
			if err != nil {
				return nil, "", err
			}
			resolved[i] = string(content)
		} else if strings.HasPrefix(arg, "{{path:") && strings.HasSuffix(arg, "}}") {
			name := arg[len("{{path:") : len(arg)-2]
			path, content, err := readFileRef(name, files)
			if err != nil {
				return nil, "", err
			}
			if tempDir == "" {
				tempDir, err = os.MkdirTemp("", "obs-files-*")
				if err != nil {
					return nil, "", fmt.Errorf("create temp dir: %w", err)
				}
			}
			tmpPath := filepath.Join(tempDir, filepath.Base(path))
			if err := os.WriteFile(tmpPath, content, 0o644); err != nil {
				return nil, tempDir, fmt.Errorf("write temp file: %w", err)
			}
			resolved[i] = tmpPath
		} else {
			resolved[i] = arg
		}
	}
	return resolved, tempDir, nil
}
