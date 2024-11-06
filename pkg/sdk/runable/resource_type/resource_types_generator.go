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

// Embed the JSON file containing resource types
//
//go:embed resource-types.json
var ResourceTypes string

// Define the ResourceType struct with Labels and Annotations
type ResourceType struct {
	ResourceName      string
	Tags              map[string][]string
	TagsString        string `json:"-"`
	ListDescriber     string
	GetDescriber      string
	SteampipeTable    string
	Model             string
	Annotations       map[string]string
	Labels            map[string]string
	AnnotationsString string `json:"-"`
	LabelsString      string `json:"-"`
}

var (
	output   = flag.String("output", "", "Path to the output file for resource types")
	indexMap = flag.String("index-map", "", "Path to the output file for index map")
)

func main() {
	flag.Parse()

	var resourceTypes []ResourceType

	// Parse the embedded JSON into resourceTypes slice
	if err := json.Unmarshal([]byte(ResourceTypes), &resourceTypes); err != nil {
		panic(err)
	}

	// Define the template with Labels and Annotations included
	tmpl, err := template.New("").Parse(fmt.Sprintf(`
	"{{ .ResourceName }}": {
		IntegrationType:      configs.IntegrationName,
		ResourceName:         "{{ .ResourceName }}",
		Tags:                 {{ .TagsString }},
		Labels:               {{ .LabelsString }},
		Annotations:          {{ .AnnotationsString }},
		ListDescriber:        {{ .ListDescriber }},
		GetDescriber:         {{ if .GetDescriber }}{{ .GetDescriber }}{{ else }}nil{{ end }},
	},
`))
	if err != nil {
		panic(err)
	}

	// Set default output paths if not provided
	if output == nil || len(*output) == 0 {
		v := "resource_types.go"
		output = &v
	}

	if indexMap == nil || len(*indexMap) == 0 {
		v := "table_index_map.go"
		indexMap = &v
	}

	// Initialize a strings.Builder to construct the output file content
	b := &strings.Builder{}
	b.WriteString(fmt.Sprintf(`package provider
import (
	"%[1]s/provider/describer"
	"%[1]s/provider/configs"
	model "github.com/opengovern/og-describer-azure/pkg/sdk/models"
)
var ResourceTypes = map[string]model.ResourceType{
`, configs.OGPluginRepoURL))

	// Iterate over each resource type to build its string representations
	for _, resourceType := range resourceTypes {
		var arr []string

		// Build TagsString
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

		sort.Strings(tagsLines) // Sort for consistency
		for _, l := range tagsLines {
			tagsStringBuilder.WriteString(l)
		}

		tagsStringBuilder.WriteString("        }")
		resourceType.TagsString = tagsStringBuilder.String()

		// Build LabelsString
		labelsStringBuilder := strings.Builder{}
		labelsStringBuilder.WriteString("map[string]string{\n")

		var labelsLines []string
		for k, v := range resourceType.Labels {
			// Escape quotes in keys and values
			escapedKey := escapeString(k)
			escapedValue := escapeString(v)
			labelsLines = append(labelsLines, fmt.Sprintf("            \"%s\": \"%s\",\n", escapedKey, escapedValue))
		}

		sort.Strings(labelsLines) // Sort for consistency
		for _, l := range labelsLines {
			labelsStringBuilder.WriteString(l)
		}

		labelsStringBuilder.WriteString("        }")
		resourceType.LabelsString = labelsStringBuilder.String()

		// Build AnnotationsString
		annotationsStringBuilder := strings.Builder{}
		annotationsStringBuilder.WriteString("map[string]string{\n")

		var annotationsLines []string
		for k, v := range resourceType.Annotations {
			// Escape quotes in keys and values
			escapedKey := escapeString(k)
			escapedValue := escapeString(v)
			annotationsLines = append(annotationsLines, fmt.Sprintf("            \"%s\": \"%s\",\n", escapedKey, escapedValue))
		}

		sort.Strings(annotationsLines) // Sort for consistency
		for _, l := range annotationsLines {
			annotationsStringBuilder.WriteString(l)
		}

		annotationsStringBuilder.WriteString("        }")
		resourceType.AnnotationsString = annotationsStringBuilder.String()

		// Execute the template with the current resourceType
		err = tmpl.Execute(b, resourceType)
		if err != nil {
			panic(err)
		}
	}
	b.WriteString("}\n")

	// Write the generated content to the output file
	err = os.WriteFile(*output, []byte(b.String()), os.ModePerm)
	if err != nil {
		panic(err)
	}

	// Generate the index map file as before
	b = &strings.Builder{}
	b.WriteString(fmt.Sprintf(`package steampipe

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

	// Build the reverse map
	for _, resourceType := range resourceTypes {
		b.WriteString(fmt.Sprintf("  \"%s\": \"%s\",\n", resourceType.SteampipeTable, resourceType.ResourceName))
	}
	b.WriteString("}\n")

	// Write the index map to the specified file
	err = os.WriteFile(*indexMap, []byte(b.String()), os.ModePerm)
	if err != nil {
		panic(err)
	}
}

// escapeString ensures that any quotes in the strings are properly escaped
func escapeString(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}
