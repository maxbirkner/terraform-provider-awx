# Terraform Provider AWX

This is a Terraform provider for AWX (Ansible Tower) that allows you to manage
resources like organizations, projects, inventories, job templates, and more.

The AWX Terraform provider relies on the [goawx](https://github.com/josh-silvas/goawx) SDK.

_This repository is a fork from [denouche/terraform-provider-awx](https://github.com/josh-silvas/terraform-provider-awx)
to extend functionality and support._

This repository is built on Terraform scaffolding for providers and contains the following:

- A resource and a data source (`internal/provider/`),
- Examples (`examples/`) and generated documentation (`docs/`),
- Miscellaneous meta files.

# Table of Contents
1. [Requirements](#requirements)
2. [Building The Provider](#building-the-provider)
3. [Adding Dependencies](#adding-dependencies)
4. [Using The Provider](#using-the-provider)
5. [Developing The Provider](#developing-the-provider)

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command or the Makefile `build` target.

```shell
make build
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up-to-date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get -u github.com/author/dependency
go mod tidy && go mod vendor
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

### Provider Configuration

```hcl
terraform {
    required_version = ">= 1.5.3"
    required_providers {
        awx = {
            source  = "josh-silvas/awx"
            version = "0.24.3"
        }
    }
}

// Configure the AWX provider. This example relies on AWX_HOSTNAME and AWX_TOKEN to exist
// in the environment. If not, you will need to specify them here.
provider "awx" {
    hostname = "https://awx.example.com"
    token    = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
    insecure = true
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org)
installed on your machine (see [Requirements](#requirements) above).

See the [DEVELOPMENT](develop/README.md) documentation for more information.


## Resources

* [Josh Silvas](mailto:josh@jsilvas.com)

