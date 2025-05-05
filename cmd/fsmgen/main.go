package fsmgen

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/egoodhall/fsm"
)

func main() {
	input, opts, err := fsm.ParseOptions()
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

	code, err := fsm.Generate(model)
	if err != nil {
		log.Fatalf("generate code: %s", err)
	}

	if err := os.MkdirAll(opts.Out, 0755); err != nil {
		log.Fatalf("mkdir: %s", err)
	}

	if err := os.WriteFile(filepath.Join(opts.Out, fmt.Sprintf("%s.fsm.go", model.Name)), code, 0644); err != nil {
		log.Fatalf("write file: %s", err)
	}
}
