package fsm

import (
	"errors"
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v3"
)

type FsmModel struct {
	Name   string               `yaml:"name"`
	Types  map[string]TypeModel `yaml:"types"`
	States []StateModel         `yaml:"states"`
}

func (s *FsmModel) InitialState() StateModel {
	for _, state := range s.States {
		if state.Entrypoint {
			return state
		}
	}
	panic("no initial state found")
}
func (s *FsmModel) GetState(name State) StateModel {
	for _, state := range s.States {
		if state.Name == name {
			return state
		}
	}
	panic(fmt.Sprintf("state %s not found", name))
}

func (s *FsmModel) FsmName() string {
	return strcase.ToCamel(s.Name) + "FSM"
}

func (s *FsmModel) StateTypeName() string {
	return strcase.ToCamel(s.Name) + "State"
}

func (s *FsmModel) StateName(state StateModel) string {
	return s.StateTypeName() + strcase.ToCamel(string(state.Name))
}

func (s *FsmModel) FsmInternalName() string {
	return strcase.ToLowerCamel(s.Name) + "FSM"
}

func (s *FsmModel) FsmBuilderConstructorName() string {
	return "New" + s.FsmBuilderName()
}

func (s *FsmModel) FsmBuilderName() string {
	return s.FsmName() + "Builder"
}

func (s *FsmModel) FsmBuilderStageName(state StateModel) string {
	return fmt.Sprintf("%s_%sStage", s.FsmBuilderName(), strcase.ToCamel(string(state.Name)))
}

func (s *FsmModel) FsmBuilderStageMethodName(state StateModel) string {
	return fmt.Sprintf("%sState", strcase.ToCamel(string(state.Name)))
}

func (s *FsmModel) FsmStateMessageName(state StateModel) string {
	return strcase.ToLowerCamel(string(state.Name)) + "Params"
}

func (s *FsmModel) FsmStateInternalName(state StateModel) string {
	return strcase.ToLowerCamel(string(state.Name)) + "State"
}

func (s *FsmModel) FsmStateQueueInternalName(state StateModel) string {
	return strcase.ToLowerCamel(string(state.Name)) + "Queue"
}

func (s *FsmModel) FsmStateProcessorName(state StateModel) string {
	return strcase.ToLowerCamel(string(state.Name)) + "Processor"
}

func (s *FsmModel) FsmBuilderFinalStageName() string {
	return fmt.Sprintf("%s__FinalStage", s.FsmBuilderName())
}

func (s *FsmModel) RenderType(name string) jen.Code {
	def, ok := s.Types[name]
	if !ok || def.Package == "" {
		return jen.Id(name)
	}
	return jen.Qual(def.Package, def.Type)
}

func (s *FsmModel) TransitionToName(to State) string {
	return fmt.Sprintf("To%s", strcase.ToCamel(string(to)))
}

func (s *FsmModel) TransitionsParamTypeName(state StateModel) string {
	return strcase.ToCamel(string(state.Name)) + "Transitions"
}

type TypeModel struct {
	Type    string `yaml:"type"`
	Package string `yaml:"package,omitempty"`
}

type StateModel struct {
	Name        State    `yaml:"name"`
	Entrypoint  bool     `yaml:"entrypoint"`
	Terminal    bool     `yaml:"terminal"`
	Workers     int      `yaml:"workers"`
	Queue       int      `yaml:"queue"`
	Inputs      []string `yaml:"inputs"`
	Transitions []State  `yaml:"transitions"`
}

func ParseModel(p []byte) (*FsmModel, error) {
	var model FsmModel
	if err := yaml.Unmarshal(p, &model); err != nil {
		return nil, err
	}
	if err := validateModel(&model); err != nil {
		return nil, err
	}
	return &model, nil
}

func validateModel(model *FsmModel) error {
	if model.Name == "" {
		return errors.New("name is required")
	}
	var entrypoints, terminals int
	for _, state := range model.States {
		if state.Name == "" {
			return errors.New("state name is required")
		}
		if state.Terminal && len(state.Transitions) > 0 {
			return errors.New("terminal state cannot have transitions")
		}
		if state.Terminal && len(state.Inputs) > 0 {
			return errors.New("terminal state cannot have inputs")
		}
		if state.Workers == 0 {
			state.Workers = 1
		}
		if state.Workers < 0 {
			return errors.New("each state must have at least one worker")
		}
		if state.Queue < 0 {
			state.Queue = 16
		}
		if state.Entrypoint {
			entrypoints++
		}
		if state.Terminal {
			terminals++
		}
	}
	if entrypoints != 1 {
		return errors.New("exactly one entrypoint is required")
	}
	if terminals < 1 {
		return errors.New("at least one terminal state is required")
	}
	return nil
}
