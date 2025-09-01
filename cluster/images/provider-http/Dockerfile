FROM gcr.io/distroless/static@sha256:a01d47d4036cae5a67a9619e3d06fa14a6811a2247b4da72b4233ece4efebd57

# OpenContainers Image Spec Labels
LABEL org.opencontainers.image.title="provider-http"
LABEL org.opencontainers.image.description="Crossplane provider for Generic HTTP request resources"
LABEL org.opencontainers.image.vendor="rossigee"
LABEL org.opencontainers.image.licenses="Apache-2.0"
LABEL org.opencontainers.image.source="https://github.com/rossigee/provider-http"
LABEL org.opencontainers.image.url="https://github.com/rossigee/provider-http"
LABEL org.opencontainers.image.documentation="https://github.com/rossigee/provider-http/blob/master/README.md"

ARG TARGETOS
ARG TARGETARCH

ADD bin/$TARGETOS\_$TARGETARCH/provider /usr/local/bin/crossplane-http-provider

USER 65532
ENTRYPOINT ["crossplane-http-provider"]
