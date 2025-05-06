package fsm

import (
	"fmt"

	"github.com/dave/jennifer/jen"
)

func Generate(pkg string, model *FsmModel) *jen.File {
	file := jen.NewFile(pkg)

	// Public interfaces
	for _, c := range generatePublicInterfaces(model) {
		file.Add(c).Line()
	}

	file.Comment("FSM type checks")
	file.Var().Id("_").Id(model.FsmName()).Op("=").New(jen.Id(model.FsmInternalName()))
	file.Var().Id("_").Qual("github.com/egoodhall/fsm", "SupportsOptions").Op("=").New(jen.Id(model.FsmInternalName()))
	for _, state := range model.States {
		file.Var().Id("_").Id(model.FsmBuilderStageName(state)).Op("=").New(jen.Id(model.FsmInternalName()))
	}
	file.Var().Id("_").Id(model.FsmBuilderFinalStageName()).Op("=").New(jen.Id(model.FsmInternalName()))
	for _, state := range model.States {
		if !state.Terminal {
			file.Var().Id("_").Id(model.TransitionsParamTypeName(state)).Op("=").New(jen.Id(model.FsmInternalName()))
		}
	}

	file.Line().Commentf("%s implementation", model.FsmName())

	// FSM implementation
	for _, c := range generateFSMImplementation(model) {
		file.Add(c).Line()
	}

	return file
}

func generatePublicInterfaces(model *FsmModel) []jen.Code {
	code := make([]jen.Code, 0)

	code = append(code, jen.Type().Id(model.StateTypeName()).Qual("github.com/egoodhall/fsm", "State"))

	code = append(code, jen.Const().DefsFunc(func(g *jen.Group) {
		for _, state := range model.States {
			g.Id(model.StateName(state)).Id(model.StateTypeName()).Op("=").Lit(string(state.Name))
		}
	}))

	// FSM interface
	code = append(code, jen.Type().Id(model.FsmName()).Interface(
		jen.Qual("github.com/egoodhall/fsm", "SupportsOptions"),
		jen.Id("Submit").
			ParamsFunc(func(g *jen.Group) {
				g.Id("ctx").Qual("context", "Context")
				for _, param := range model.InitialState().Inputs {
					g.Id(param).Add(model.RenderType(param))
				}
			}).
			Params(jen.Qual("github.com/egoodhall/fsm", "TaskID"), jen.Error()),
	))

	// FSM builder constructor
	code = append(code, jen.Func().Id(model.FsmBuilderConstructorName()).Params().
		Id(model.FsmBuilderStageName(model.InitialState())).Block(
		jen.Return(jen.New(jen.Id(model.FsmInternalName()))),
	))

	// FSM transition interfaces
	for _, state := range model.States {
		if state.Terminal {
			continue
		}

		code = append(code, jen.Type().Id(model.TransitionsParamTypeName(state)).InterfaceFunc(func(g *jen.Group) {
			for _, transition := range state.Transitions {
				params := []jen.Code{jen.Qual("context", "Context")}
				for _, param := range model.GetState(transition).Inputs {
					params = append(params, model.RenderType(param))
				}
				g.Id(model.TransitionToName(transition)).Params(params...).Error()
			}
		}))
	}

	// FSM builder stage interfaces
	for i, state := range model.States {
		method := jen.Id(model.FsmBuilderStageMethodName(state)).Params(
			generateFSMStateMethodSignature(model, state),
		)
		if i == len(model.States)-1 {
			method = method.Id(model.FsmBuilderFinalStageName())
		} else {
			method = method.Id(model.FsmBuilderStageName(model.States[i+1]))
		}

		code = append(code, jen.Type().Id(model.FsmBuilderStageName(state)).Interface(method).Line())
	}

	// FSM builder final stage
	code = append(code, jen.Type().Id(model.FsmBuilderName()+"__FinalStage").Interface(
		jen.Id("BuildAndStart").Params(jen.Qual("context", "Context"), jen.Op("...").Qual("github.com/egoodhall/fsm", "Option")).Params(jen.Id(model.FsmName()), jen.Error()),
	))

	return code
}

func generateFSMImplementation(model *FsmModel) []jen.Code {
	code := make([]jen.Code, 0)

	for _, state := range model.States {
		code = append(code, jen.Type().Id(model.FsmStateMessageName(state)).StructFunc(func(g *jen.Group) {
			g.Id("ID").Qual("github.com/egoodhall/fsm", "TaskID")
			for i, input := range state.Inputs {
				g.Id(fmt.Sprintf("P%d", i)).Add(model.RenderType(input))
			}
		}))
	}

	// FSM struct
	code = append(code,
		jen.Type().Id(model.FsmInternalName()).StructFunc(func(g *jen.Group) {
			g.Id("lock").Qual("sync", "Mutex")
			g.Id("ctx").Qual("context", "Context")
			g.Line()
			g.Comment("Configuration options")
			g.Id("store").Qual("github.com/egoodhall/fsm", "Store")
			g.Id("logger").Op("*").Qual("log/slog", "Logger")
			g.Id("onTransition").Qual("github.com/egoodhall/fsm", "TransitionListener")
			g.Id("onCompletion").Qual("github.com/egoodhall/fsm", "CompletionListener")
			g.Line()
			g.Comment("FSM state transitions")
			for _, state := range model.States {
				g.Id(model.FsmStateInternalName(state)).Add(generateFSMStateMethodSignature(model, state))
			}
			g.Line()
			g.Comment("FSM queues")
			for _, state := range model.States {
				g.Id(model.FsmStateQueueInternalName(state)).Chan().Id(model.FsmStateMessageName(state))
			}
		}),
		jen.Comment("FSM builder methods"),
	)

	// FSM builder stage methods
	for i, state := range model.States {
		method := jen.Func().
			Params(jen.Id("f").Op("*").Id(model.FsmInternalName())).
			Id(model.FsmBuilderStageMethodName(state)).Params(
			jen.Id("fn").Add(generateFSMStateMethodSignature(model, state)),
		)

		if i == len(model.States)-1 {
			method = method.Id(model.FsmBuilderFinalStageName())
		} else {
			method = method.Id(model.FsmBuilderStageName(model.States[i+1]))
		}

		code = append(code,
			method.Block(
				jen.Id("f").Dot(model.FsmStateQueueInternalName(state)).Op("=").Make(jen.Chan().Id(model.FsmStateMessageName(state)), jen.Lit(100)),
				jen.Id("f").Dot(model.FsmStateInternalName(state)).Op("=").Id("fn"),
				jen.Return(jen.Id("f")),
			),
		)
	}

	// FSM builder final stage
	code = append(code,
		jen.Func().
			Params(jen.Id("f").Op("*").Id(model.FsmInternalName())).
			Id("BuildAndStart").
			Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("opts").Op("...").Qual("github.com/egoodhall/fsm", "Option")).
			Params(jen.Id(model.FsmName()), jen.Error()).
			BlockFunc(func(g *jen.Group) {
				// Check if FSM is already started
				g.Comment("Check if FSM is already started")
				g.If(jen.Op("!").Id("f").Dot("lock").Dot("TryLock").Call()).Block(
					jen.Return(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("FSM already started"))),
				)
				g.Line()
				g.Comment("Set context")
				g.Id("f").Dot("ctx").Op("=").Id("ctx")
				g.Line()
				// Apply options
				g.Comment("Apply options")
				g.For(jen.List(jen.Id("_"), jen.Id("opt")).Op(":=").Range().Id("opts")).Block(
					jen.If(jen.Err().Op(":=").Id("opt").Call(jen.Id("f")), jen.Err().Op("!=").Nil()).Block(
						jen.Return(jen.Nil(), jen.Err()),
					),
				)
				g.If(jen.Id("f").Dot("store").Op("==").Nil()).Block(
					jen.If(jen.Err().Op(":=").Qual("github.com/egoodhall/fsm", "InMemory").Call().Call(jen.Id("f")), jen.Err().Op("!=").Nil()).Block(
						jen.Return(jen.Nil(), jen.Err()),
					),
				)
				g.Line()
				g.Comment("Start FSM processors")
				for _, state := range model.States {
					g.Go().Id("f").Dot(model.FsmStateProcessorName(state)).Call()
				}
				g.Line()
				// Return FSM
				g.Return(jen.Id("f"), jen.Nil())
			}),
	)

	// FSM option methods
	code = append(code,
		jen.Comment("FSM options"),
		jen.Func().
			Params(jen.Id("f").Op("*").Id(model.FsmInternalName())).
			Id("WithStore").
			Params(jen.Id("store").Qual("github.com/egoodhall/fsm", "Store")).
			Block(
				jen.Id("f").Dot("store").Op("=").Id("store"),
			),
		jen.Func().
			Params(jen.Id("f").Op("*").Id(model.FsmInternalName())).
			Id("WithLogger").
			Params(jen.Id("logger").Op("*").Qual("log/slog", "Logger")).
			Block(
				jen.Id("f").Dot("logger").Op("=").Id("logger"),
			),
		jen.Func().
			Params(jen.Id("f").Op("*").Id(model.FsmInternalName())).
			Id("WithTransitionListener").
			Params(jen.Id("listener").Qual("github.com/egoodhall/fsm", "TransitionListener")).
			Block(
				jen.Id("f").Dot("onTransition").Op("=").Id("listener"),
			),
		jen.Func().
			Params(jen.Id("f").Op("*").Id(model.FsmInternalName())).
			Id("WithCompletionListener").
			Params(jen.Id("listener").Qual("github.com/egoodhall/fsm", "CompletionListener")).
			Block(
				jen.Id("f").Dot("onCompletion").Op("=").Id("listener"),
			),
	)

	// FSM transition methods
	code = append(code, jen.Comment("FSM transition methods"))
	for _, state := range model.States {
		code = append(code,
			jen.Func().
				Params(jen.Id("f").Op("*").Id(model.FsmInternalName())).
				Id(model.TransitionToName(state.Name)).
				ParamsFunc(func(g *jen.Group) {
					g.Id("ctx").Qual("context", "Context")
					for i, param := range state.Inputs {
						g.Id(fmt.Sprintf("P%d", i)).Add(model.RenderType(param))
					}
				}).
				Error().
				Block(
					jen.Id("id").Op(":=").Qual("github.com/egoodhall/fsm", "GetTaskID").Call(jen.Id("ctx")),
					jen.Id("msg").Op(":=").Id(model.FsmStateMessageName(state)).ValuesFunc(func(g *jen.Group) {
						g.Id("ID").Op(":").Id("id")
						for i := range state.Inputs {
							g.Id(fmt.Sprintf("P%d", i)).Op(":").Id(fmt.Sprintf("P%d", i))
						}
					}),
					jen.Line(),
					// Encode and save transition
					jen.Id("buf").Op(":=").New(jen.Qual("bytes", "Buffer")),
					jen.If(jen.Err().Op(":=").Qual("encoding/gob", "NewEncoder").Call(jen.Id("buf")).Dot("Encode").Call(jen.Id("msg")), jen.Err().Op("!=").Nil()).Block(
						jen.Return(jen.Err()),
					),
					jen.Line(),
					jen.If(jen.Err().Op(":=").Id("f").Dot("store").Dot("Q").Call().Dot("RecordTransition").Call(
						jen.Line().Id("ctx"),
						jen.Line().Int64().Call(jen.Id("id")),
						jen.Line().String().Call(jen.Qual("github.com/egoodhall/fsm", "GetState").Call(jen.Id("ctx"))),
						jen.Line().String().Call(jen.Id(model.StateName(state))),
						jen.Line().Id("buf").Dot("Bytes").Call(),
					), jen.Err().Op("!=").Nil()).Block(
						jen.Return(jen.Err()),
					),
					jen.Line(),
					jen.Select().Block(
						jen.Case(jen.Id("f").Dot(model.FsmStateQueueInternalName(state)).Op("<-").Id("msg")).Block(
							jen.Return(jen.Nil()),
						),
						jen.Case(jen.Op("<-").Id("ctx").Dot("Done").Call()).Block(
							jen.Return(jen.Qual("errors", "New").Call(jen.Lit("task submission cancelled"))),
						),
					),
				),
		)
	}

	// FSM processing methods
	for _, state := range model.States {
		code = append(code,
			jen.Func().
				Params(jen.Id("f").Op("*").Id(model.FsmInternalName())).
				Id(model.FsmStateProcessorName(state)).
				Params().
				Block(
					jen.Id("ctx").Op(":=").Qual("github.com/egoodhall/fsm", "PutState").Call(jen.Id("f").Dot("ctx"), jen.Qual("github.com/egoodhall/fsm", "State").Call(jen.Id(model.StateName(state)))),
					jen.For(jen.Id("msg").Op(":=").Range().Id("f").Dot(model.FsmStateQueueInternalName(state))).Block(
						jen.Id("f").Dot(model.FsmStateInternalName(state)).CallFunc(func(g *jen.Group) {
							g.Qual("github.com/egoodhall/fsm", "PutTaskID").Call(jen.Id("ctx"), jen.Id("msg").Dot("ID"))
							if !state.Terminal {
								g.Id("f")
							}
							for i := range state.Inputs {
								g.Id("msg").Dot(fmt.Sprintf("P%d", i))
							}
						}),
					),
				),
		)
	}

	// FSM submit method
	code = append(code,
		jen.Comment("Submit FSM tasks"),
		jen.Func().
			Params(jen.Id("f").Op("*").Id(model.FsmInternalName())).
			Id("Submit").
			ParamsFunc(func(g *jen.Group) {
				g.Id("ctx").Qual("context", "Context")
				for i, param := range model.InitialState().Inputs {
					g.Id(fmt.Sprintf("P%d", i)).Add(model.RenderType(param))
				}
			}).
			Params(jen.Qual("github.com/egoodhall/fsm", "TaskID"), jen.Error()).
			Block(
				// Construct message without ID
				jen.Id("msg").Op(":=").Id(model.FsmStateMessageName(model.InitialState())).ValuesFunc(func(g *jen.Group) {
					for i := range model.InitialState().Inputs {
						g.Id(fmt.Sprintf("P%d", i)).Op(":").Id(fmt.Sprintf("P%d", i))
					}
				}),
				jen.Line(),
				// Encode and create task
				jen.Id("buf").Op(":=").New(jen.Qual("bytes", "Buffer")),
				jen.If(jen.Err().Op(":=").Qual("encoding/gob", "NewEncoder").Call(jen.Id("buf")).Dot("Encode").Call(jen.Id("msg")), jen.Err().Op("!=").Nil()).Block(
					jen.Return(jen.Lit(0), jen.Err()),
				),
				jen.Line(),
				jen.List(jen.Id("task"), jen.Id("err")).Op(":=").Id("f").Dot("store").Dot("Q").Call().Dot("CreateTask").Call(jen.Id("ctx"), jen.Id("buf").Dot("Bytes").Call()),
				jen.If(jen.Id("err").Op("!=").Nil()).Block(
					jen.Return(jen.Lit(0), jen.Err()),
				),
				jen.Id("msg").Dot("ID").Op("=").Qual("github.com/egoodhall/fsm", "TaskID").Call(jen.Id("task").Dot("ID")),
				jen.Line(),
				jen.Select().Block(
					jen.Case(jen.Id("f").Dot(model.FsmStateQueueInternalName(model.InitialState())).Op("<-").Id("msg")).Block(
						jen.Return(jen.Id("msg").Dot("ID"), jen.Nil()),
					),
					jen.Case(jen.Op("<-").Id("ctx").Dot("Done").Call()).Block(
						jen.Return(jen.Lit(0), jen.Qual("errors", "New").Call(jen.Lit("task submission cancelled"))),
					),
				),
			),
	)

	return code
}

func generateFSMStateMethodSignature(model *FsmModel, state StateModel) jen.Code {
	params := []jen.Code{
		jen.Qual("context", "Context"),
	}
	if !state.Terminal {
		params = append(params, jen.Id(model.TransitionsParamTypeName(state)))
	}
	for _, param := range state.Inputs {
		params = append(params, model.RenderType(param))
	}
	return jen.Func().Params(params...).Error()
}
