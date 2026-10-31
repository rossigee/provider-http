# Migrating Request and DisposableRequest to namespaced groups

Request and DisposableRequest were cluster-scoped and served from
`http.m.crossplane.io/v1beta1`. They are now namespaced and served from their own
groups:

- `Request`: `request.http.m.crossplane.io/v1beta1`
- `DisposableRequest`: `disposablerequest.http.m.crossplane.io/v1beta1`

The ProviderConfig stays on `http.m.crossplane.io/v1beta1`.

Installing this package removes the old CRDs. Deleting a CRD deletes every
custom resource under it, and each resource's deletion policy then runs against
the remote endpoint. Migrate before upgrading.

## Steps

1. Set `spec.deletionPolicy: Orphan` on every existing Request and
   DisposableRequest, so removing the old objects does not send any remote call.
2. Record the existing resources (`kubectl get <kind> -o yaml`), then delete the
   old ones once they are orphaned.
3. Upgrade the provider package.
4. Recreate each Request under `request.http.m.crossplane.io/v1beta1` with
   `metadata.namespace` set. DisposableRequests are one-shot and are not carried
   over; recreate them only if you need to run them again.
