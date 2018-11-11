package core

func init() {
	g_is_debug = false //调试开关
	g_error = NewError()
	g_error_cnt = 1
	g_statement_stack = NewStatementStack()
	g_symbols = NewScopedSymbolTable("global", 1, nil)
	g_builtin = NewBuiltinSymbolTable("builtin")
	g_builtin.set("print", NewFunc(true, nil, "print", &Param{flag: -1}, nil))
	g_builtin.set("list", NewFunc(true, nil, "list", &Param{flag: -1}, nil))
	g_builtin.set("exit", NewFunc(true, nil, "exit", &Param{flag: 0}, nil))

	inScope := NewScopedSymbolTable(g_symbols.scopeName+"_class", g_symbols.scopeLevel+1, g_symbols)
	obj := &Class{scope: inScope,
		name: &Variable{name: "Ojbect"},
		mems: []AstNode{}}
	cls_init_func := &Func{scope: inScope,
		name:      "Object",
		isBuiltin: false,
		params:    &Param{flag: 0},
		body: &LocalCompoundStatement{
			nodes: []AstNode{&ReturnStatement{results: []AstNode{obj}}}}}
	inScope.set("Object", cls_init_func)

	g_builtin.set("Object", obj)
}
