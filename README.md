# Terraform Provider for metalstack.cloud

Manage the lifecycle of your bare-metal Kubernetes clusters on [metalstack.cloud](https://metalstack.cloud) using Terraform.

> **Note:** this project is in an early development stage. It might break in the future and getting started might be cumbersome.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building the provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install
```

## Using the provider

First, let's add the provider to your project:

```terraform
terraform {
  required_providers {
    metal = {
      source = "metal-stack-cloud/metal"
    }
  }
}
```

To obtain an `api token` for creating resources, visit [metalstack.cloud](https://metalstack.cloud). Head to the the `Access Tokens` section and create a new one with the desired permissions, name and validity.
**Note:** Watch out to first select the desired organization and project you want the token to be valid for.

Configure the provider by providing your token:

```terraform
provider "metal" {
    api_token = "<YOUR_TOKEN>" # or set env METAL_STACK_CLOUD_API_TOKEN

    # project and organization will be derived from the api_token
}
```

Now you are ready to go!
