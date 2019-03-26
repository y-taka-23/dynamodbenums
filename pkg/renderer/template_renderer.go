package renderer

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/xerrors"

	"github.com/y-taka-23/dynamodbenums/pkg/enum"
)

type TemplateRenderer struct {
	template template.Template
	dir      string
	prefix   string
	suffix   string
}

func NewTemplateRenderer(tpl string, dir, prefix, suffix string) (TemplateRenderer, error) {

	t, err := template.New("").Parse(tpl)
	if err != nil {
		return TemplateRenderer{},
			xerrors.Errorf("failed to parse the template: %w", err)
	}

	return TemplateRenderer{
		template: *t,
		dir:      dir,
		prefix:   prefix,
		suffix:   suffix,
	}, nil
}

func (r TemplateRenderer) Render(enum enum.Enum) error {

	if len(enum.Values) == 0 {
		return nil
	}

	output := outputFile(enum.TypeName, r.dir, r.prefix, r.suffix)

	f, err := os.Create(output)
	if err != nil {
		return xerrors.Errorf("failed to create a file: %w", err)
	}
	defer f.Close()

	return r.template.Execute(f, enum)
}

func outputFile(typeName, dir, prefix, suffix string) string {
	base := prefix + strings.ToLower(typeName) + suffix + ".go"
	return filepath.Join(dir, base)
}
