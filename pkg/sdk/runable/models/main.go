package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"html/template"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	file       = flag.String("file", "", "Location of the model file")
	sourceType = flag.String("type", "", "Type of resource clients (e.g. , azure). Should match the model import path")
	output     = flag.String("output", "", "Location of the output file")
)

type SourceType struct {
	Name        string
	Index       string
	ListFilters map[string]string
	GetFilters  map[string]string
	SourceType  string
}

type ResourceType struct {
	ResourceName         string
	ResourceLabel        string
	ServiceName          string
	ListDescriber        string
	GetDescriber         string
	TerraformName        []string
	TerraformNameString  string `json:"-"`
	TerraformServiceName string
	FastDiscovery        bool
	SteampipeTable       string
	Model                string
}

func main() {
	rt := "../resourceType/resource-types.json"
	b, err := os.ReadFile(rt)
	if err != nil {
		panic(err)
	}
	var resourceTypes []ResourceType
	err = json.Unmarshal(b, &resourceTypes)
	if err != nil {
		panic(err)
	}

	flag.CommandLine.Init("gen", flag.ExitOnError)
	flag.Parse()

	tpl := template.New("types")
	_, err = tpl.Parse(`
// ==========================  START: {{ .Name }} =============================

type {{ .Name }} struct {
	Description   {{ .SourceType }}.{{ .Name }}Description 	` + "`json:\"description\"`" + `
	Metadata      {{ .SourceType }}.Metadata 					` + "`json:\"metadata\"`" + `
	ResourceJobID int ` + "`json:\"resource_job_id\"`" + `
	SourceJobID   int ` + "`json:\"source_job_id\"`" + `
	ResourceType  string ` + "`json:\"resource_type\"`" + `
	SourceType    string ` + "`json:\"source_type\"`" + `
	ID            string ` + "`json:\"id\"`" + `
	ARN            string ` + "`json:\"arn\"`" + `
	SourceID      string ` + "`json:\"source_id\"`" + `
}

type {{ .Name }}Hit struct {
	ID      string            ` + "`json:\"_id\"`" + `
	Score   float64           ` + "`json:\"_score\"`" + `
	Index   string            ` + "`json:\"_index\"`" + `
	Type    string            ` + "`json:\"_type\"`" + `
	Version int64             ` + "`json:\"_version,omitempty\"`" + `
	Source  {{ .Name }}       ` + "`json:\"_source\"`" + `
	Sort    []interface{}     ` + "`json:\"sort\"`" + `
}

type {{ .Name }}Hits struct {
	Total essdk.SearchTotal       ` + "`json:\"total\"`" + `
	Hits  []{{ .Name }}Hit ` + "`json:\"hits\"`" + `
}

type {{ .Name }}SearchResponse struct {
	PitID string          ` + "`json:\"pit_id\"`" + `
	Hits  {{ .Name }}Hits ` + "`json:\"hits\"`" + `
}

type {{ .Name }}Paginator struct {
	paginator *essdk.BaseESPaginator
}

func (k Client) New{{ .Name }}Paginator(filters []essdk.BoolFilter, limit *int64) ({{ .Name }}Paginator, error) {
	paginator, err := essdk.NewPaginator(k.ES(), "{{ .Index }}", filters, limit)
	if err != nil {
		return {{ .Name }}Paginator{}, err
	}

	p := {{ .Name }}Paginator{
		paginator: paginator,
	}

	return p, nil
}

func (p {{ .Name }}Paginator) HasNext() bool {
	return !p.paginator.Done()
}

func (p {{ .Name }}Paginator) Close(ctx context.Context) error {
	return p.paginator.Deallocate(ctx)
}

func (p {{ .Name }}Paginator) NextPage(ctx context.Context) ([]{{ .Name }}, error) {
	var response {{ .Name }}SearchResponse
	err := p.paginator.Search(ctx, &response)
	if err != nil {
		return nil, err
	}

	var values []{{ .Name }}
	for _, hit := range response.Hits.Hits {
		values = append(values, hit.Source)
	}

	hits := int64(len(response.Hits.Hits))
	if hits > 0 {
		p.paginator.UpdateState(hits, response.Hits.Hits[hits-1].Sort, response.PitID)
	} else {
		p.paginator.UpdateState(hits, nil, "")
	}

	return values, nil
}

var list{{ .Name }}Filters = map[string]string{
	{{ range $key, $value := .ListFilters }}"{{ $key }}": "{{ $value }}",
	{{ end }}
}

func List{{ .Name }}(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("List{{ .Name }}")
	runtime.GC()
	// create service
	cfg := essdk.GetConfig(d.Connection)
	ke, err := essdk.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		plugin.Logger(ctx).Error("List{{ .Name }} NewClientCached", "error", err)
		return nil, err
	}
	k := Client{Client: ke}

	sc, err := steampipesdk.NewSelfClientCached(ctx, d.ConnectionCache)
	if err != nil {
		plugin.Logger(ctx).Error("List{{ .Name }} NewSelfClientCached", "error", err)
		return nil, err
	}
	accountId, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyAccountID)
	if err != nil {
		plugin.Logger(ctx).Error("List{{ .Name }} GetConfigTableValueOrNil for OpenGovernanceConfigKeyAccountID", "error", err)
		return nil, err
	}
	encodedResourceCollectionFilters, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyResourceCollectionFilters)
	if err != nil {
		plugin.Logger(ctx).Error("List{{ .Name }} GetConfigTableValueOrNil for OpenGovernanceConfigKeyResourceCollectionFilters", "error", err)
		return nil, err
	}
	clientType, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyClientType)
	if err != nil {
		plugin.Logger(ctx).Error("List{{ .Name }} GetConfigTableValueOrNil for OpenGovernanceConfigKeyClientType", "error", err)
		return nil, err
	}

	paginator, err := k.New{{ .Name }}Paginator(essdk.BuildFilter(ctx, d.QueryContext, list{{ .Name }}Filters, "{{ .SourceType }}", accountId, encodedResourceCollectionFilters, clientType), d.QueryContext.Limit)
	if err != nil {
		plugin.Logger(ctx).Error("List{{ .Name }} New{{ .Name }}Paginator", "error", err)
		return nil, err
	}

	for paginator.HasNext() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("List{{ .Name }} paginator.NextPage", "error", err)
			return nil, err
		}

		for _, v := range page {
			d.StreamListItem(ctx, v)
		}
	}

	err = paginator.Close(ctx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}


var get{{ .Name }}Filters = map[string]string{
	{{ range $key, $value := .GetFilters }}"{{ $key }}": "{{ $value }}",
	{{ end }}
}

func Get{{ .Name }}(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("Get{{ .Name }}")
	runtime.GC()
	// create service
	cfg := essdk.GetConfig(d.Connection)
	ke, err := essdk.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		return nil, err
	}
	k := Client{Client: ke}

	sc, err := steampipesdk.NewSelfClientCached(ctx, d.ConnectionCache)
	if err != nil {
		return nil, err
	}
	accountId, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyAccountID)
	if err != nil {
		return nil, err
	}
	encodedResourceCollectionFilters, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyResourceCollectionFilters)
	if err != nil {
		return nil, err
	}
	clientType, err := sc.GetConfigTableValueOrNil(ctx, steampipesdk.OpenGovernanceConfigKeyClientType)
	if err != nil {
		return nil, err
	}

	limit := int64(1)
	paginator, err := k.New{{ .Name }}Paginator(essdk.BuildFilter(ctx, d.QueryContext, get{{ .Name }}Filters, "{{ .SourceType }}", accountId, encodedResourceCollectionFilters, clientType), &limit)
	if err != nil {
		return nil, err
	}

	for paginator.HasNext() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page {
			return v, nil
		}
	}

	err = paginator.Close(ctx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ==========================  END: {{ .Name }} =============================

`)
	if err != nil {
		log.Fatal(err)
	}

	out, err := os.Create(*output)
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, *file, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(&buf, "// Code is generated by go generate. DO NOT EDIT.")
	fmt.Fprintf(&buf, "package opengovernance")

	

	var sources []SourceType

	ast.Inspect(node, func(n ast.Node) bool {
		decl, ok := n.(*ast.GenDecl)
		if !ok || decl.Tok != token.TYPE {
			return true
		}

		for _, spec := range decl.Specs {
			t := spec.(*ast.TypeSpec)

			if !strings.HasSuffix(t.Name.String(), "Description") {
				fmt.Println("Ignoring type " + t.Name.String())
				continue
			}

			s := SourceType{
				Name:        strings.TrimSuffix(t.Name.String(), "Description"),
				SourceType:  *sourceType,
				GetFilters:  map[string]string{},
				ListFilters: map[string]string{},
			}
			exists := false
			for _, resourceType := range resourceTypes {
				if resourceType.Model == s.Name {
					exists = true
					var stopWordsRe = regexp.MustCompile(`\W+`)
					index := stopWordsRe.ReplaceAllString(resourceType.ResourceName, "_")
					index = strings.ToLower(index)
					s.Index = index

					tableAliasMap := map[string]string{
						"aws_api_gateway_authorizer": "table_aws_api_gateway_api_authorizer.go",
					}

					tableFile := fmt.Sprintf("table_%s.go", resourceType.SteampipeTable)
					if v, ok := tableAliasMap[resourceType.SteampipeTable]; ok {
						tableFile = v
					}
					plugin := "steampipe-plugin-aws/aws"
					fileName := "../../" + plugin + "/" + tableFile
					tableFileSet := token.NewFileSet()
					tableNode, err := parser.ParseFile(tableFileSet, fileName, nil, parser.ParseComments)
					if err != nil {
						panic(err)
					}

					ast.Inspect(tableNode, func(tnode ast.Node) bool {
						if c, ok := tnode.(*ast.CompositeLit); ok {

							var columnName, transformer string
							for _, arg := range c.Elts {
								if kv, ok := arg.(*ast.KeyValueExpr); ok {
									if i, ok := kv.Key.(*ast.Ident); ok {
										if i.Name == "Name" {
											if bl, ok := kv.Value.(*ast.BasicLit); ok {
												columnName = strings.Trim(bl.Value, "\"")
											}
										} else if i.Name == "Transform" {
											if cl, ok := kv.Value.(*ast.CallExpr); ok {
												transformer = extractTransformer(cl)
											}
										}
									}
								}
							}

							if columnName != "" && transformer != "" {
								transformer = strings.ToLower(transformer[:1]) + transformer[1:]
								if !strings.Contains(transformer, ".") {
									transformer = strings.ToLower(transformer)
								}
								s.GetFilters[columnName] = transformer
								s.ListFilters[columnName] = transformer
							}
							return true
						}
						return true
					})
				}
			}
			if !exists {
				fmt.Println("resourceType not found", s.Name)
			}

			if decl.Doc != nil {
				for _, c := range decl.Doc.List {
					if strings.HasPrefix(c.Text, "//index:") {
						//s.Index = strings.TrimSpace(strings.TrimPrefix(c.Text, "//index:"))
					} else if strings.HasPrefix(c.Text, "//getfilter:") {
						f := strings.TrimSpace(strings.TrimPrefix(c.Text, "//getfilter:"))
						fparts := strings.Split(f, "=")
						s.GetFilters[fparts[0]] = fparts[1]
					} else if strings.HasPrefix(c.Text, "//listfilter:") {
						f := strings.TrimSpace(strings.TrimPrefix(c.Text, "//listfilter:"))
						fparts := strings.Split(f, "=")
						s.ListFilters[fparts[0]] = fparts[1]
					}
				}
				s.GetFilters["og_account_id"] = "metadata.SourceID"
				s.ListFilters["og_account_id"] = "metadata.SourceID"
			}

			if s.Index != "" {
				sources = append(sources, s)
			} else {
				fmt.Println("ignoring due to empty index", s)
			}
		}

		return false
	})

	if len(sources) > 0 {
		fmt.Fprintln(&buf, `
		import (
			"context"
			"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
			essdk "github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
			steampipesdk "github.com/opengovern/og-util/pkg/steampipe"
			`+*sourceType+` "github.com/opengovern/og-`+*sourceType+`-describer/`+*sourceType+`/model"
            "runtime"
		)

		type Client struct {
			essdk.Client
		}
		`)
	}

	for _, source := range sources {
		err := tpl.Execute(&buf, source)
		if err != nil {
			log.Fatal(err)
		}
	}

	source, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	_, err = out.Write(source)
	if err != nil {
		log.Fatal(err)
	}
}

func extractTransformer(cl *ast.CallExpr) string {
	if sl, ok := cl.Fun.(*ast.SelectorExpr); ok {
		if sl.Sel.Name == "Transform" {
			return ""
		}
		if call, ok := sl.X.(*ast.CallExpr); ok {
			return extractTransformer(call)
		}
		if sl.Sel.Name == "FromField" {
			for _, arg := range cl.Args {
				if bl, ok := arg.(*ast.BasicLit); ok {
					return strings.Trim(bl.Value, "\"")
				}
			}
		}
	}
	return ""
}
