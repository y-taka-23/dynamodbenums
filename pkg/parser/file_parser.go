package parser

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"

	"golang.org/x/xerrors"
)

type FileParser struct {
	path string
}

func NewFileParser(path string) FileParser {
	return FileParser{path: path}
}

func (p FileParser) Parse() (map[string]*types.Package, error) {

	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, p.path, nil, 0)
	if err != nil {
		return map[string]*types.Package{},
			xerrors.Errorf("failed to parse the file: %w", err)
	}

	pkgName := f.Name.Name

	conf := types.Config{Importer: importer.Default()}
	checked, err := conf.Check(pkgName, fset, []*ast.File{f}, nil)
	if err != nil {
		return map[string]*types.Package{},
			xerrors.Errorf("failed to check types: %w", err)
	}

	return map[string]*types.Package{pkgName: checked}, nil
}
