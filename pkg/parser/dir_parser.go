package parser

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"

	"golang.org/x/xerrors"
)

type DirParser struct {
	path string
}

func NewDirParser(path string) DirParser {
	return DirParser{path: path}
}

func (p DirParser) Parse() (map[string]*types.Package, error) {

	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, p.path, nil, 0)
	if err != nil {
		return map[string]*types.Package{},
			xerrors.Errorf("failed to parse the dir: %w", err)
	}

	conf := types.Config{Importer: importer.Default()}
	checkedPkgs := map[string]*types.Package{}

	for pkgName, pkg := range pkgs {

		fs := []*ast.File{}
		for _, f := range pkg.Files {
			fs = append(fs, f)
		}

		checked, err := conf.Check(pkgName, fset, fs, nil)
		if err != nil {
			return map[string]*types.Package{},
				xerrors.Errorf("failed to check types: %w", err)
		}

		checkedPkgs[pkgName] = checked
	}

	return checkedPkgs, nil
}
