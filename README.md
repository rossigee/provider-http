# provider-http

**✅ BUILD STATUS: WORKING** - Successfully builds with excellent test coverage (v1.1.0)

`provider-http` is a Crossplane Provider designed to facilitate sending HTTP requests as resources.

## Features
- **HTTP Requests**: Send HTTP requests as declarative Kubernetes resources
- **Request Management**: Full lifecycle management of HTTP-based resources
- **Response Processing**: Handle and transform HTTP responses
- **Provider Status**: ✅ Builds successfully with excellent test coverage (32.7%-100%)

## Container Registry
- **Primary**: `ghcr.io/rossigee/provider-http:v1.1.0`
- **Harbor**: Available via environment configuration
- **Upbound**: Available via environment configuration

## Installation

To install `provider-http`, you have two options:

1. Using the Crossplane CLI in a Kubernetes cluster where Crossplane is installed:

   ```console
   crossplane xpkg install provider ghcr.io/rossigee/provider-http:v1.1.0
   ```

2. Manually creating a Provider by applying the following YAML:

   ```yaml
   apiVersion: pkg.crossplane.io/v1
   kind: Provider
   metadata:
     name: provider-http
   spec:
     package: "ghcr.io/rossigee/provider-http:v1.1.0"
   ```

## Supported Resources

`provider-http` supports the following resources:

- **DisposableRequest:** Initiates a one-time HTTP request. See [DisposableRequest CRD documentation](resources-docs/disposablerequest_docs.md).
- **Request:** Manages a resource through HTTP requests. See [Request CRD documentation](resources-docs/request_docs.md).

## Usage

### DisposableRequest

Create a `DisposableRequest` resource to initiate a single-use HTTP interaction:

```yaml
apiVersion: http.crossplane.io/v1alpha2
kind: DisposableRequest
metadata:
  name: example-disposable-request
spec:
  # Add your DisposableRequest specification here
```

For more detailed examples and configuration options, refer to the [examples directory](examples/sample/).

### Request

Manage a resource through HTTP requests with a `Request` resource:

```yaml
apiVersion: http.crossplane.io/v1alpha2
kind: Request
metadata:
  name: example-request
spec:
  # Add your Request specification here
```

For more detailed examples and configuration options, refer to the [examples directory](examples/sample/).

## Developing locally

Run controller against the cluster:

```
make run
```

## Run tests

```
make test
make e2e
```

## Troubleshooting

If you encounter any issues during installation or usage, refer to the [troubleshooting guide](https://docs.crossplane.io/knowledge-base/guides/troubleshoot/) for common problems and solutions.
