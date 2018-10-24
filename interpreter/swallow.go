package interpreter

type Swallow struct {
	parser *Parser
}

func NewSwallow(text, file string) *Swallow {
	this := &Swallow{}
	this.parser = NewParser(NewLexer(text, file))

	return this
}

func (this *Swallow) interpreter() {
	this.parser.parser()
}
