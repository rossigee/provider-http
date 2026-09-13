# HTTP Provider API Reference

Complete API specification for provider-http v1beta1 resources.

## API Groups

HTTP provider resources are under the `http.m.crossplane.io` API group:

- `http.m.crossplane.io/v1beta1` — ProviderConfig, Request, DisposableRequest

## Common Fields

All resources share standard Crossplane fields:

```yaml
apiVersion: http.m.crossplane.io/v1beta1
kind: <ResourceType>
metadata:
  name: <resource-name>
  namespace: <namespace>
spec:
  providerConfigRef:
    name: <provider-config-name>
  forProvider:
    # Resource-specific configuration
  deletionPolicy: Delete|Orphan     # Deletion behavior
status:
  conditions: []                     # Readiness conditions
  atProvider:
    # Observed provider state
```

## ProviderConfig

Connection configuration for HTTP provider.

### Spec

```yaml
apiVersion: http.m.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: None|Secret            # Credential source
    secretRef:
      name: <secret-name>           # Secret containing credentials
      namespace: <namespace>        # Secret namespace
      key: <key>                    # Key in secret
```

### Credential Source Options

- **None** — No credentials (public APIs)
- **Secret** — Read credentials from Kubernetes Secret

### Example

```yaml
apiVersion: http.m.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: public-api
spec:
  credentials:
    source: None
---
apiVersion: http.m.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: private-api
spec:
  credentials:
    source: Secret
    secretRef:
      name: api-token
      namespace: default
      key: token
```

## Request

Manages resources through complete HTTP lifecycle (Create/Observe/Update/Delete).

### Spec

```yaml
apiVersion: http.m.crossplane.io/v1beta1
kind: Request
metadata:
  name: my-resource
spec:
  providerConfigRef:
    name: default
  forProvider:
    mappings:                       # HTTP action mappings (required)
      - action: CREATE              # Crossplane action
        method: POST                # HTTP method
        url: http://api.example.com/resources
        body: "{...}"              # Request body
        headers: {...}             # Request headers
    
    headers: {...}                 # Default headers for all requests
    waitTimeout: "5m"              # Max wait time
    insecureSkipTLSVerify: false   # Skip TLS validation
    tlsConfig: {...}               # TLS configuration
    
    secretInjectionConfigs: [...]  # Extract response to Secret
    expectedResponseCheck: {...}   # Validate CREATE/OBSERVE response
    isRemovedCheck: {...}          # Validate DELETE confirmation
```

### Mappings

Define HTTP calls for each Crossplane action:

```yaml
mappings:
  - action: CREATE                 # Required: CREATE|OBSERVE|UPDATE|REMOVE
    method: POST                   # HTTP method: GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS
    url: "http://api/resources"    # URL (supports JQ templating)
    body: "{...}"                  # Request body (supports JQ)
    headers: {...}                # Override headers for this mapping
  
  # Observe - check resource state
  - action: OBSERVE
    method: GET
    url: "http://api/resources/{{ .status.atProvider.id }}"
  
  # Update - modify resource
  - action: UPDATE
    method: PUT
    url: "http://api/resources/{{ .status.atProvider.id }}"
    body: "{...}"
  
  # Remove - delete resource
  - action: REMOVE
    method: DELETE
    url: "http://api/resources/{{ .status.atProvider.id }}"
```

### Headers

Default headers applied to all mappings:

```yaml
headers:
  Content-Type:
    - application/json
  Authorization:
    - "Bearer {{ secrets:api-creds:default:token }}"
  X-Custom-Header:
    - value1
    - value2
```

Headers from mappings override defaults.

### Secret Injection Configs

Extract sensitive data from API response into Kubernetes Secret:

```yaml
secretInjectionConfigs:
  - name: extracted-secret         # Secret name to create/update
    namespace: default              # Secret namespace
    fieldPath: .data.token          # JQ path to extract from response
    keys:                           # Map response field to secret key
      - "token"                     # Secret key
```

### Expected Response Check

Validate API responses:

```yaml
expectedResponseCheck:
  type: DEFAULT                    # Use default validation
  # OR
  type: CUSTOM                     # Custom JQ logic
  logic: ".status == \"ok\" and .code == 200"
```

**Types:**
- `DEFAULT` — Check for HTTP 2xx status and valid JSON
- `CUSTOM` — Evaluate logic JQ expression against response body

**Logic Examples:**
```yaml
# Check specific field
logic: ".status == \"success\""

# Multiple conditions
logic: ".status == \"ok\" and .statusCode == 200"

# Nested paths
logic: ".data.ready == true"

# Negation
logic: ".error == null"

# Array check
logic: ".items | length > 0"
```

### Is Removed Check

Validate resource was deleted:

```yaml
isRemovedCheck:
  type: CUSTOM
  logic: ".notFound == true"
```

Confirms resource no longer exists after REMOVE action.

### TLS Config

Configure TLS for HTTPS:

```yaml
tlsConfig:
  insecureSkipVerify: true         # Skip certificate validation
  caBundle:
    source: Secret
    secretRef:
      name: ca-cert
      namespace: default
      key: ca.crt
```

### Payload

Shared payload configuration:

```yaml
payload:
  # (Additional payload configuration)
```

### Wait Timeout

Maximum duration to wait for requests:

```yaml
waitTimeout: "5m"                  # Duration format
```

Valid formats: `30s`, `5m`, `1h`

## DisposableRequest

Fire-and-forget one-shot HTTP request (webhooks, notifications).

### Spec

```yaml
apiVersion: http.m.crossplane.io/v1beta1
kind: DisposableRequest
metadata:
  name: webhook
spec:
  providerConfigRef:
    name: default
  forProvider:
    url: "http://webhook.example.com/notify"    # Webhook URL
    method: POST                                  # HTTP method
    body: "{...}"                                # Request body
    headers: {...}                               # Request headers
    expectedResponse:                            # Validate response
      statusCode: 200
```

## JQ Templating Reference

All string fields support JQ templating:

### Basic Syntax

- `{{ expression }}` — JQ expression
- `{{ . }}` — Current object
- `{{ .field }}` — Access field
- `{{ .nested.field }}` — Nested access

### Built-in Variables

- `.spec` — CR spec
- `.status` — CR status (after observe)
- `.metadata` — Kubernetes metadata
- `now` — Current Unix timestamp

### Common Filters

| Filter | Example | Output |
|--------|---------|--------|
| `todate` | `{{ now \| todate }}` | ISO 8601 date |
| `ascii_upcase` | `{{ .name \| ascii_upcase }}` | UPPERCASE |
| `ascii_downcase` | `{{ .name \| ascii_downcase }}` | lowercase |
| `length` | `{{ .items \| length }}` | Array length |
| `keys` | `{{ . \| keys }}` | Object keys |
| `values` | `{{ . \| values }}` | Object values |
| `tojson` | `{{ . \| tojson }}` | JSON string |
| `fromjson` | `{{ .json \| fromjson }}` | Parsed JSON |

### Conditional Logic

```jq
# If-then-else
if .status == "ready" then "yes" else "no" end

# Alternative operator
.field // "default"

# Boolean operators
(.a == true) and (.b != null)
```

### Common Patterns

```yaml
# Current timestamp
"timestamp": "{{ now | todate }}"

# Reference CR name
"name": "{{ .metadata.name }}"

# Reference status field
"resource_id": "{{ .status.atProvider.id }}"

# Upper case conversion
"code": "{{ .spec.forProvider.name | ascii_upcase }}"

# JSON encoding
"config": {{ .spec.forProvider.config | tojson }}

# Array length
"count": {{ .spec.forProvider.items | length }}

# Conditional
"env": "{{ if .metadata.namespace == "prod" then "production" else "development" end }}"
```

## Secret Reference Syntax

Reference Kubernetes Secrets in values:

```
{{ secrets:secret-name:namespace:key }}
```

### Example

```yaml
headers:
  Authorization:
    - "Bearer {{ secrets:api-creds:default:token }}"

body: |
  {
    "password": "{{ secrets:db-creds:database:password }}"
  }
```

## Status Fields

### Conditions

Standard Crossplane conditions:

| Type | Status | Reason | Description |
|------|--------|--------|-------------|
| `Ready` | True | `Available` | Resource created successfully |
| `Ready` | False | `ReconcileFailed` | Sync failed |
| `Synced` | True | `ReconcileSuccess` | Latest operation succeeded |
| `Synced` | False | `ReconcileFailed` | Latest operation failed |

### AtProvider (Observed State)

Request stores observed state from HTTP responses:

```yaml
status:
  atProvider:
    # Varies by API response
    # Any fields from HTTP response stored here
    id: resource-123
    status: active
    created_at: "2025-09-13T10:00:00Z"
```

## Deletion Policies

Control deletion behavior:

| Policy | Behavior |
|--------|----------|
| `Delete` (default) | Sends HTTP DELETE request before removing CR |
| `Orphan` | Delete only CR, don't call HTTP endpoint |

### Example

```yaml
spec:
  deletionPolicy: Orphan  # Keep resource when deleting CR
```

## Management Policies

Control which operations are allowed:

```yaml
spec:
  managementPolicies:
    - Create
    - Update
    - Delete
```

Options: `Create`, `Update`, `Delete`

Default: `[Create, Update, Delete]`

Example - read-only:

```yaml
spec:
  managementPolicies: []  # No operations allowed
```

## Examples

### Simple POST Webhook

```yaml
apiVersion: http.m.crossplane.io/v1beta1
kind: DisposableRequest
metadata:
  name: notify-slack
spec:
  providerConfigRef:
    name: default
  forProvider:
    url: "https://hooks.slack.com/services/xxx/yyy/zzz"
    method: POST
    headers:
      Content-Type:
        - application/json
    body: |
      {
        "text": "Deployment {{ .metadata.name }} created in {{ .metadata.namespace }}"
      }
```

### Full CRUD Resource

```yaml
apiVersion: http.m.crossplane.io/v1beta1
kind: Request
metadata:
  name: managed-api-resource
spec:
  providerConfigRef:
    name: default
  forProvider:
    mappings:
      - action: CREATE
        method: POST
        url: "http://api.example.com/resources"
        body: |
          {
            "name": "{{ .metadata.name }}",
            "config": {{ .spec.forProvider.config | tojson }}
          }
      
      - action: OBSERVE
        method: GET
        url: "http://api.example.com/resources/{{ .status.atProvider.id }}"
      
      - action: UPDATE
        method: PUT
        url: "http://api.example.com/resources/{{ .status.atProvider.id }}"
        body: |
          {
            "config": {{ .spec.forProvider.config | tojson }}
          }
      
      - action: REMOVE
        method: DELETE
        url: "http://api.example.com/resources/{{ .status.atProvider.id }}"
    
    headers:
      Authorization:
        - "Bearer {{ secrets:api-token:default:token }}"
      Content-Type:
        - application/json
    
    expectedResponseCheck:
      type: CUSTOM
      logic: ".statusCode == 200"
```

### Secret Extraction

```yaml
apiVersion: http.m.crossplane.io/v1beta1
kind: Request
metadata:
  name: get-api-token
spec:
  forProvider:
    mappings:
      - action: CREATE
        method: POST
        url: "http://auth.example.com/tokens"
        body: '{"client_id":"app"}'
      
      - action: OBSERVE
        method: GET
        url: "http://auth.example.com/tokens/{{ .status.atProvider.token_id }}"
    
    secretInjectionConfigs:
      - name: api-token
        namespace: default
        fieldPath: .access_token
        keys:
          - "token"
```

## Compatibility

- **Kubernetes Version**: 1.20+
- **Crossplane Version**: 1.14+
- **HTTP Methods**: GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS

## Rate Limits

HTTP provider respects target API rate limits. Implement custom retry logic via `expectedResponseCheck` if needed.

## See Also

- [Getting Started Guide](getting-started.md)
- [JQ Manual](https://stedolan.github.io/jq/manual/)
- [Crossplane Documentation](https://docs.crossplane.io)
