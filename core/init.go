package core

func init() {
	gIsDebug = false //调试开关
	gError = NewError()
	gErrorCnt = 1
	gStatementStack = NewStatementStack()
	gSymbols = NewScopedSymbolTable("global", 1, nil)
	gBuiltin = NewBuiltinSymbolTable("builtin")
	gBuiltin.set("print", NewFunc(true, nil, "print", &Param{flag: -1}, nil))
	gBuiltin.set("list", NewFunc(true, nil, "list", &Param{flag: -1}, nil))
	gBuiltin.set("exit", NewFunc(true, nil, "exit", &Param{flag: 0}, nil))

	inScope := NewScopedSymbolTable(gSymbols.scopeName+"_class", gSymbols.scopeLevel+1, gSymbols)
	obj := &Class{scope: inScope,
		name: &Variable{name: "Ojbect"},
		mems: []AstNode{}}
	clsInitFunc := &Func{scope: inScope,
		name:      "Object",
		isBuiltin: false,
		params:    &Param{flag: 0},
		body: &LocalCompoundStatement{
			nodes: []AstNode{&ReturnStatement{results: []AstNode{obj}}}}}
	inScope.set("Object", clsInitFunc)

	gBuiltin.set("Object", obj)
}
