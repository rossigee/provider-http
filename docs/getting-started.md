# Getting Started with Provider HTTP

Complete guide to installing and using provider-http v1beta1 with Crossplane for generic HTTP resource management.

## Overview

Provider HTTP enables managing arbitrary resources through HTTP requests. Use it to:
- Integrate with REST APIs that lack a native Crossplane provider
- Create fire-and-forget webhooks and notifications
- Manage resources through HTTP endpoints
- Synchronize external systems with Kubernetes declaratively

## Prerequisites

- **Kubernetes cluster** with Crossplane v1.14+ installed
- **HTTP endpoint** to manage resources against
- `kubectl` configured to access your cluster
- Basic understanding of HTTP methods (GET, POST, PUT, DELETE)

## Installation

### 1. Install Crossplane

```bash
helm repo add crossplane-stable https://charts.crossplane.io/stable
helm install crossplane \
  crossplane-stable/crossplane \
  -n crossplane-system \
  --create-namespace
```

Wait for Crossplane to be ready:
```bash
kubectl wait -n crossplane-system --for=condition=Ready pods -l app.kubernetes.io/instance=crossplane --timeout=300s
```

### 2. Install Provider HTTP

```bash
# Using Crossplane CLI
kubectl crossplane install provider ghcr.io/rossigee/provider-http:v1.2.1

# Or using Helm
helm repo add crossplane-contrib https://charts.crossplane.io/contrib
helm install provider-http \
  crossplane-contrib/provider-http \
  -n crossplane-system \
  --version ">=1.2.1"
```

Verify installation:
```bash
kubectl get providers
kubectl describe provider provider-http
```

## Configuration

### Create ProviderConfig

Provider HTTP can work with or without credentials:

```yaml
apiVersion: http.m.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: None
```

For APIs requiring authentication, create a secret first:

```bash
kubectl create secret generic api-credentials \
  -n crossplane-system \
  --from-literal=token=your-api-token
```

Then reference it in ProviderConfig:

```yaml
apiVersion: http.m.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: api-auth
spec:
  credentials:
    source: Secret
    secretRef:
      name: api-credentials
      namespace: crossplane-system
      key: token
```

Apply the configuration:
```bash
kubectl apply -f provider-config.yaml
```

## Resources

### Resource Types

Provider HTTP provides two main resource types:

1. **Request** — Full lifecycle management (Create/Observe/Update/Delete)
2. **DisposableRequest** — Fire-and-forget one-shot HTTP calls

## First Resource: Simple REST API Call

### 1. Create a DisposableRequest (Webhook)

Send a one-shot HTTP POST to notify an external system:

```yaml
apiVersion: http.m.crossplane.io/v1beta1
kind: DisposableRequest
metadata:
  name: deployment-webhook
  namespace: default
spec:
  providerConfigRef:
    name: default
  forProvider:
    url: http://webhook.example.com/v1/notify
    method: POST
    body: |
      {
        "event": "deployment-created",
        "timestamp": "{{ now | todate }}"
      }
    headers:
      Content-Type:
        - application/json
      Authorization:
        - Bearer token123
```

Apply:
```bash
kubectl apply -f webhook.yaml
```

Monitor:
```bash
kubectl describe disposablerequest deployment-webhook
kubectl get disposablerequest deployment-webhook -o jsonpath='{.status.conditions}'
```

### 2. Create a Request (Full Lifecycle)

Manage a resource through complete CRUD lifecycle:

```yaml
apiVersion: http.m.crossplane.io/v1beta1
kind: Request
metadata:
  name: managed-resource
  namespace: default
spec:
  providerConfigRef:
    name: default
  forProvider:
    mappings:
      # Create - POST to create resource
      - action: CREATE
        method: POST
        url: http://api.example.com/v1/resources
        body: |
          {
            "name": "{{ .spec.forProvider.name }}",
            "description": "{{ .spec.forProvider.description }}"
          }
      # Observe - GET to check resource status
      - action: OBSERVE
        method: GET
        url: http://api.example.com/v1/resources/{{ .status.atProvider.id }}
      # Update - PUT to update resource
      - action: UPDATE
        method: PUT
        url: http://api.example.com/v1/resources/{{ .status.atProvider.id }}
        body: |
          {
            "name": "{{ .spec.forProvider.name }}",
            "description": "{{ .spec.forProvider.description }}"
          }
      # Remove - DELETE to delete resource
      - action: REMOVE
        method: DELETE
        url: http://api.example.com/v1/resources/{{ .status.atProvider.id }}
    
    headers:
      Content-Type:
        - application/json
      Authorization:
        - "Bearer token123"
```

Apply:
```bash
kubectl apply -f managed-resource.yaml
```

## Common Patterns

### Pattern 1: Simple Status Check

Poll an endpoint and extract response into status:

```yaml
apiVersion: http.m.crossplane.io/v1beta1
kind: Request
metadata:
  name: status-monitor
spec:
  providerConfigRef:
    name: default
  forProvider:
    mappings:
      - action: OBSERVE
        method: GET
        url: http://monitoring.example.com/api/status
    
    expectedResponseCheck:
      type: CUSTOM
      logic: .status == "healthy"
```

### Pattern 2: Secret Injection from API Response

Extract sensitive data from API response into a Kubernetes Secret:

```yaml
apiVersion: http.m.crossplane.io/v1beta1
kind: Request
metadata:
  name: api-token-fetcher
spec:
  providerConfigRef:
    name: default
  forProvider:
    mappings:
      - action: CREATE
        method: POST
        url: http://auth.example.com/v1/tokens
        body: '{"client_id":"my-app"}'
      
      - action: OBSERVE
        method: GET
        url: http://auth.example.com/v1/tokens/{{ .status.atProvider.token_id }}
    
    # Extract token from response into Secret
    secretInjectionConfigs:
      - name: api-token
        namespace: default
        fieldPath: .token
        keys:
          - "token"
```

Result: Token is stored in Secret `api-token` in `default` namespace

### Pattern 3: Headers from Secret

Reference credentials stored in Kubernetes Secret:

```yaml
apiVersion: http.m.crossplane.io/v1beta1
kind: Request
metadata:
  name: authenticated-request
spec:
  forProvider:
    mappings:
      - action: CREATE
        method: POST
        url: http://api.example.com/v1/resources
    
    headers:
      Authorization:
        - "Bearer {{ secrets:api-credentials:default:token }}"
      X-API-Key:
        - "{{ secrets:api-keys:default:key1 }}"
```

Syntax: `{{ secrets:secret-name:namespace:key }}`

### Pattern 4: Conditional Removal

Check if resource is already deleted:

```yaml
apiVersion: http.m.crossplane.io/v1beta1
kind: Request
metadata:
  name: conditional-delete
spec:
  forProvider:
    mappings:
      - action: REMOVE
        method: DELETE
        url: http://api.example.com/v1/resources/{{ .status.atProvider.id }}
    
    # Check if removed - custom logic returns true if deleted
    isRemovedCheck:
      type: CUSTOM
      logic: .notFound == true
```

### Pattern 5: Multiple Endpoints (Resource Composition)

Manage related resources through multiple API calls:

```yaml
apiVersion: http.m.crossplane.io/v1beta1
kind: Request
metadata:
  name: complex-resource
spec:
  forProvider:
    mappings:
      # Create parent resource
      - action: CREATE
        method: POST
        url: http://api.example.com/v1/parents
        body: '{"name":"{{ .spec.forProvider.name }}"}'
      
      # Create child resource using parent ID
      - action: CREATE
        method: POST
        url: http://api.example.com/v1/children
        body: |
          {
            "parent_id": "{{ .status.atProvider.parent_id }}",
            "config": {{ .spec.forProvider.config | tojson }}
          }
      
      # Observe parent
      - action: OBSERVE
        method: GET
        url: http://api.example.com/v1/parents/{{ .status.atProvider.parent_id }}
```

## Advanced Features

### JQ Templating

All fields (URL, body, headers) support JQ expressions:

```yaml
body: |
  {
    "timestamp": "{{ now | todate }}",
    "upper": "{{ .spec.forProvider.name | ascii_upcase }}",
    "id": "{{ .status.atProvider.resource_id }}",
    "values": {{ .spec.forProvider.values | tojson }},
    "count": {{ .spec.forProvider.items | length }}
  }
```

Common JQ filters:
- `todate` — Convert timestamp to ISO date
- `ascii_upcase` / `ascii_downcase` — String case conversion
- `tojson` — Convert to JSON string
- `length` — Get array/string length
- `keys` — Get object keys
- `values` — Get object values

### TLS Configuration

Configure TLS for HTTPS endpoints:

```yaml
forProvider:
  tlsConfig:
    insecureSkipVerify: true  # Skip cert validation (not recommended)
    caBundle:
      source: Secret
      secretRef:
        name: ca-cert
        namespace: default
        key: ca.crt
```

### Custom Response Validation

Use custom JQ logic to validate API responses:

```yaml
expectedResponseCheck:
  type: CUSTOM
  logic: ".status == \"success\" and .code == 200"
```

The logic is evaluated against the API response. Returns `true` if valid.

## Troubleshooting

### Request Fails with Connection Error

```bash
# Check provider logs
kubectl logs -n crossplane-system -l pkg.crossplane.io/provider=provider-http

# Verify endpoint is reachable
kubectl run -it --rm debug --image=curlimages/curl --restart=Never -- \
  curl http://api.example.com/v1/health
```

### Response Parsing Fails

Check the actual API response format:
```bash
curl -v http://api.example.com/v1/resources/123
```

Update JQ expressions in mappings to match actual response structure.

### Secret Reference Not Working

Verify secret exists:
```bash
kubectl get secret api-credentials -n default
kubectl get secret api-credentials -n default -o jsonpath='{.data.token}' | base64 -d
```

Verify syntax: `{{ secrets:secret-name:namespace:key }}`

### Headers Not Being Sent

Check provider config credentials source - must be set correctly:
```bash
kubectl describe providerconfig default
```

## Best Practices

### 1. Use Secrets for Sensitive Data

Never embed credentials in manifests:

```yaml
# ❌ Bad
headers:
  Authorization:
    - "Bearer my-secret-token"

# ✅ Good
headers:
  Authorization:
    - "{{ secrets:api-creds:default:token }}"
```

### 2. Validate Responses Explicitly

Always define expected response checks:

```yaml
expectedResponseCheck:
  type: CUSTOM
  logic: ".status == \"ok\" and .statusCode == 200"
```

### 3. Set Appropriate Timeouts

For long-running operations:

```yaml
forProvider:
  waitTimeout: "10m"  # Wait up to 10 minutes
```

### 4. Use Meaningful Names

```yaml
metadata:
  name: webhook-notify-slack  # Descriptive name
```

### 5. Document Complex Mappings

Use comments in body for clarity:

```yaml
body: |
  {
    # Resource name from CR spec
    "name": "{{ .spec.forProvider.name }}",
    # Current timestamp
    "created": "{{ now | todate }}"
  }
```

## Examples Repository

Complete working examples available at:
https://github.com/rossigee/provider-http/tree/master/examples

## Next Steps

1. **Explore Examples**: Check examples/ directory for real-world patterns
2. **Advanced JQ**: Learn more at https://stedolan.github.io/jq/
3. **API Reference**: See [API_REFERENCE.md](API_REFERENCE.md) for complete specification
4. **Troubleshooting**: Check [Troubleshooting Guide](troubleshooting.md)

## Getting Help

- **GitHub Issues**: https://github.com/rossigee/provider-http/issues
- **Crossplane Slack**: https://slack.crossplane.io
- **JQ Documentation**: https://stedolan.github.io/jq/
