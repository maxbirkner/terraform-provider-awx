# Dev Environment for terraform-provider-awx
This directory is a development environment to test and verify the terraform
provider on your local workstation.

If you wish to work on the provider, you'll first need to install
[Go](http://www.golang.org) on your machine (see [Requirements](../Makefile) for more information).

## Table of Contents
1. [Prerequisites](#prerequisites)
2. [Running The Tests](#running-the-tests)
3. [Example Usage](#example-usage)


### Prerequisites

When you run `make build` to generate the provider binary, it will be placed in
this directory. You will need to modify your `~/.terraformrc` file to include
this directory in the `providers` section. For example:

Add this to your `~/.terraformrc` file:
```hcl
provider_installation {

  dev_overrides {
      "registry.terraform.io/josh-silvas/awx" = "/<full-path-to-this-directory>/develop"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

Once this is defined, you can use the provider in your terraform code like so:
```hcl
terraform {
  required_providers {
    awx = {
      source = "registry.terraform.io/josh-silvas/awx"
    }
  }
}
```

> NOTE -- Skip terraform init when using provider development overrides. It is not necessary and may error unexpectedly.

### Running The Tests
At this point, you can run `terraform plan` and `terraform apply` to test the
provider from within this directory.

You can optionally run `make terraform-plan` and `make terraform-apply` to run
the test in the root directory as helper make targets.

### Example Usage

To compile the provider, run `make build`. This will build the provider and
put the provider binary in the `<this-repository-directory>/develop` directory.

To generate or update documentation, run `make docs`.

In order to run the full suite of Acceptance tests, run `make testacc`.

> *Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make lint             | Run golangci-lint on all sub-packages within docker
make _lint            | Run golangci-lint on all sub-packages
                      |
make testacc          | Run acceptance tests on all sub-packages within docker
make _testacc         | Run acceptance tests
                      |
make cli              | Launch a bash shell inside the running container.
make destroy          | Destroy the docker-compose environment and volumes
make develop          | Build the development docker image and push to registry
                      |
release               | Run goreleaser to create a release
                      |
make docs             | Run go generate to create documentation in the docs subfolder
                      |
make build            | Build the provider for local development
make terraform-apply  | Run terraform apply to test the provider
make terraform-plan   | Run terraform plan to test the provider
                      |
make tidy             | Run go mod tidy and go mod vendor
make help             | Display this help screen
```
