package command

import (
	"go/types"

	"golang.org/x/xerrors"

	"github.com/y-taka-23/dynamodbenums/pkg/enum"
)

type Parser interface {
	Parse() (map[string]*types.Package, error)
}

type Extractor interface {
	Extract(string, *types.Package) (enum.Enum, error)
}

type Renderer interface {
	Render(enum.Enum) error
}

type Command struct {
	parser    Parser
	extractor Extractor
	renderer  Renderer
}

func New(psr Parser, ext Extractor, rnd Renderer) Command {
	return Command{
		parser:    psr,
		extractor: ext,
		renderer:  rnd,
	}
}

func (c Command) Run(typeNames []string) error {

	pkgs, err := c.parser.Parse()
	if err != nil {
		return xerrors.Errorf("failed to parse the input: %w", err)
	}

	enums := []enum.Enum{}

	for _, typeName := range typeNames {
		for _, pkg := range pkgs {
			enum, err := c.extractor.Extract(typeName, pkg)
			if err != nil {
				return xerrors.Errorf("failed to extract the values: %w", err)
			}
			enums = append(enums, enum)
		}
	}

	for _, enum := range enums {
		if err := c.renderer.Render(enum); err != nil {
			return xerrors.Errorf("failed to render the values: %w", err)
		}
	}

	return nil
}
