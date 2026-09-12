# Provider HTTP Documentation

A Crossplane v2 provider mapping configurable HTTP requests to managed-resource lifecycles. All resources are namespaced `http.m.crossplane.io/v1beta1` with multi-tenancy support.

## Resource Documentation

| Resource | API Group | Description |
|----------|-----------|-------------|
| Request | `http.m.crossplane.io/v1beta1` | Full CRUD lifecycle mapped to configurable HTTP requests |
| DisposableRequest | `http.m.crossplane.io/v1beta1` | One-shot HTTP request with no ongoing lifecycle |
| ProviderConfig | `http.m.crossplane.io/v1beta1` | Provider-level credentials (cluster-scoped) |

See `examples/sample/` for full Request coverage (CREATE/OBSERVE/UPDATE/REMOVE mappings, JQ up-to-date checks, secret extraction).

## API Coverage Gaps

By design this provider is transport-generic rather than API-covering, but notable missing capabilities: response pagination/fan-out across pages, retry/backoff policy per mapping, request signing (HMAC/OAuth2 flows beyond static headers), webhook-triggered reconcile, and dry-run/diff preview of mappings.
