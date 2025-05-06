package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/egoodhall/fsm"
)

func main() {
	input, opts, err := fsm.ParseFlags()
	if err != nil {
		log.Fatalf("parse options: %s", err)
	}

	content, err := os.ReadFile(input)
	if err != nil {
		log.Fatalf("read file: %s", err)
	}

	model, err := fsm.ParseModel(content)
	if err != nil {
		log.Fatalf("parse model: %s", err)
	}

	generated := fsm.Generate(opts.Pkg, model)

	if err := os.MkdirAll(opts.Out, 0755); err != nil {
		log.Fatalf("mkdir: %s", err)
	}

	out := filepath.Join(opts.Out, buildOutfileName(input))
	file, err := os.OpenFile(out, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0640)
	if err != nil {
		log.Fatalf("open file: %s", err)
	}
	defer file.Close()

	if err := generated.Render(file); err != nil {
		log.Fatalf("render file: %s", err)
	}
}

func buildOutfileName(name string) string {
	file := filepath.Base(name)
	ext := filepath.Ext(file)
	fname := strings.TrimSuffix(file, ext)
	return fmt.Sprintf("%s.fsm.go", fname)
}
