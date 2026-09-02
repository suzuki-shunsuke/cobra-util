// Package docs implements the docs command of CLIs built with spf13/cobra, which
// lists and outputs the documents embedded in the binary so that a coding agent can
// read them before answering questions about the CLI or troubleshooting its errors.
//
// The documents are Markdown files with a YAML frontmatter holding their
// description. They are passed as an fs.FS, usually one built by go:embed, so that
// they always describe the version that is running.
//
//	//go:embed docs/*.md
//	var docsFS embed.FS
package docs

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"

	"gopkg.in/yaml.v3"
)

// Ext is the file extension of the documents. A document is named by its path in the
// documents' fs.FS without this extension, which is what 'docs show' takes.
const Ext = ".md"

// errNoDocs is returned rather than reporting an empty list of documents, because a
// CLI that embeds no documents shouldn't add the command in the first place, and a
// coding agent reading an empty list would conclude there is nothing to read.
var errNoDocs = errors.New("no documents are embedded in this binary")

// Result is what 'docs list' reports about a document: enough for a coding agent to
// decide whether it is worth reading.
type Result struct {
	// Name is the document name, which comes from the file name rather than from the
	// frontmatter, so that the two can't drift apart.
	Name string `json:"name" yaml:"-"`
	// Description says what the document covers and when to read it. It comes from
	// the frontmatter.
	Description string `json:"description" yaml:"description"`
}

// Names lists the names of the documents in fsys, in lexical order by path.
// Anything that isn't a Markdown file is skipped, so a directory of documents can
// also hold the images they embed.
//
// Subdirectories are walked, and a document in one is named by its path without the
// extension, such as "codes/001", which is the name 'docs show' takes for it. A CLI
// that groups its documents in directories therefore serves them all without having
// to flatten the directory the documentation is written in.
func Names(fsys fs.FS) ([]string, error) {
	names := []string{}
	if err := fs.WalkDir(fsys, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || !strings.HasSuffix(path, Ext) {
			return nil
		}
		names = append(names, strings.TrimSuffix(path, Ext))
		return nil
	}); err != nil {
		return nil, fmt.Errorf("read the docs directory: %w", err)
	}
	return names, nil
}

// Parse reads the YAML frontmatter of a document into result, leaving the fields the
// frontmatter doesn't set as they are, so that the name taken from the file name
// survives.
//
// It is exported so that a test of the CLI can check that every document it embeds
// has a frontmatter, which is otherwise only found out when 'docs list' is run.
func Parse(b []byte, result *Result) error {
	content := string(b)
	const delim = "---\n"
	if !strings.HasPrefix(content, delim) {
		return errors.New("the document has no frontmatter")
	}
	front, _, found := strings.Cut(content[len(delim):], "\n---")
	if !found {
		return errors.New("the document frontmatter is not closed")
	}
	if err := yaml.Unmarshal([]byte(front), result); err != nil {
		return fmt.Errorf("parse the frontmatter as YAML: %w", err)
	}
	return nil
}
