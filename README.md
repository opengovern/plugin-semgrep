# OpenGovernance Describer Template

This document is a GitHub repository temmplate for write describers for any Provider you want.

# Instructions

## 1. Create a new repository using this template

First, you need to fork this repository to your account. Then, you can create a new repository using this template.

## 2. Fill the Provider information

Fill the information of the Provider you want to describe in the [configs folder](./provider/configs/).

### 2.1. Provider information

Fill the Credential information of the Provider in the [credentials.go](./provider/configs/credentials.go) file.

```go
package configs

type IntegrationCredentials struct {
	// You should provide Credentials for any Provider.
}
```

### 2.2. Integration information

Fill the Integration information of the Provider in the [general.go](./provider/configs/general.go) file.

```go
package configs

import "github.com/opengovern/og-util/pkg/integration"

const (
	IntegrationTypeLower = "integrationType"                    // example: aws, azure
	IntegrationName      = integration.Type("INTEGRATION_NAME") // example: AWS_ACCOUNT, AZURE_SUBSCRIPTION
	OGPluginRepoURL      = "repo-url"                           // example: github.com/opengovern/og-describer-aws
)

```

## 3. Fill the describer wrapper

Fill the describer wrapper in the [describer_wrapper.go](./provider/descrdescriber_wrapperiber.go) file.

You should implement two functions:

DescribeListByProvider: This function should return a list of resources of the Provider.

DescribeSingleByProvider: This function should return a single resource of the Provider.

Theses functions are wrapper for the describer any resource of the Provider.

## 4. Create the describer file and implement the describer

### 4.1 Create the describer file

Create a new file in the [describers folder](./provider/describers/) with the name of the resource you want to describe.

### 4.2 Implement the describer

Implement the describer in the file you created in the previous step.
Implement two List and Get functions:

List: This function should return a list of resources of the Provider.

Get: This function should return a single resource of the Provider.

**Note:** You can use the [example describer](./provider/describers/example.go) as a reference. Example is for describing CohereAI datasets resource.

### 4.3 Fill model

You should fill the model of the resource in the [models.go](./provider/models/models.go).

You can define models for the resource you want to describe. The main model of the resource should have `Description` suffix.

```go
type DatasetDescription struct {
	ID                 string        
	Name               string        
	CreatedAt          time.Time     
	UpdatedAt          time.Time     
	DatasetType        string        
	ValidationStatus   string        
	ValidationError    string        
	Schema             string        
	RequiredFields     []string      
	PreserveFields     []string      
	DatasetParts       []DatasetPart 
	ValidationWarnings []string      
	TotalUsage         float64       
}

type DatasetPart struct {
	ID   string 
	Name string 
}

```

All models without `Description` suffix should be used for the response of the Provider API and they will be ignored in the main files.

**Note:** Please Do not add `json:"-"` tag to the models which has Description suffix. Also any model refrenced in these models.

## 5. Run the auto generators

For genertaing the all neccessary files, you should run this three commands:

```bash
go run pkg/sdk/runable/resurce_type/main.go
go run pkg/sdk/runable/steampipe_es_client_generator/main.go
go run pkg/sdk/runable/steampipe_index_map/main.go
```

## 6. Test the describer

First you nedd to add credentials to the [describer.go](./command/cmd/describer.go).
Then you can run the describer with the following command:

```bash
go run command/main.go
```

result will be saved in the output.json file.

**Note:** Next steps are optional.

## 7. Connect the describer to steampipe

You can connect the describer to steampipe. For this, you should implement the steampipe plugin.

You can use the [example plugin](./plugin/cohereai) as a reference. Example is for describing CohereAI datasets resource.

### 7.1 Add Table for the resource

Add a file with this format: `table_resource.go` in the [plugin folder](./plugin/cohereai).
You Should implement the table definition for the resource. [Example file](./plugin/cohereai/table_cohereai_datasets.go) is for describing CohereAI datasets resource.

**Note:** Transform Field should have `Description.` prefix.

### 7.2 Add the describer to the plugin

Add you function to [plugin.go] file in the [plugin.go](./plugin/cohereai/plugin.go).

## 8. Run the plugin

Build the plugin with the [Dockefile](./plugin/cohereai/Dockerfile) and run the plugin with the following command:

## 9. Test the plugin

You can import the plugin to steampipe and test it.










