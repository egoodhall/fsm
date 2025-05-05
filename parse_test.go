package fsm

import (
	"testing"
)

func TestParseFsmDefinition(t *testing.T) {
	src := `
		fsm TrafficLight {
			start Red;
			state Yellow;
			end Green;
			
			transition Green to Yellow;
			transition Yellow to Red;
			transition Red to Green;
		}
	`

	ast, err := Parse(src)
	if err != nil {
		t.Fatalf("Failed to parse FSM: %v", err)
	}

	// Verify the AST structure
	if len(ast.Definitions) != 1 {
		t.Fatalf("Expected 1 definition, got %d", len(ast.Definitions))
	}

	fsm, ok := ast.Definitions[0].(*FsmDefinition)
	if !ok {
		t.Fatalf("Expected FsmDefinition, got %T", ast.Definitions[0])
	}

	if fsm.Name != "TrafficLight" {
		t.Errorf("Expected FSM name 'TrafficLight', got '%s'", fsm.Name)
	}

	if len(fsm.Body) != 6 {
		t.Fatalf("Expected 6 body items (3 states + 3 transitions), got %d", len(fsm.Body))
	}

	// Check start state
	startState, ok := fsm.Body[0].(*StartState)
	if !ok {
		t.Fatalf("Expected StartState, got %T", fsm.Body[0])
	}
	if startState.Name != "Red" {
		t.Errorf("Expected start state name 'Red', got '%s'", startState.Name)
	}

	// Check regular state
	regularState, ok := fsm.Body[1].(*RegularState)
	if !ok {
		t.Fatalf("Expected RegularState, got %T", fsm.Body[1])
	}
	if regularState.Name != "Yellow" {
		t.Errorf("Expected regular state name 'Yellow', got '%s'", regularState.Name)
	}

	// Check end state
	endState, ok := fsm.Body[2].(*EndState)
	if !ok {
		t.Fatalf("Expected EndState, got %T", fsm.Body[2])
	}
	if endState.Name != "Green" {
		t.Errorf("Expected end state name 'Green', got '%s'", endState.Name)
	}

	// Check first transition
	transition, ok := fsm.Body[3].(*TransitionDeclaration)
	if !ok {
		t.Fatalf("Expected TransitionDeclaration, got %T", fsm.Body[3])
	}
	if transition.SourceState != "Red" {
		t.Errorf("Expected source state 'Red', got '%s'", transition.SourceState)
	}
	if len(transition.TargetStates) != 1 || transition.TargetStates[0] != "Yellow" {
		t.Errorf("Expected target state 'Yellow', got %v", transition.TargetStates)
	}
}

func TestParseTypeAndOptions(t *testing.T) {
	src := `
		type Event;
		option lang = "go";
		option package = "github.com/egoodhall/fsm/gen/parser";
		option eventBased = true;
		option timeout = 3.14;
	`

	ast, err := Parse(src)
	if err != nil {
		t.Fatalf("Failed to parse FSM: %v", err)
	}

	// Verify the AST structure
	if len(ast.Definitions) != 5 {
		t.Fatalf("Expected 5 definitions, got %d", len(ast.Definitions))
	}

	// Check type declaration
	typeDef, ok := ast.Definitions[0].(*TypeDeclaration)
	if !ok {
		t.Fatalf("Expected TypeDeclaration, got %T", ast.Definitions[0])
	}
	if typeDef.Name != "Event" {
		t.Errorf("Expected type name 'Event', got '%s'", typeDef.Name)
	}

	// Check string option
	option1, ok := ast.Definitions[1].(*Option)
	if !ok {
		t.Fatalf("Expected Option, got %T", ast.Definitions[1])
	}
	if option1.Name != "lang" {
		t.Errorf("Expected option name 'lang', got '%s'", option1.Name)
	}
	stringOpt, ok := option1.Value.(*StringOption)
	if !ok {
		t.Fatalf("Expected StringOption, got %T", option1.Value)
	}
	if stringOpt.Value != "go" {
		t.Errorf("Expected option value 'go', got '%s'", stringOpt.Value)
	}

	// Check bool option
	option3, ok := ast.Definitions[3].(*Option)
	if !ok {
		t.Fatalf("Expected Option, got %T", ast.Definitions[3])
	}
	if option3.Name != "eventBased" {
		t.Errorf("Expected option name 'eventBased', got '%s'", option3.Name)
	}
	boolOpt, ok := option3.Value.(*BoolOption)
	if !ok {
		t.Fatalf("Expected BoolOption, got %T", option3.Value)
	}
	if boolOpt.Value != true {
		t.Errorf("Expected option value 'true', got '%v'", boolOpt.Value)
	}

	// Check float option
	option4, ok := ast.Definitions[4].(*Option)
	if !ok {
		t.Fatalf("Expected Option, got %T", ast.Definitions[4])
	}
	if option4.Name != "timeout" {
		t.Errorf("Expected option name 'timeout', got '%s'", option4.Name)
	}
	floatOpt, ok := option4.Value.(*FloatOption)
	if !ok {
		t.Fatalf("Expected FloatOption, got %T", option4.Value)
	}
	if floatOpt.Value != 3.14 {
		t.Errorf("Expected option value '3.14', got '%v'", floatOpt.Value)
	}
}
