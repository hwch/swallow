package core

func init() {
	g_is_debug = false //调试开关
	g_error = NewError()
	g_error_cnt = 1
	g_statement_stack = NewStatementStack()
	g_symbols = NewScopedSymbolTable("global", 1, nil)
	g_builtin = NewBuiltinSymbolTable("builtin")
	g_builtin.set("print", NewFunc(true, nil, "print", NewParam(&Token{}, -1, nil, nil), nil, nil))
	g_builtin.set("list", NewFunc(true, nil, "list", NewParam(&Token{}, -1, nil, nil), nil, nil))
	g_builtin.set("exit", NewFunc(true, nil, "exit", NewParam(&Token{}, 0, nil, nil), nil, nil))

	inScope := NewScopedSymbolTable(g_symbols.scopeName+"_class", g_symbols.scopeLevel+1, g_symbols)
	obj := NewClass(&Token{}, "Object", nil /*parent*/, nil /*member*/, inScope)
	inScope.set("Object", NewFunc(false, &Token{}, "Object",
		NewParam(&Token{}, 0, nil, inScope),
		NewLocalCompoundStatement(&Token{}, []AstNode{NewReturnStatement(&Token{}, []AstNode{obj})}, inScope),
		inScope))

	g_builtin.set("Object", obj)
}
