# Step-by-Step: Creating a DigitalOcean Integration

## 1. Define General Information

```json
{
  "integration_type_id": "digitalocean_account",
  "integration_name": "DigitalOcean Account",
  "help_text_md": "DigitalOcean Account integration facilitates secure connections to your DigitalOcean resources. [Documentation](https://docs.digitalocean.com).",
  "platform_documentation": "https://docs.digitalocean.com",
  "provider_documentation": "https://digitalocean.com",
  "icon": "digitalocean.svg",
  ...
}
```

## 2. Specify Definitions

### Credentials

```json
{
  "discover": {
      "credentials":
          [
            {
              "label": "DigitalOcean API Token",
              "priority": 1,
              "type": "do_api_token",
              "fields": [
                  {
                    "name": "do_api_token",
                    "label": "DigitalOcean API Token",
                    "inputType": "password",
                    "required": true,
                    "order": 1,
                    "validation": {
                      "pattern": "^[a-zA-Z0-9]{32}$",
                      "errorMessage": "DigitalOcean API Token must be a 32-character alphanumeric string."
                    },
                    "info": "Your DigitalOcean API Token.",
                    "external_help_url": "https://docs.digitalocean.com/reference/api/create-personal-access-token/"
                  }
              ]
            }
            ...
          ]
  }
}
```

### Integrations

```json
{
  "discover": {
    ...
    "integrations": 
      [
        {
        "label": "DigitalOcean Account",
        "type": "digitalocean_account",
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

## 3. Configure Render Section



### List

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
      ...
    ]
      }
  }
}
```


## 4. Set Up Actions

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
        "editableFields": ["do_api_token", "account_id", "additional_notes"]
      },
      {
        "type": "delete",
        "label": "Delete",
        "confirm": {
        "message": "Are you sure you want to delete this credential? This action cannot be undone.",
        "condition": {
          "field": "integration_count",
          "operator": "==",
          "value": 0,
          "errorMessage": "Credential cannot be deleted because it is used by active integrations."
        }
        }
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
          "account_name",
          "region",
          "state",
          "additional_notes"
        ]
      },
      {
        "type": "delete",
        "label": "Delete",
        "confirm": {
          "message": "Are you sure you want to delete this integration? This action cannot be undone."
        }
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

### Actions Overview

- **Credentials and Integrations:** Define actions such as `view`, `update`, `delete`, and `health_check` with relevant configurations.
