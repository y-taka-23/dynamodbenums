package extractor

import (
	"go/types"
	"sort"

	"golang.org/x/xerrors"

	"github.com/y-taka-23/dynamodbenums/pkg/enum"
)

type Extractor struct{}

func New() Extractor {
	return Extractor{}
}

func (e Extractor) Extract(typeName string, pkg *types.Package) (enum.Enum, error) {

	if typeName == "" {
		return enum.Enum{}, xerrors.New("type name cannot be empty")
	}

	scope := pkg.Scope()
	names := scope.Names()

	var typ types.Type
	consts := []*types.Const{}

	for _, name := range names {

		obj := scope.Lookup(name)

		if t, ok := obj.(*types.TypeName); ok && t.Name() == typeName {
			typ = t.Type()
		}

		if c, ok := obj.(*types.Const); ok {
			consts = append(consts, c)
		}
	}

	if typ == nil {
		return enum.Enum{
			PackageName: pkg.Name(),
			TypeName:    typeName,
			Values:      []string{},
		}, nil
	}

	values := []string{}

	for _, c := range consts {
		if c.Type() == typ {
			values = append(values, c.Name())
		}
	}

	sort.Strings(values)

	return enum.Enum{
		PackageName: pkg.Name(),
		TypeName:    typeName,
		Values:      values,
	}, nil
}
