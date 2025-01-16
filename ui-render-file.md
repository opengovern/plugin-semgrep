# Understanding the JSON UI Render File

**Note: Please follow the JSON UI Render file structure and fill in the all fields to ensure the integration is displayed correctly in the OpenGovernance platform.**

The JSON UI Render file is divided into four main sections:

1. **General Information**: Contains metadata about the integration.
2. **Discover**: Defines the structure of credentials and integrations.
3. **Render**: Specifies how to display content in different views.
4. **Actions**: Lists possible user actions on credentials and integrations.

## General Information

```json
{
  "integration_type_id": "aws_cloud_account",
  "integration_name": "AWS Cloud Account",
  "help_text_md": "AWS Cloud Account integration facilitates secure connections to your AWS resources. [Documentation](https://docs.aws.amazon.com).",
  "platform_documentation": "https://docs.aws.amazon.com",
  "provider_documentation": "https://aws.amazon.com",
  "icon": "aws.svg",
  ...
}
```

- **integration_type_id**: Unique identifier (e.g., `aws_cloud_account`).
- **integration_name**: Human-readable name.
- **help_text_md**: Markdown-formatted help text with documentation links.
- **platform_documentation** & **provider_documentation**: URLs to official documentation.
- **icon**: Filename for the integration's icon.

## Discover

### Credentials

Defines the authentication details required for the integration.

```json
"discover": {
  "credentials": [
    {
        "type": "aws_single_account",
        "label": "AWS Single Account",
        "priority": 1,
        "fields": [
        {
          "name": "aws_access_key_id",
          "label": "AWS Access Key ID",
          "inputType": "text",
          "required": true,
          "order": 1,
          "validation": {
              "pattern": "^[A-Z0-9]{20}$",
              "errorMessage": "AWS Access Key ID must be a 20-character uppercase alphanumeric string."
          },
          "info": "Your AWS Access Key ID.",
          "external_help_url": "https://docs.aws.amazon.com/access-key-id"
        },
        ...
      ]

    },
    ...
  ]
}
```

- **credentials**: Contains different credential types as Object (e.g., `aws_single_account`).
- **fields**: Array defining each credential field with properties like `name`, `label`, `inputType`, `required`, `order`, `validation`, `info`, and `external_help_url`.
- **fields[x].type** Supported input types: `text`, `password`, `file`.

### Integrations

Defines the structure of the integration objects managed within OpenGovernance.

```json
{
"discover": {
  ...
  "integrations": [ {
        "label": "AWS Cloud Account",
        "type": "aws_cloud_account",
        "priority": 1,
        "fields": [
          {
            "name": "uuid",
            "label": "Integration UUID",
            "fieldType": "text",
            "required": true,
            "order": 1,
            "info": "Unique identifier (UUID) for the integration."
          },
          ...
          ]
  }
  ...
  ]
}
}
```

- **integrations**: Contains different integration types Object (e.g., `aws_cloud_account`).
- **fields**: Array defining each integration field with properties similar to credentials.

## Render

Specifies how the integration's content is displayed in various contexts.

```json
{
"render": {
    "credentials": {
    "defaultPageSize": 10,
    "fields": [
      ...
    ]
    },
    "integrations": {
      "defaultPageSize": 15,
      "fields": [
        {
          "name": "name",
          "label": "Name",
          "fieldType": "text",
          "order": 1,
          "sortable": true,
          "filterable": true,
          "info": "Name.",
          "detail": false,
          "detail_order": 1
        },
        {
          "name": "state",
          "label": "State",
          "fieldType": "status",
          "order": 3,
          "sortable": true,
          "filterable": true,
          "detail": false,
          "detail_order": 3,
          "info": "Current state of the Azure Subscription integration.",
          "statusOptions": [
            {
              "value": "ACTIVE",
              "label": "Active",
              "color": "green"
            },
           ...
          ]
        }
        ...
      ]
    }
}
}
```

- **defaultPageSize**: Number of items per page.
- **fields**: Array of fields to show in the list view with properties like `sortable`, `filterable`, `detail`,`detail_order`.
- **detail**: Display fields in the detail view. If `detail` is `true`, the field will be shown in the detail view.
- **detail_order**: Order of fields in the detail view. If `detail` is `false` or not present, this field is ignored.
- **fieldType**: Supported field types: `text`, `status`, `date`.
- **statusOptions**: Array of status options with `value`, `label`, and `color` properties. If `fieldType` is `status`, this property is required.

## Actions

Defines the actions users can perform on credentials and integrations.

```json
{
"actions": {
  "credentials": [
    {
      "type": "view",
      "label": "View"
    },
    {
      "type": "update",
      "label": "Update",
      "editableFields": [
        ...
      ]
    },
    {
      "type": "delete",
      "label": "Delete",
      "confirm": { ... }
    }
  ],
  "integrations": [
    {
      "type": "view",
      "label": "View"
    },
    {
      "type": "update",
      "label": "Update",
      "editableFields": [
        ...
      ]
    },
    {
      "type": "delete",
      "label": "Delete",
      "confirm": { ... }
    },
    {
      "type": "health_check",
      "label": "Health Check",
      "tooltip": "Run a health check on the integration to verify connectivity and configuration."
    }
  ]
}
}
```

- **type**: Action identifier (e.g., `view`, `update`, `delete`).
- **label**: Display name for the action.
- **editableFields**: Fields that can be edited during an update action.
- **confirm**: Confirmation dialog configuration for delete actions.

---

**Note: The JSON UI Render file is a crucial part of the integration development process. It defines how the integration is displayed and interacted with in the OpenGovernance platform.**
