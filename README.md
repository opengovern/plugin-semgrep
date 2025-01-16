# OpenGovernance Describer Template

This document is a GitHub repository temmplate for write describers for any Provider you want.

- [OpenGovernance Describer Template](#opengovernance-describer-template)
- [Instructions](#instructions)
	- [1. Create a new repository using this template](#1-create-a-new-repository-using-this-template)
	- [2. Fill the Provider information](#2-fill-the-provider-information)
		- [2.1. Provider information](#21-provider-information)
		- [2.2. Integration information](#22-integration-information)
	- [3. Fill the describer wrapper](#3-fill-the-describer-wrapper)
	- [4. Create the describer file and implement the describer](#4-create-the-describer-file-and-implement-the-describer)
		- [4.1 Create the describer file](#41-create-the-describer-file)
		- [4.2 Implement the describer](#42-implement-the-describer)
		- [4.3 Fill model](#43-fill-model)
		- [4.4 Fill resource-types.json](#44-fill-resource-typesjson)
	- [5. Run the auto generators](#5-run-the-auto-generators)
	- [6. Test the describer](#6-test-the-describer)
	- [7. Connect the describer to steampipe](#7-connect-the-describer-to-steampipe)
		- [7.1 Add Table for the resource](#71-add-table-for-the-resource)
		- [7.2 Add the describer to the plugin](#72-add-the-describer-to-the-plugin)
	- [8. Connect the describer to opencomply ui](#8-connect-the-describer-to-opencomply-ui)
		- [8.1 Discover integrations](#81-discover-integrations)
		- [8.2 Health Check integration](#82-health-check-integration)
		- [8.3 Complete Interfaces](#83-complete-interfaces)
		- [8.4 UI Spec](#84-ui-spec)
	- [9 Test Plugins](#9-test-plugins)


# Instructions

## 1. Create a new repository using this template

First, you need to fork this repository to your account. Then, you can create a new repository using this template.

## 2. Fill the Provider information

Fill the information of the Provider you want to describe in the [global folder](./global).

### 2.1. Provider information

Fill the Credential information of the Provider in the [configs.go](./global/configs.go) file.

```go
package global

type IntegrationCredentials struct {
	// TODO
}
```

### 2.2. Integration information

Fill the Integration information of the Provider in the [configs.go](./global/configs.go) file.

```go
const (
IntegrationTypeLower = "template"                                    // example: aws, azure
IntegrationName      = integration.Type("template,github")          // example: aws_account, github_account
OGPluginRepoURL      = "github.com/opengovern/og-describer-template" // example: github.com/opengovern/og-describer-aws
)
```

## 3. Fill the describer wrapper

Fill the describer wrapper in the [describer_wrapper.go](./discovery/provider/descrdescriber_wrapperiber.go) file.

You should implement two functions:

DescribeListByProvider: This function should return a list of resources of the Provider.

DescribeSingleByProvider: This function should return a single resource of the Provider.

Theses functions are wrapper for the describer any resource of the Provider.

## 4. Create the describer file and implement the describer

### 4.1 Create the describer file

Create a new file in the [describers folder](./discovery/describers/) with the name of the resource you want to describe.

### 4.2 Implement the describer

Implement the describer in the file you created in the previous step.
Implement two List and Get functions:

List: This function should return a list of resources of the Provider.

Get: This function should return a single resource of the Provider.

**Note:** You can use the [example describer](./discovery/describers/example.go) as a reference. Example is for describing CohereAI datasets resource.

### 4.3 Fill model

You should fill the model of the resource in the [models.go](./discovery/provider/models.go).

You can define models for the resource you want to describe. The main model of the resource should have `Description` suffix.

```go
type ArtifactDockerFileDescription struct {
Sha  *string
Name *string
LastUpdatedAt *string
HTMLURL *string
DockerfileContent       string
DockerfileContentBase64 *string
Repository              map[string]interface{}
Images                  []string 
}
```

### 4.4 Fill resource-types.json

You should fill the resource-types.json file in the [resource-types.json](./global/maps/resource-types.json) folder.

```json
[
 {
   "ResourceName": "Github/Artifact/DockerFile",
   "Tags": {
     "category": ["artifact_dockerfile"]
   },
   "ListDescriber": "DescribeByIntegration(describers.ListType)",
   "GetDescriber": "",
   "SteampipeTable": "template_artifact_dockerfile",
   "Model": "ArtifactDockerFile",
   "Params": [
     {
       "Name": "repository",
       "Description": "Please provide the repo name (i.e. internal-tools)",
       "Required": false
     },
     {
       "Name": "organization",
       "Description": "Please provide the organization name",
       "Required": false
     }
   ]
 }
]
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

You can use the [example plugin](./cloudql/template) as a reference. Example is for describing CohereAI datasets resource.

### 7.1 Add Table for the resource

Add a file with this format: `table_template_resource.go` in the [plugin folder](./cloudql/template).
You Should implement the table definition for the resource. [Example file](./cloudql/template/table_template_artifact_dockerfile.go) is for describing CohereAI datasets resource.

**Note:** Transform Field should have `Description.` prefix.

### 7.2 Add the describer to the plugin

Add your function to [plugin.go] file in the [plugin.go](./cloudql/template/plugin.go).

## 8. Connect the describer to opencomply ui

You can connect the describer to opencomply ui. For this, you should follow next steps.

### 8.1 Discover integrations

write discovery function to find all integrations with the given credentials in the [discovery.go](./platform/discovery.go) file.

### 8.2 Health Check integration

write health check function to check the health of the integration in the [healthcheck.go](./platform/healthcheck.go) file.

### 8.3 Complete Interfaces
Complete discovery and healthCheck functions in the [integration.go](./platform/integration.go) file.

### 8.4 UI Spec

for rendering the integration in the UI, you should write the UI spec in the [ui_spec.json](./platform/constants/ui-spec.json) file.

You can follow guides on the [helper](spec-helper.md)   file for writing the UI spec.
Also there is an example for write a UI spec for digitalocean in the [example](how-to-digital-ocean) file.

## 9 Test Plugins

Change the [build.yaml.txt](.github/workflows/build.yaml.txt) to `build.yaml` and change the names of the plugins in the file.

Then you can test the plugins with runing action on the github.
