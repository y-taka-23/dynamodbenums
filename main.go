package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/y-taka-23/dynamodbenums/pkg/command"
	"github.com/y-taka-23/dynamodbenums/pkg/extractor"
	"github.com/y-taka-23/dynamodbenums/pkg/parser"
	"github.com/y-taka-23/dynamodbenums/pkg/renderer"
)

var (
	types  = flag.String("type", "", "comma-separated list of type names (required)")
	prefix = flag.String("prefix", "", "prefix for the output file")
	suffix = flag.String("suffix", "_dynamodbenums", "suffix for the output file")
)

func main() {

	flag.Parse()

	if *types == "" {
		log.Fatal("flag -type is required")
	}
	typeNames := strings.Split(*types, ",")

	var input string
	args := flag.Args()

	switch len(args) {
	case 0:
		input = "."
	case 1:
		input = args[0]
	default:
		log.Fatal("specify a single directory/file at once")
	}

	path, err := filepath.Abs(input)
	if err != nil {
		log.Fatalf("cannot resolve the absolute path of %s: %v", input, err)
	}

	stat, err := os.Stat(path)
	if err != nil {
		log.Fatalf("cannot stat %s: %v", path, err)
	}

	var psr command.Parser
	var dir string

	if stat.IsDir() {
		psr = parser.NewDirParser(path)
		dir = path
	} else {
		psr = parser.NewFileParser(path)
		dir = filepath.Dir(path)
	}

	ext := extractor.New()

	rnd, err := renderer.NewTemplateRenderer(renderer.Template, dir, *prefix, *suffix)
	if err != nil {
		log.Fatalf("cannot load the template: %v", err)
	}

	c := command.New(psr, ext, rnd)

	if err := c.Run(typeNames); err != nil {
		log.Fatalf("something wrong during execution: %v", err)
	}
}
