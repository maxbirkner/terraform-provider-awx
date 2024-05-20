# Defines the arguments that will be passed into this Dockerfile via the
# build-args for FROM commands.
ARG ALPINE_VERSION
ARG GO_VERSION

# Pull the latest golang container from upstream. Go version and
# Alpine versions are taken as --build-args during the build process.
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION}

# Defines the arguments that will be passed into this Dockerfile via the
# build-args for all remaining commands.
ARG GORELEASER_VERSION
ARG GOLANGCILINT_VERSION
ARG TERRAFORM_VERSION

# Install any additional packages that need to be included in the build container. Since
# we run yamllint as part of our linting process, it needs to be installed here.
RUN apk update && apk --no-cache add yamllint make git gcc libc-dev findutils bash curl dbus-x11 gcc-go gpg gpg-agent

# Install terraform from the specified version tag variable.
RUN wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip
RUN unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /usr/local/bin/
RUN chmod +x /usr/local/bin/terraform
RUN rm terraform_${TERRAFORM_VERSION}_linux_amd64.zip

# Install goreleaser from the specified tag in Github.
RUN wget https://github.com/goreleaser/goreleaser/releases/download/v${GORELEASER_VERSION}/goreleaser_${GORELEASER_VERSION}_x86_64.apk
RUN apk add --allow-untrusted goreleaser_${GORELEASER_VERSION}_x86_64.apk

# Install golangci-lint, also specifying the version tag to be used.
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v${GOLANGCILINT_VERSION}

# Work around new security warning introduced in Git 2.34.2
RUN git config --global --add safe.directory /workdir

# Compose will map the current development directory on the users machine to the `/source` directory
# within the container. We set the working directory to be the mapped volume.
WORKDIR /workdir
