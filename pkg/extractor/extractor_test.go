package extractor_test

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"reflect"
	"testing"

	"github.com/y-taka-23/dynamodbenums/pkg/enum"
	"github.com/y-taka-23/dynamodbenums/pkg/extractor"
)

func TestExtractor_Extract(t *testing.T) {

	ext := extractor.New()

	cases := []struct {
		desc     string
		typeName string
		src      string
		expected enum.Enum
		hasErr   bool
	}{
		{
			desc:     "happy path",
			typeName: "Pill",
			src:      painkiller,
			expected: enum.Enum{
				PackageName: "painkiller",
				TypeName:    "Pill",
				Values: []string{
					"Acetaminophen",
					"Aspirin",
					"Ibuprofen",
					"Paracetamol",
					"Placebo",
				},
			},
			hasErr: false,
		},
		{
			desc:     "placeholders should be skipped",
			typeName: "Pill",
			src:      placeholders,
			expected: enum.Enum{
				PackageName: "painkiller",
				TypeName:    "Pill",
				Values: []string{
					"Acetaminophen",
					"Aspirin",
					"Paracetamol",
				},
			},
			hasErr: false,
		},
		{
			desc:     "var should be ignored",
			typeName: "Pill",
			src:      definedInVar,
			expected: enum.Enum{
				PackageName: "painkiller",
				TypeName:    "Pill",
				Values:      []string{},
			},
			hasErr: false,
		},
		{
			desc:     "constants of the other types should be ignored",
			typeName: "Pill",
			src:      multiTypes,
			expected: enum.Enum{
				PackageName: "painkiller",
				TypeName:    "Pill",
				Values: []string{
					"Aspirin",
					"Placebo",
				},
			},
			hasErr: false,
		},
		{
			desc:     "the identifier not found",
			typeName: "Granule",
			src:      painkiller,
			expected: enum.Enum{
				PackageName: "painkiller",
				TypeName:    "Granule",
				Values:      []string{},
			},
			hasErr: false,
		},
		{
			desc:     "not a type but a function",
			typeName: "Pill",
			src:      definedAsFunc,
			expected: enum.Enum{
				PackageName: "painkiller",
				TypeName:    "Pill",
				Values:      []string{},
			},
			hasErr: false,
		},
		{
			desc:     "type name cannot be empty",
			typeName: "",
			src:      painkiller,
			expected: enum.Enum{},
			hasErr:   true,
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {

			checked := parseSrc(c.src)
			actual, err := ext.Extract(c.typeName, checked)

			switch {
			case err != nil && !c.hasErr:
				t.Error("unexpected error: ", err)
			case err == nil && c.hasErr:
				t.Error("expected error has not occurred")
			case !reflect.DeepEqual(actual, c.expected):
				t.Errorf("want %+v got %+v", c.expected, actual)
			}

		})
	}
}

func parseSrc(src string) *types.Package {

	fset := token.NewFileSet()

	f, _ := parser.ParseFile(fset, "", src, 0)

	conf := types.Config{Importer: importer.Default()}
	checked, _ := conf.Check(f.Name.Name, fset, []*ast.File{f}, nil)

	return checked
}

const painkiller = `package painkiller

import (
	"fmt"
)

type Pill int

const (
	Placebo Pill = iota
	Aspirin
	Ibuprofen
	Paracetamol
	Acetaminophen = Paracetamol
)
`

const placeholders = `package painkiller

type Pill int

const (
	_ Pill = iota
	Aspirin
	_
	Paracetamol
	Acetaminophen = Paracetamol
)
`

const definedInVar = `package painkiller

type Pill int

var (
	Placebo       Pill = 0
	Aspirin       Pill = 1
	Ibuprofen     Pill = 2
	Paracetamol   Pill = 3
	Acetaminophen      = Paracetamol
)
`

const multiTypes = `package painkiller

type Pill int
type Granule int

const (
	Placebo Pill = iota
	Aspirin
	Ibuprofen Granule = iota
	Paracetamol
	Acetaminophen = Paracetamol
)
`

const definedAsFunc = `package painkiller

func Pill(n int) error {
	return nil
}
`
