# Grammar mode: syntax=source.shell
# This image contains the package-registry binary and packages storage.

FROM alpine:3.15 as base

ARG DOCKER_IMAGE_BASE="Alpine 3.15"
ARG BUILD_DATE=""
ARG REGISTRY_VER="v1.8.0"
ENV REGISTRY_VER="$REGISTRY_VER"

# Build binary
FROM base AS builder

ENV GO111MODULE=on
# Download build dependencies
RUN set -eux; \
      apk add --no-cache go git
# Build app
RUN set -eux; \
      git clone https://github.com/elastic/package-registry -b "$REGISTRY_VER"; \
      git clone https://github.com/elastic/package-storage -b production; \
      cd package-registry; \
      go build ./


# Run binary
FROM base

LABEL maintainer="$DOCKER_IMAGE_MAINTAINER" \
      org.label-schema.build-date="$BUILD_DATE" \
      org.label-schema.description="Elastic package registry v$REGISTRY_VER" \
      org.label-schema.docker.cmd="docker run --rm -p 8080:8080 -d $DOCKER_IMAGE:$REGISTRY_VER" \
      org.label-schema.name="$DOCKER_IMAGE" \
      org.label-schema.schema-version="1.0" \
      org.label-schema.url="https://git.intermb.ru/$DOCKER_IMAGE" \
      org.label-schema.vcs-url="https://git.intermb.ru/$DOCKER_IMAGE" \
      org.label-schema.version="$REGISTRY_VER-$DOCKER_IMAGE_BASE" \
      versions.os="$DOCKER_IMAGE_BASE"

# Get dependencies
RUN set -eux \
    && apk add --no-cache \
        zip \
        rsync \
        curl \
        shared-mime-info

# Start registry when container is run an expose it on port 8080
EXPOSE 8080
ENV EPR_ADDRESS=0.0.0.0:8080
ENV EPR_DISABLE_PACKAGE_VALIDATION=true

# Move binary from the builder image
COPY --from=builder /package-registry/package-registry /package-registry/package-registry
# Get in config which expects packages in /packages
COPY --from=builder /package-storage/deployment/package-registry.yml /package-registry/config.yml
# Get packages
COPY --from=builder /package-storage/packages /packages/production

# Change to new working directory
WORKDIR /package-registry
# Smoke test
RUN ./package-registry -dry-run


HEALTHCHECK --interval=10s --retries=3 CMD curl --silent --fail localhost:8080/health || exit 1
ENTRYPOINT ["./package-registry"]
