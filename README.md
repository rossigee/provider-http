# provider-http

[![Build](https://github.com/rossigee/provider-http/actions/workflows/ci.yml/badge.svg)][build]
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

[build]: https://github.com/rossigee/provider-http/actions/workflows/ci.yml
[releases]: https://github.com/rossigee/provider-http/releases

## Overview

A generic [Crossplane](https://crossplane.io/) provider for managing arbitrary resources through HTTP requests. `Request` maps Crossplane's Create/Observe/Update/Delete lifecycle onto configurable HTTP calls (with JQ-based response mapping); `DisposableRequest` fires a one-shot HTTP call for cases like webhooks or notifications with no create/observe/update/delete lifecycle to track.

## Container Registry

- **Primary**: `ghcr.io/rossigee/provider-http:v1.2.1`

## Features

- **Generic HTTP CRUD mapping**: define per-action (CREATE/OBSERVE/UPDATE/REMOVE) HTTP requests, or let method (GET/POST/PUT/DELETE) imply the action
- **JQ-based templating**: request bodies, URLs, and up-to-date/removed checks are all evaluated with JQ expressions against prior responses
- **Secret injection**: reference Kubernetes Secret values in headers/body via `{{ name:namespace:key }}` syntax
- **Secret extraction**: patch response fields into new or existing Secrets via `secretInjectionConfigs`
- **DisposableRequest**: fire-and-forget HTTP calls (e.g. webhooks, notifications) with optional rollback retries and custom expected-response checks

## Getting Started

### Prerequisites

- Kubernetes with Crossplane installed
- An HTTP endpoint to manage resources against

### Installation

```bash
kubectl crossplane install provider ghcr.io/rossigee/provider-http:v1.2.1
```

### Configuration

```yaml
apiVersion: http.crossplane.io/v1alpha1
kind: ProviderConfig
metadata:
  name: http-conf
spec:
  credentials:
    source: None
```

## Usage

```yaml
apiVersion: http.crossplane.io/v1alpha2
kind: DisposableRequest
metadata:
  name: send-notification
spec:
  deletionPolicy: Orphan
  forProvider:
    url: http://flask-api.default.svc.cluster.local/v1/notify
    method: POST
    body: |
      {
        "recipient": "user@example.com",
        "subject": "Alert",
        "message": "Your action is required immediately."
      }
    headers:
      Content-Type:
        - application/json
    expectedResponse: '.body.status == "sent"'
    rollbackRetriesLimit: 5
  providerConfigRef:
    name: http-conf
```

See `examples/sample/` for a full `Request` example covering CREATE/OBSERVE/UPDATE/REMOVE mappings, JQ-based up-to-date checks, and secret extraction.

## Resource Types

| Resource | API Group | Description |
|----------|-----------|-------------|
| Request | `http.crossplane.io` (v1alpha1/v1alpha2), `http.m.crossplane.io` (v1beta1, namespaced) | Full CRUD lifecycle mapped to configurable HTTP requests |
| DisposableRequest | `http.crossplane.io` (v1alpha1/v1alpha2), `http.m.crossplane.io` (v1beta1, namespaced) | One-shot HTTP request with no ongoing lifecycle |
| ProviderConfig | `http.crossplane.io/v1alpha1` | Provider-level credentials configuration |

## Development

```bash
# Build
make build

# Test
make test

# Lint
make lint

# Generate
make generate
```

## Contributing

Issues and pull requests are welcome at [github.com/rossigee/provider-http](https://github.com/rossigee/provider-http).

## License

provider-http is under the Apache 2.0 license.

## Implementation

This provider is a native Crossplane controller that directly implements the provider APIs without using Terraform or upjet scaffolding. This approach yields smaller binaries, simpler code, and reduced dependencies.
