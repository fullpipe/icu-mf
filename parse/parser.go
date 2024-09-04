package parse

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type PlainArg struct {
	Name string `"{" @Ident "}"`
}

type Func struct {
	ArgName string `"{" @Ident `
	Func    string `"," @Ident`
	Param   string `("," @Ident )? "}"`
}

type Expr struct {
	Name   string  `"{" @Ident`
	Func   string  `("," @Ident)?`
	Offset int     `("," "offset" ":" @Int)?`
	Cases  []*Case `(@@*)? "}"`
}

type Case struct {
	Name    string   `(@Ident | @Case)`
	Message *Message `"{" @@ "}"`
}

type Message struct {
	Fragments []*Fragment `@@*`
}

type Fragment struct {
	Text       string    `(@String | @SubMessageString)`
	PlainArg   *PlainArg `| @@`
	Func       *Func     `| @@`
	Expr       *Expr     `| @@`
	Octothorpe bool      `| @"#"`
}

func NewParser() *participle.Parser[Message] {
	def := lexer.MustStateful(lexer.Rules{
		"Root": {
			{Name: `String`, Pattern: `[^{]+`, Action: nil},
			{Name: `Expr`, Pattern: `{`, Action: lexer.Push("Expr")},
		},
		"Expr": {
			{Name: `Whitespace`, Pattern: `\s+`, Action: nil},
			{Name: `Punctuation`, Pattern: `[,:]`, Action: nil},
			{Name: `Int`, Pattern: `\d+`, Action: nil},
			{Name: `Ident`, Pattern: `\w+`, Action: nil},
			{Name: `Case`, Pattern: `=\d+`, Action: nil},
			{Name: `ExprEnd`, Pattern: `}`, Action: lexer.Pop()},
			{Name: `SubMessage`, Pattern: `{`, Action: lexer.Push("SubMessage")},
		},
		"SubMessage": {
			{Name: `SubMessageString`, Pattern: `[^{^}^#]+`, Action: nil},
			{Name: `Octothorpe`, Pattern: `#`, Action: nil},
			{Name: `Expr`, Pattern: `{`, Action: lexer.Push("Expr")},
			{Name: `SubMessageEnd`, Pattern: `}`, Action: lexer.Pop()},
		},
	})

	parser := participle.MustBuild[Message](
		participle.Lexer(def),
		participle.Elide("Whitespace", "Punctuation"),
		participle.UseLookahead(4),
	)

	return parser
}
