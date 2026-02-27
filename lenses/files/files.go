package files

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/stikypiston/spyglass/lens"
)

type filesLens struct {
	files []string
	home  string
}

func New() lens.Lens {
	home, _ := os.UserHomeDir()
	l := &filesLens{
		home: home,
	}
	l.index()
	return l
}

func (f *filesLens) Name() string {
	return "Files"
}

func (f *filesLens) index() {
	filepath.WalkDir(f.home, func(path string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			f.files = append(f.files, path)
		}
		return nil
	})
}

func (f *filesLens) Search(query string) ([]lens.Entry, error) {
	var results []lens.Entry
	query = strings.ToLower(query)

	for _, file := range f.files {
		if strings.Contains(strings.ToLower(file), query) {
			results = append(results, lens.Entry{
				ID:          file,
				Title:       shortenPath(file, f.home),
				Icon:        "󰈔",
				Description: file,
			})
		}
	}
	return results, nil
}

func (f *filesLens) Enter(entry lens.Entry) error {
	return exec.Command("xdg-open", entry.ID).Start()
}

func (f *filesLens) ContextActions(entry lens.Entry) []lens.Action {
	return nil
}

// Shorten filepath (inspired by Fish shell prompt)
func shortenPath(fullPath, home string) string {
	// Replace home with ~
	if strings.HasPrefix(fullPath, home) {
		fullPath = "~" + strings.TrimPrefix(fullPath, home)
	}

	parts := strings.Split(fullPath, "/")
	if len(parts) <= 2 {
		return fullPath
	}

	// Keep first element (~ or root)
	result := []string{parts[0]}

	// Compress all middle directories
	for i := 1; i < len(parts)-1; i++ {
		dir := parts[i]
		if dir == "" {
			continue
		}

		if len(dir) <= 2 {
			result = append(result, dir)
		} else {
			result = append(result, dir[:2])
		}
	}

	// Keep full filename
	result = append(result, parts[len(parts)-1])

	return strings.Join(result, "/")
}
