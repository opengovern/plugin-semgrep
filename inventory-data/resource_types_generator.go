package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"
)

//go:embed aws-resource-types.json
var awsResourceTypes string

type DiscoveryStatus string

const (
	DiscoveryStatus_COMPLETE = "COMPLETE"
	DiscoveryStatus_FAST     = "FAST"
	DiscoveryStatus_COST     = "COST"
	DiscoveryStatus_DISABLED = "DISABLED"
)

type ResourceType struct {
	ResourceName         string
	ResourceLabel        string
	Category             []string
	Tags                 map[string][]string
	TagsString           string `json:"-"`
	ServiceName          string
	ListDescriber        string
	GetDescriber         string
	TerraformName        []string
	TerraformNameString  string `json:"-"`
	TerraformServiceName string
	Discovery            DiscoveryStatus
	IgnoreSummarize      bool
	SteampipeTable       string
	Model                string
}

var (
	provider = flag.String("provider", "", "")
	output   = flag.String("output", "", "")
	indexMap = flag.String("index-map", "", "")
)

func main() {
	flag.Parse()

	if provider == nil || *provider == "" {
		v := "aws"
		provider = &v
	}

	var resourceTypes []ResourceType
	var cloud string
	var upperProvider string

	if err := json.Unmarshal([]byte(awsResourceTypes), &resourceTypes); err != nil {
		panic(err)
	}
	cloud = "CloudAWS"
	upperProvider = "AWS"

	tmpl, err := template.New("").Parse(fmt.Sprintf(`
	"{{ .ResourceName }}": {
		Connector:            source.%s,
		ResourceName:         "{{ .ResourceName }}",
		ResourceLabel:        "{{ .ResourceLabel }}",
		Tags:                 {{ .TagsString }},
		ServiceName:          "{{ .ServiceName }}",
		ListDescriber:        {{ .ListDescriber }},
		GetDescriber:         {{ if .GetDescriber }}{{ .GetDescriber }}{{ else }}nil{{ end }},
		TerraformName:        {{ .TerraformNameString }},
		TerraformServiceName: "{{ .TerraformServiceName }}",
		FastDiscovery:        {{ if eq .Discovery "FAST" }}true{{ else }}false{{ end }},{{ if eq .Discovery "COST" }}
		CostDiscovery:		  true,{{ end }}
		Summarize:            {{ if .IgnoreSummarize }}false{{ else }}true{{ end }},
	},
`, cloud))
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
package %[1]s
import (
	"github.com/opengovern/og-%[1]s-describer/%[1]s/describer"
	"github.com/opengovern/og-util/pkg/source"
)
var resourceTypes = map[string]ResourceType{
`, *provider))
	for _, resourceType := range resourceTypes {
		if resourceType.Discovery == DiscoveryStatus_DISABLED {
			continue
		}
		var arr []string
		for _, t := range resourceType.TerraformName {
			arr = append(arr, "\""+t+"\"")
		}
		resourceType.TerraformNameString = "[]string{" + strings.Join(arr, ",") + "}"

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
	"github.com/opengovern/og-%[1]s-describer/pkg/opengovernance-es-sdk"
)

var %[1]sMap = map[string]string{
`, *provider))
	for _, resourceType := range resourceTypes {
		b.WriteString(fmt.Sprintf("  \"%s\": \"%s\",\n", resourceType.ResourceName, resourceType.SteampipeTable))
	}
	b.WriteString(fmt.Sprintf(`}

var %sDescriptionMap = map[string]interface{}{
`, upperProvider))
	for _, resourceType := range resourceTypes {
		b.WriteString(fmt.Sprintf("  \"%s\": opengovernance.%s{},\n", resourceType.ResourceName, resourceType.Model))
	}
	b.WriteString(fmt.Sprintf(`}

var %[1]sReverseMap = map[string]string{
`, upperProvider))

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
