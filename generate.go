package fsm

func Generate(pkg string, model *FSMSchema) ([]byte, error) {
	// TODO: Implement this. Some notes:
	// - The generated code should use the desired package.
	// - The generated code should be named after the schema's name.
	// - The generated code should have a staged builder interface, which satisfies:
	//   - Each state should have a method in the builder that takes a context, an object for building transitions to the next state, and the input(s) for the current state.
	//   - The builder should have a Build method that returns the FSM.
	//   - An example of desired usage can be found in example/main.go
	// - The generated FSM should be named {SchemaName}FSM
	// - The generated FSM should have a method Submit{StateName} for each initial state that takes a context and the input(s) for the current state.
	return nil, nil
}
