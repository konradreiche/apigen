package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/konradreiche/apigen/parser"
	"github.com/october93/engine/kit/globalid"
	"github.com/october93/engine/rpc/protocol"

	"golang.org/x/tools/imports"
)

type DocGenerator struct {
	Parser   *parser.Parser
	file     *ast.File
	Examples map[string][]*ExamplePair
}

type ExamplePair struct {
	Request   *Example
	Response  *Example
	RequestID globalid.ID
}

func NewExamplePair() *ExamplePair {
	return &ExamplePair{RequestID: globalid.Next()}
}

type Example struct {
	RPC         string
	Description string
	Payload     string
	Kind        string
	RequestID   globalid.ID
}

func NewExample(rpc, description, kind string) *Example {
	return &Example{
		RPC:         rpc,
		Description: description,
		Kind:        kind,
	}
}

func toLower(rpc string) string {
	return fmt.Sprintf("%s%s", strings.ToLower(string(rpc[0])), rpc[1:])
}

func (e *Example) JSON() (string, error) {
	m := protocol.NewMessage(e.RPC)
	m.RPC = toLower(e.RPC)

	if e.Kind == "Request" {
		m.RequestID = e.RequestID
	} else {
		m.RequestID = globalid.Nil
		m.Ack = e.RequestID
	}

	if e.Payload != "{}" {
		m.Data = json.RawMessage([]byte(e.Payload))
	}
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	var indented bytes.Buffer
	err = json.Indent(&indented, b, "", "  ")
	if err != nil {
		return "", err
	}
	return indented.String(), nil
}

func NewDocGenerator(parser *parser.Parser, fn string) (*DocGenerator, error) {
	fset := token.NewFileSet()
	file, err := goparser.ParseFile(fset, fn, nil, goparser.ParseComments)
	if err != nil {
		return nil, err
	}
	return &DocGenerator{
		Parser:   parser,
		file:     file,
		Examples: make(map[string][]*ExamplePair),
	}, nil
}

func (dg *DocGenerator) ParseExamples() error {
	examplePair := NewExamplePair()
	var example *Example

	ast.Inspect(dg.file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.GenDecl:
			value := x.Specs[0]
			vs, ok := value.(*ast.ValueSpec)
			if !ok {
				// TODO handle
				return true
			}

			name := vs.Names[0].String()
			re := regexp.MustCompile(`(.*)(Request|Response)Example\d`)
			match := re.FindStringSubmatch(name)
			if len(match) == 0 {
				return true
			}
			endpoint := match[1]

			example = NewExample(endpoint, x.Doc.Text(), match[2])

			example.RPC = endpoint
			example.Payload = examples[name]
			if match[2] == "Request" {
				examplePair.Request = example
				example.RequestID = examplePair.RequestID
				if examplePair.Request != nil && examplePair.Response != nil {
					dg.Examples[endpoint] = append(dg.Examples[endpoint], examplePair)
					examplePair = NewExamplePair()
				}
			} else if match[2] == "Response" {
				examplePair.Response = example
				example.RequestID = examplePair.RequestID
				if examplePair.Request != nil && examplePair.Response != nil {
					dg.Examples[endpoint] = append(dg.Examples[endpoint], examplePair)
					examplePair = NewExamplePair()
				}
			}
		default:
		}
		return true
	})
	return nil
}

func unfold(rpc string) string {
	var b bytes.Buffer
	for i, r := range rpc {
		if i == 0 {
			b.WriteRune(r)
		} else if unicode.IsUpper(r) {
			b.WriteString(" ")
			b.WriteRune(r)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func (dg *DocGenerator) generate() error {
	// TODO (konrad) refactor
	tmpl := template.New("documentation")
	funcs := make(template.FuncMap, 1)
	funcs["toLower"] = func(s string) string {
		return strings.ToLower(s)
	}
	funcs["toDash"] = func(s string) string {
		return strings.ToLower(s)
	}
	funcs["json"] = func(example *Example) (string, error) {
		return example.JSON()
	}

	funcs["unfold"] = unfold
	tmpl.Funcs(funcs)
	docTemplate, err := ioutil.ReadFile("../templates/doc.html.tmpl")
	if err != nil {
		return err
	}

	temp := template.Must(tmpl.Parse(string(docTemplate)))
	err = os.MkdirAll("../docs", 0700)
	if err != nil {
		return err
	}
	ls, err := ioutil.ReadDir("../cmd/subcmd/assets")
	if err != nil {
		return err
	}
	for _, file := range ls {
		// TODO (konrad) refactor to use something offered by the Go library?
		cmd := exec.Command("cp", fmt.Sprintf("../cmd/subcmd/assets/%s", file.Name()), fmt.Sprintf("../docs/%s", file.Name())) // #nosec
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	buf := bytes.NewBuffer([]byte{})
	err = temp.Execute(buf, dg)
	if err != nil {
		return err
	}
	absolutePath, err := filepath.Abs("../docs/index.html")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(absolutePath, buf.Bytes(), 0666)
}

func bootstrapExamples(filename string) error {
	fset := token.NewFileSet()
	f, err := goparser.ParseFile(fset, filename, nil, goparser.ParseComments)
	if err != nil {
		return err
	}
	t, err := ioutil.ReadFile("../templates/examples.go.tmpl")
	if err != nil {
		return err
	}
	examples := make([]string, 0)
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ValueSpec:
			examples = append(examples, x.Names[0].Name)
		}
		return true
	})
	tmpl := template.Must(template.New("examples").Parse(string(t)))
	buf := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(buf, examples)
	if err != nil {
		return err
	}
	res, err := imports.Process("../cmd/subcmd/examples.go", buf.Bytes(), nil)
	if err != nil {
		return err
	}
	return ioutil.WriteFile("../cmd/subcmd/examples.go", res, 0666)
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "examples" {
		err := bootstrapExamples(os.Getenv("GOFILE"))
		if err != nil {
			fail(err)
		}
		return
	}

	err := loadExamples()
	if err != nil {
		fail(err)
	}
	p, err := parser.NewParser("api.go")
	if err != nil {
		fail(err)
	}
	err = p.Parse()
	if err != nil {
		fail(err)
	}
	dg, err := NewDocGenerator(p, "../api/examples.go")
	if err != nil {
		fail(err)
	}
	err = dg.ParseExamples()
	if err != nil {
		fail(err)
	}
	err = dg.generate()
	if err != nil {
		fail(err)
	}

}

func fail(err error) {
	fmt.Println(err)
	os.Exit(1)
}
