package renderer

const Template = `// generated by dynamodbenums; DO NOT EDIT

package {{ .PackageName }}

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	_{{ .TypeName }}Values = []{{ .TypeName }}{

		{{ range .Values }}{{ . }},
		{{ end }}
	}

	_{{ .TypeName }}Attributes = []string{

		{{ range .Values }}"{{ . }}",
		{{ end }}
	}

	_{{ .TypeName }}ValueToAttribute = map[{{ .TypeName }}]string{}
	_{{ .TypeName }}AttributeToValue = map[string]{{ .TypeName }}{}
)

func init() {
	for i, val := range _{{ .TypeName }}Values {
		attr := _{{ .TypeName }}Attributes[i]
		_{{ .TypeName }}ValueToAttribute[val] = attr
		_{{ .TypeName }}AttributeToValue[attr] = val
	}
}

func (r {{ .TypeName }}) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	n, ok := _{{ .TypeName }}ValueToAttribute[r]
	if !ok {
		return fmt.Errorf("invalid {{ .TypeName }}: %d", r)
	}
	av.S = &n
	return nil
}

func (r *{{ .TypeName }}) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	var n string
	if err := dynamodbattribute.Unmarshal(av, &n); err != nil {
		return err
	}
	v, ok := _{{ .TypeName }}AttributeToValue[n]
	if !ok {
		return fmt.Errorf("invalid {{ .TypeName }}: %q", n)
	}
	*r = v
	return nil
}
`
