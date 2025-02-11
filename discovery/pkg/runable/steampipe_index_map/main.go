package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/opengovern/og-describer-semgrep/global"
	"os"
	"strings"
)

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
	resourceTypesFile = flag.String("resourceTypesFile", "", "Location of the resource types json file file")
	indexMap          = flag.String("index-map", "", "Path to the output file for index map")
)

func main() {
	flag.Parse()

	var resourceTypes []ResourceType

	if resourceTypesFile == nil || len(*resourceTypesFile) == 0 {
		rt := "global/maps/resource-types.json"
		resourceTypesFile = &rt
	}

	if indexMap == nil || len(*indexMap) == 0 {
		v := "global/maps/table_index_map.gen.go"
		indexMap = &v
	}

	resourceTypesFileJson, err := os.ReadFile(*resourceTypesFile)
	if err != nil {
		panic(err)
	}
	// Parse the embedded JSON into resourceTypes slice
	if err := json.Unmarshal(resourceTypesFileJson, &resourceTypes); err != nil {
		panic(err)
	}

	// Generate the index map file as before
	b := &strings.Builder{}
	b.WriteString(fmt.Sprintf(`package maps

import (
	"%[1]s/discovery/pkg/es"
)

var ResourceTypesToTables = map[string]string{
`, global.OGPluginRepoURL))
	for _, resourceType := range resourceTypes {
		b.WriteString(fmt.Sprintf("  \"%s\": \"%s\",\n", resourceType.ResourceName, resourceType.SteampipeTable))
	}
	b.WriteString(fmt.Sprintf(`}

var ResourceTypeToDescription = map[string]interface{}{
`))
	for _, resourceType := range resourceTypes {
		b.WriteString(fmt.Sprintf("  \"%s\": opengovernance.%s{},\n", resourceType.ResourceName, resourceType.Model))
	}
	b.WriteString(fmt.Sprintf(`}

var TablesToResourceTypes = map[string]string{
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
