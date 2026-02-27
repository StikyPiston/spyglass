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
}

func New() lens.Lens {
	l := &filesLens{}
	l.index()
	return l
}

func (f *filesLens) Name() string {
	return "Files"
}

func (f *filesLens) index() {
	home, _ := os.UserHomeDir()
	filepath.WalkDir(home, func(path string, d os.DirEntry, err error) error {
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
				Title:       filepath.Base(file),
				Icon:        "󰈔",
				Description: file,
			})
		}
	}
	return results, nil
}

func (f *filesLens) Enter(entry lens.Entry) error {
	return execOpen(entry.ID)
}

func execOpen(path string) error {
	return exec.Command("xdg-open", path).Start()
}

func (f *filesLens) ContextActions(entry lens.Entry) []lens.Action {
	return nil
}
