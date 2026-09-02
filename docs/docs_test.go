package docs_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/suzuki-shunsuke/cobra-util/docs"
)

func TestParse(t *testing.T) {
	t.Parallel()
	data := []struct {
		name    string
		content string
		want    *docs.Result
		wantErr bool
	}{
		{
			name:    "the description of the document",
			content: "---\ndescription: How to install.\n---\n\n# Install\n",
			want:    &docs.Result{Name: install, Description: "How to install."},
		},
		{
			name:    "the name in the frontmatter doesn't overwrite the file name",
			content: "---\nname: other\ndescription: How to install.\n---\n",
			want:    &docs.Result{Name: install, Description: "How to install."},
		},
		{
			name:    "no frontmatter",
			content: "# Install\n",
			wantErr: true,
		},
		{
			name:    "an unclosed frontmatter",
			content: "---\ndescription: How to install.\n",
			wantErr: true,
		},
		{
			name:    "a frontmatter that isn't YAML",
			content: "---\ndescription: [\n---\n",
			wantErr: true,
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			result := &docs.Result{Name: install}
			err := docs.Parse([]byte(d.content), result)
			if d.wantErr {
				if err == nil {
					t.Fatal("an error must be returned")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(d.want, result); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
