package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"regexp"
	"strings"
	"text/template"

	"golang.org/x/tools/imports"
)

var re = regexp.MustCompile(`json:"(\d|\w+)(,omitempty)?"`)

const Object = "Object"

type Parameter struct {
	Field       string
	Type        string
	Description string
	Tag         string
}

type ResponseField struct {
	Field       string
	Type        string
	Description string
	Tag         string
}

type Response struct {
	Type   string
	Fields []ResponseField
}

type Endpoint struct {
	Name        string
	Description string

	Parameters      []Parameter
	ParameterByName map[string]Parameter
	Response        *Response
}

func NewEndpoint(name string) *Endpoint {
	return &Endpoint{
		Name:            name,
		Parameters:      make([]Parameter, 0),
		ParameterByName: make(map[string]Parameter),
	}
}

type Parser struct {
	Endpoints       []*Endpoint
	EndpointsByName map[string]*Endpoint

	file *ast.File

	serverTemplate          *template.Template
	clientTemplate          *template.Template
	loggingTemplate         *template.Template
	instrumentationTemplate *template.Template
	recorderTemplate        *template.Template
}

func NewParser(filename string) (*Parser, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	clientTemplate, err := ioutil.ReadFile("../templates/client.go.tmpl")
	if err != nil {
		return nil, err
	}
	serverTemplate, err := ioutil.ReadFile("../templates/server.go.tmpl")
	if err != nil {
		return nil, err
	}
	loggingTemplate, err := ioutil.ReadFile("../templates/logging.go.tmpl")
	if err != nil {
		return nil, err
	}
	instrumentationTemplate, err := ioutil.ReadFile("../templates/instrumentation.go.tmpl")
	if err != nil {
		return nil, err
	}
	recorderTemplate, err := ioutil.ReadFile("../templates/recorder.go.tmpl")
	if err != nil {
		return nil, err
	}

	return &Parser{
		file:                    f,
		Endpoints:               make([]*Endpoint, 0),
		EndpointsByName:         make(map[string]*Endpoint),
		serverTemplate:          template.Must(template.New("server").Parse(string(serverTemplate))),
		clientTemplate:          template.Must(template.New("client").Parse(string(clientTemplate))),
		loggingTemplate:         template.Must(template.New("logging").Parse(string(loggingTemplate))),
		instrumentationTemplate: template.Must(template.New("instrumentation").Parse(string(instrumentationTemplate))),
		recorderTemplate:        template.Must(template.New("instrumentation").Parse(string(recorderTemplate))),
	}, nil
}

func (p *Parser) Parse() error {
	ast.Inspect(p.file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if err := p.parseType(x); err != nil {
				return false
			}
		case *ast.FuncDecl:
			p.parseFunction(x)
		}
		return true
	})
	return p.generate()
}

func (p *Parser) generate() error {
	err := p.generateCode(p.clientTemplate, "../client/endpoints.go")
	if err != nil {
		return err
	}
	err = p.generateCode(p.serverTemplate, "../server/endpoints.go")
	if err != nil {
		return err
	}
	err = p.generateCode(p.loggingTemplate, "logging.go")
	if err != nil {
		return err
	}
	err = p.generateCode(p.instrumentationTemplate, "instrumentation.go")
	if err != nil {
		return err
	}
	return p.generateCode(p.recorderTemplate, "recorder.go")
}

func (p *Parser) generateCode(tmpl *template.Template, fn string) error {
	fmt.Printf("Generating %s\n", fn)
	buf := bytes.NewBuffer([]byte{})
	err := tmpl.Execute(buf, p)
	if err != nil {
		return err
	}
	res, err := imports.Process(fn, buf.Bytes(), nil)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fn, res, 0666)
}

func (p *Parser) parseType(st *ast.TypeSpec) error {
	if strings.HasSuffix(st.Name.Name, "Request") {
		endpoint := strings.Replace(st.Name.Name, "Request", "", -1)
		p.addParameter(endpoint, st.Type.(*ast.StructType))
	}
	if strings.HasSuffix(st.Name.Name, "Response") {
		endpoint := strings.Replace(st.Name.Name, "Response", "", -1)
		p.addResponseField(endpoint, st.Type)
	}
	return nil
}

func (p *Parser) addParameter(endpoint string, st *ast.StructType) {
	for _, field := range st.Fields.List {
		params := Parameter{
			Field:       field.Names[0].Name,
			Description: field.Doc.Text(),
			Tag:         parseTag(field.Tag.Value),
			Type:        mapFieldType(field.Type),
		}

		if p.EndpointsByName[endpoint] == nil {
			p.EndpointsByName[endpoint] = NewEndpoint(endpoint)
		}
		p.EndpointsByName[endpoint].Parameters = append(p.EndpointsByName[endpoint].Parameters, params)
		p.EndpointsByName[endpoint].ParameterByName[params.Field] = params
	}
}

func (p *Parser) addResponseField(endpoint string, expr ast.Expr) {
	if p.EndpointsByName[endpoint] == nil {
		p.EndpointsByName[endpoint] = NewEndpoint(endpoint)
	}
	response := &Response{}
	switch x := expr.(type) {
	case *ast.StructType:
		response.Type = Object
		response.Fields = make([]ResponseField, 0)
		for _, field := range x.Fields.List {
			responseField := ResponseField{
				Description: field.Doc.Text(),
				Field:       parseTag(field.Tag.Value),
				Type:        mapFieldType(field.Type),
			}
			response.Fields = append(response.Fields, responseField)
		}
	default:
		response.Type = mapFieldType(x)
	}
	if response.Type != Object || len(response.Fields) != 0 {
		p.EndpointsByName[endpoint].Response = response
	}
}

func mapFieldType(expr ast.Expr) string {
	switch x := expr.(type) {
	case *ast.Ident:
		return x.Name
	case *ast.StarExpr:
		ident, ok := x.X.(*ast.Ident)
		if ok {
			return ident.Name
		}
		return Object
	case *ast.SelectorExpr:
		name := fmt.Sprintf("%v.%s", x.X, x.Sel.Name)
		switch name {
		case "globalid.ID":
			return "UUID"
		case "model.ReactionType":
			return "string"
		case "model.CardsResponse", "model.CardResponse", "model.Draft":
			return Object
		}
		return name
	case *ast.ArrayType:
		return "Array"
	default:
		panic(fmt.Sprintf("Unmapped type %T %v", x, x))
	}
}

func parseTag(tag string) string {
	match := re.FindStringSubmatch(tag)
	return match[1]
}

func (p *Parser) parseFunction(fd *ast.FuncDecl) {
	if fd.Recv == nil {
		return
	}
	if recv, ok := fd.Recv.List[0].Type.(*ast.StarExpr); ok {
		if ident, ok := recv.X.(*ast.Ident); ok {
			name := fd.Name.Name
			description := fd.Doc.Text()
			firstChar := string(name[0])
			if ident.Name == "api" && firstChar == strings.ToUpper(firstChar) {
				p.AddEndpoint(name, description)
			}
		}
	}
}

func (p *Parser) AddEndpoint(name, description string) {
	endpoint := p.EndpointsByName[name]
	if endpoint == nil {
		endpoint = NewEndpoint(name)
	}
	endpoint.Description = enhanceDescription(description, name)

	p.EndpointsByName[name] = endpoint
	p.Endpoints = append(p.Endpoints, endpoint)
}

func enhanceDescription(description, rpc string) string {
	return toUpper(strings.Replace(description, fmt.Sprintf("%s ", rpc), "", -1))
}

func toUpper(s string) string {
	if s == "" {
		return ""
	}
	return fmt.Sprintf("%s%s", strings.ToUpper(string(s[0])), s[1:])
}
