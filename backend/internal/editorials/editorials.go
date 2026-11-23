// Package editorials has the Editorial of each learning set task and each bank problem, by problem slug.
package editorials

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"strings"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// Editorial is the solutions of one problem: HTML in the lesson format and its figures in the order of the placeholders.
type Editorial struct {
	HTML    string        `json:"html"`
	Figures []sets.Figure `json:"figures"`
}

//go:embed *.html *.json
var files embed.FS

// All is the Editorials by problem slug. A file without its pair or a bad figure file stops the start. TestEditorialContent reads the same map.
var All = func() map[string]Editorial {
	all, err := load(files)
	if err != nil {
		panic(err)
	}
	return all
}()

// load reads each <slug>.html with its figures in <slug>.json.
func load(fsys fs.FS) (map[string]Editorial, error) {
	jsons, err := fs.Glob(fsys, "*.json")
	if err != nil {
		return nil, err
	}
	for _, name := range jsons {
		if _, err := fs.Stat(fsys, strings.TrimSuffix(name, ".json")+".html"); err != nil {
			return nil, fmt.Errorf("editorial %s: %w", name, err)
		}
	}
	htmls, err := fs.Glob(fsys, "*.html")
	if err != nil {
		return nil, err
	}
	all := map[string]Editorial{}
	for _, name := range htmls {
		slug := strings.TrimSuffix(name, ".html")
		html, err := fs.ReadFile(fsys, name)
		if err != nil {
			return nil, fmt.Errorf("editorial %s: %w", name, err)
		}
		data, err := fs.ReadFile(fsys, slug+".json")
		if err != nil {
			return nil, fmt.Errorf("editorial %s: %w", name, err)
		}
		// A typo in a figure field, such as mark for marked, fails the start.
		decoder := json.NewDecoder(bytes.NewReader(data))
		decoder.DisallowUnknownFields()
		var figures []sets.Figure
		if err := decoder.Decode(&figures); err != nil {
			return nil, fmt.Errorf("editorial %s.json: %w", slug, err)
		}
		all[slug] = Editorial{HTML: string(html), Figures: figures}
	}
	return all, nil
}
