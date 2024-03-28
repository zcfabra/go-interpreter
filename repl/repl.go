package repl

import (
	"bufio"
	"fmt"
	"io"
	"lang/evaluator"
	"lang/lexer"
	"lang/object"
	"lang/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			handleParserErrors(out, p.Errors())
			continue
		}

		if evaluated := evaluator.Eval(program, env); evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}

}

func handleParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
