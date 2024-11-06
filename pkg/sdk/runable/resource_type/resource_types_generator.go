package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/opengovern/og-describer-template/provider/configs"
	"os"
	"sort"
	"strings"
	"text/template"
)

//go:embed resource-types.json
var ResourceTypes string

type ResourceType struct {
	ResourceName   string
	Tags           map[string][]string
	TagsString     string `json:"-"`
	ListDescriber  string
	GetDescriber   string
	SteampipeTable string
	Model          string
}

var (
	output   = flag.String("output", "", "")
	indexMap = flag.String("index-map", "", "")
)

func main() {
	flag.Parse()

	var resourceTypes []ResourceType

	if err := json.Unmarshal([]byte(ResourceTypes), &resourceTypes); err != nil {
		panic(err)
	}

	tmpl, err := template.New("").Parse(fmt.Sprintf(`
	"{{ .ResourceName }}": {
		IntegrationType:      configs.IntegrationName,
		ResourceName:         "{{ .ResourceName }}",
		Tags:                 {{ .TagsString }},
		ListDescriber:        {{ .ListDescriber }},
		GetDescriber:         {{ if .GetDescriber }}{{ .GetDescriber }}{{ else }}nil{{ end }},
	},
`))
	if err != nil {
		panic(err)
	}

	if output == nil || len(*output) == 0 {
		v := "resource_types.go"
		output = &v
	}

	if indexMap == nil || len(*indexMap) == 0 {
		v := "table_index_map.go"
		indexMap = &v
	}

	b := &strings.Builder{}
	b.WriteString(fmt.Sprintf(`
package provider
import (
	"%[1]s/provider/describer"
	"%[1]s/provider/configs"
	model "github.com/opengovern/og-describer-azure/pkg/sdk/models"
)
var ResourceTypes = map[string]model.ResourceType{
`, configs.OGPluginRepoURL))
	for _, resourceType := range resourceTypes {
		var arr []string

		tagsStringBuilder := strings.Builder{}
		tagsStringBuilder.WriteString("map[string][]string{\n")

		var tagsLines []string
		for k, v := range resourceType.Tags {
			arr = []string{}
			for _, t := range v {
				arr = append(arr, "\""+t+"\"")
			}
			tagsLines = append(tagsLines, fmt.Sprintf("            \"%s\": {%s},\n", k, strings.Join(arr, ",")))
		}

		sort.Strings(tagsLines)
		for _, l := range tagsLines {
			tagsStringBuilder.WriteString(l)
		}

		tagsStringBuilder.WriteString("        }")
		resourceType.TagsString = tagsStringBuilder.String()
		err = tmpl.Execute(b, resourceType)
		if err != nil {
			panic(err)
		}
	}
	b.WriteString("}\n")

	err = os.WriteFile(*output, []byte(b.String()), os.ModePerm)
	if err != nil {
		panic(err)
	}

	b = &strings.Builder{}
	b.WriteString(fmt.Sprintf(`
package steampipe

import (
	"%[1]s/pkg/sdk/es"
)

var Map = map[string]string{
`, configs.OGPluginRepoURL))
	for _, resourceType := range resourceTypes {
		b.WriteString(fmt.Sprintf("  \"%s\": \"%s\",\n", resourceType.ResourceName, resourceType.SteampipeTable))
	}
	b.WriteString(fmt.Sprintf(`}

var DescriptionMap = map[string]interface{}{
`))
	for _, resourceType := range resourceTypes {
		b.WriteString(fmt.Sprintf("  \"%s\": opengovernance.%s{},\n", resourceType.ResourceName, resourceType.Model))
	}
	b.WriteString(fmt.Sprintf(`}

var ReverseMap = map[string]string{
`))

	// reverse map
	for _, resourceType := range resourceTypes {
		b.WriteString(fmt.Sprintf("  \"%s\": \"%s\",\n", resourceType.SteampipeTable, resourceType.ResourceName))
	}
	b.WriteString("}\n")

	err = os.WriteFile(*indexMap, []byte(b.String()), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
