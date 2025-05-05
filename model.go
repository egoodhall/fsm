package fsm

import (
	"errors"

	"gopkg.in/yaml.v3"
)

type Model struct {
	Name   string  `yaml:"name"`
	States []State `yaml:"states"`
}

type State struct {
	Name        string   `yaml:"name"`
	Terminal    bool     `yaml:"terminal"`
	Inputs      []string `yaml:"inputs"`
	Transitions []string `yaml:"transitions"`
}

func ParseModel(p []byte) (*Model, error) {
	var model Model
	if err := yaml.Unmarshal(p, &model); err != nil {
		return nil, err
	}
	return &model, nil
}

func ValidateModel(model *Model) error {
	if model.Name == "" {
		return errors.New("name is required")
	}
	for _, state := range model.States {
		if state.Name == "" {
			return errors.New("state name is required")
		}
		if state.Terminal && len(state.Transitions) > 0 {
			return errors.New("terminal state cannot have transitions")
		}
	}
	return nil
}
