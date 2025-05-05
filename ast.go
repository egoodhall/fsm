package fsm

// Position represents a position in the source code.
type Position struct {
	Offset    int
	EndOffset int
	Line      int
}

// Node is the interface implemented by all AST nodes.
type Node interface {
	Pos() Position
}

// BaseNode provides common position information for all AST nodes.
type BaseNode struct {
	Position Position
}

func (n BaseNode) Pos() Position {
	return n.Position
}

// File represents a complete FSM file with all its definitions.
type File struct {
	BaseNode
	Definitions []Definition
}

// Definition is an interface for top-level definitions in the FSM file.
type Definition interface {
	Node
	isDefinition()
}

// TypeDeclaration represents a type declaration in the FSM file.
type TypeDeclaration struct {
	BaseNode
	Name string
}

func (TypeDeclaration) isDefinition() {}

// OptionValue represents a value that can be assigned to an option.
type OptionValue interface {
	Node
	isOptionValue()
}

// StringOption represents a string option value.
type StringOption struct {
	BaseNode
	Value string
}

func (StringOption) isOptionValue() {}

// BoolOption represents a boolean option value.
type BoolOption struct {
	BaseNode
	Value bool
}

func (BoolOption) isOptionValue() {}

// IntOption represents an integer option value.
type IntOption struct {
	BaseNode
	Value int64
}

func (IntOption) isOptionValue() {}

// FloatOption represents a floating-point option value.
type FloatOption struct {
	BaseNode
	Value float64
}

func (FloatOption) isOptionValue() {}

// Option represents an option declaration in the FSM file.
type Option struct {
	BaseNode
	Name  string
	Value OptionValue
}

func (Option) isDefinition() {}

// StateDeclaration is an interface for different kinds of state declarations.
type StateDeclaration interface {
	Node
	isStateDeclaration()
}

// StartState represents a start state declaration.
type StartState struct {
	BaseNode
	Name string
}

func (StartState) isStateDeclaration() {}

// RegularState represents a regular state declaration.
type RegularState struct {
	BaseNode
	Name        string
	ContextType string // Optional, empty if not specified
}

func (RegularState) isStateDeclaration() {}

// EndState represents an end state declaration.
type EndState struct {
	BaseNode
	Name string
}

func (EndState) isStateDeclaration() {}

// TransitionDeclaration represents a transition declaration between states.
type TransitionDeclaration struct {
	BaseNode
	SourceState  string
	TargetStates []string
}

// FsmBodyItem is an interface for items in the FSM body.
type FsmBodyItem interface {
	Node
	isFsmBodyItem()
}

func (StartState) isFsmBodyItem()            {}
func (RegularState) isFsmBodyItem()          {}
func (EndState) isFsmBodyItem()              {}
func (TransitionDeclaration) isFsmBodyItem() {}

// FsmDefinition represents an FSM definition.
type FsmDefinition struct {
	BaseNode
	Name      string
	InputType string // Optional, empty if not specified
	Body      []FsmBodyItem
}

func (FsmDefinition) isDefinition() {}
