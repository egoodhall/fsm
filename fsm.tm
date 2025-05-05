language fsm(go);

eventBased = true
lang = "go"
package = "github.com/egoodhall/fsm/gen/parser"

:: lexer

%x initial,inFsmDefinition,inTypeDeclaration,inOption;

ident = /[a-zA-Z_][a-zA-Z0-9_]+/

stringLiteral = /"([^"]|\\")*"/
intLiteral = /-?[0-9]+/
floatLiteral = /-?[0-9]+\.[0-9]+/
boolLiteral = /true|false/

<initial> {
  'fsm': /fsm/ { l.State = StateInFsmDefinition }
  'type': /type/ { l.State = StateInTypeDeclaration }
  'option': /option/ { l.State = StateInOption }
}

<*> {
  Whitespace: /[ \t\n\r]+/ (space)
  EolComment: /\/\/[^\n]*\n/ (space)
  BlockComment: /\/\*([^*]|\*+[^*\/])*\**\*\// (space)
  error:
}

<inFsmDefinition> {
  Name: /{ident}/ (class)
  '[': /\[/
  ']': /\]/
  '{': /\{/
  '}': /\}/ { l.State = StateInitial }
  ';': /;/
  'start': /start/
  'state': /state/
  'end': /end/
  'transition': /transition/
  'to': /to/
  'or': /or/
}

<inTypeDeclaration> {
  Name: /{ident}/
  ';': /;/ { l.State = StateInitial }
}

<inOption> {
  Name: /({ident}\.)?{ident}/
  '=': /=/
  StringLiteral: /{stringLiteral}/
  BoolLiteral: /{boolLiteral}/ 1
  IntLiteral: /{intLiteral}/
  FloatLiteral: /{floatLiteral}/
  ';': /;/ { l.State = StateInitial }
}

:: parser

%input FsmFile;

# Option handling
OptionName -> OptionName: Name;
OptionString -> OptionString: StringLiteral;
OptionBool -> OptionBool: BoolLiteral;
OptionInt -> OptionInt: IntLiteral;
OptionFloat -> OptionFloat: FloatLiteral;
OptionValue: OptionString | OptionBool | OptionInt | OptionFloat;
Option: 'option' OptionName '=' OptionValue ';';

# Type declarations
TypeName -> TypeName: Name;
TypeDeclaration: 'type' TypeName ';';

# FSM definitions
FsmName -> FsmName: Name;
InputType -> InputType: Name;

# State declarations
StateName -> StateName: Name;
ContextType -> ContextType: Name;
StateContext: '[' ContextType ']';

StartState: 'start' StateName ';';
RegularState: 'state' StateName StateContext? ';';
EndState: 'end' StateName ';';
StateDeclaration: StartState | RegularState | EndState;

# Transition declarations
SourceState -> SourceState: Name;
TargetState -> TargetState: Name;
TargetStateRest: 'or' TargetState TargetStateRest?;
TargetStateList: TargetState TargetStateRest?;
TransitionDeclaration: 'transition' SourceState 'to' TargetStateList ';';

# FSM body content
FsmBodyItem: StateDeclaration | TransitionDeclaration;
FsmBodyRest: FsmBodyItem FsmBodyRest?;
FsmBody: FsmBodyRest?;

FsmInputType: '[' InputType ']';
FsmDefinition: 'fsm' FsmName FsmInputType? '{' FsmBody '}';

# Top-level elements
Definition: TypeDeclaration | FsmDefinition | Option | error (';'|'}');
Definitions: Definition Definitions?;

# File definition
FsmFile: Definitions?;
