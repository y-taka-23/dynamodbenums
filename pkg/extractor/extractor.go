package extractor

import (
	"go/types"

	"github.com/y-taka-23/dynamodbenums/pkg/enum"
)

type Extractor struct{}

func New() Extractor {
	return Extractor{}
}

func (e Extractor) Extract(typeName string, pkg *types.Package) (enum.Enum, error) {
	return enum.Enum{}, nil
}
