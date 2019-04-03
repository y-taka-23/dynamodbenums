dynamodbenums
=============

[![CircleCI](https://circleci.com/gh/y-taka-23/dynamodbenums.svg?style=svg)](https://circleci.com/gh/y-taka-23/dynamodbenums)

dynamodbenums is a small utility tool to convert Go enums (a.k.a iota) to DynamoDB strings, vice versa.

Description
-----------

For the specified name of a type defined as constants, dynamodbenums generates a Go source file including marshaling/unmarshaling functions for the given type. The functions make the type as an implementation of the following interfaces using its string repsentation.

```go
type Marshaler interface {
    MarshalDynamoDBAttributeValue(*dynamodb.AttributeValue) error
}
```

```go
type Unmarshaler interface {
    UnmarshalDynamoDBAttributeValue(*dynamodb.AttributeValue) error
}
```

If the type has the `String()` method, it converts the enum values to DynamoDB strings. Otherwise simply their identifiers will be used.

For the detail of the interfaces, see [the official document of AWS SDK for Go](https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/dynamodbattribute/).

Install
-------

```console
$ git clone https://github.com/y-taka-23/dynamodbenums.git
$ cd dynamodbenums
$ go install
```

Walkthrough
-----------

As an example. consider the following simple Go module `github.com/user/repo`.

```
repo
+-- go.mod
+-- main.go
+-- painkiller
    +--- painkiller.go
```

The `painkiller.go` defines the data types you store in DynamoDB.

```go
package painkiller

type Pill int

const (
	Placebo Pill = iota
	Aspirin
	Ibuprofen
	Paracetamol
	Acetaminophen = Paracetamol
)

type Prescription struct {
  ID       string
  Name     Pill
  Quantity int
}
```

And the `main.go` has marshaling procedure for the data types.

```go
package main

import (
  "fmt"
  
  "github.com/user/repo/painkiller"
  
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func main() {

  p := painkiller.Prescription{
    ID:       "PRE-0001",
    Name:     painkiller.Aspirin,
    Quantity: 42,
  }
  
  item, _ := dynamodbattribute.MarshalMap(p)
  fmt.Println(item)
}
```

Firstly run the program without dynamodbenums. The `painkiller.Aspirin` will be displayed just a number `1`.

```console
$ go install -i
$ go run main.go
map[ID:{
  S: "PRE-0001"
} Name:{
  N: "1"
} Quantity:{
  N: "42"
}]
```

Next, generate marshaler/unmarshaler correspong to the `Pill` type by dynamodbenums.

```console
$ dynamodbenums -type=Pill painkiller/
```

You can see a Go source file named `pill_dynamodbenums.go` in the same directory with `painkiller.go`, where the `Pill` defined in.

```
repo
+-- go.mod
+-- main.go
+-- painkiller
    +--- painkiller.go
    +--- pill_dynamodbenums.go
```

Run the program again. In this time, `painkiller.Aspirin` shold be displayed as a string.

```console
$ go run main.go
map[ID:{
  S: "PRE-0001"
} Name:{
  S: "Aspirin"
} Quantity:{
  N: "42"
}]
```

### go generate

dynamodbenums' default behaviour is goodly designed for cooperating with `go generate`. For example you can process code by an annotation comment like:

```go
//go:generate dynamodbenums -type=Pill
type Pill int
```

Options
-------

* `-prefix`: prefix for the output file
* `-suffix`: suffix for the output file (default `_dynamodbenums`)
* `-type`: comma-separated list of type names (required)
