# OpenGovernance Describer Template

## Introduction

This document is a GitHub repository temmplate for write describers for any Provider you want.

## Instructions

### 1. Create a new repository using this template

First, you need to fork this repository to your account. Then, you can create a new repository using this template.

### 2. Fill the environment variables

Fill the environment variables in the `.env` file with the information of the Provider you want to describe.
After that please run following command to export the environment variables:

```bash
./export_env.sh
```

### 3. Fill the resources

Complete the resources same as the example in the [resources-types.json](./SDK/runable/resourceType/resource-types.json).


### 4. Fill the index map

Complete the index map same as the example in the [table_index_map.go](./steampipe/table_index_map.go).

### 5. Copy the Steampipe plugin

Copy your  Steampipe plugin to root directory of the repository.

### 6. Complete authentication and describer files

Complete the authentication functions for your provider in the [config.go](./provider/config.go).
After that Please complete all neccessary functions which tagged with `TODO` in the [provider](./provider/) and [describer](./describer/) directories.

### 7. Run the Template

Just run the following command to run the template:

```bash
go run main.go
```
