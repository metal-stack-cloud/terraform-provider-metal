# Terraform Provider for metalstack.cloud

Manage the lifecycle of your bare-metal Kubernetes clusters on [metalstack.cloud](https://metalstack.cloud) using Terraform.

> **Note:** this project is in an early development stage. It might break in the future and getting started might be cumbersome.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Using the provider

First, let's add the provider to your project:

```terraform
terraform {
  required_providers {
    metal = {
      source = "metalstack.cloud/terraform/metal"
    }
  }
}
```

> **Note:** As [metalstack.cloud](https://metalstack.cloud) does not yet provide API Tokens, you currently need to pick your JWT from your web session. This is obviously going to change. After you logged in, open the Developer Tools of your browser, head to the console, filter for `token` and copy the JWT starting with `eyJ`.
> To get the project id, with dev tools open, switch to your project and open the clusters view. In the dev tools' network tab, search for `api.v1.ClusterService/List`, select a request with status 200, head to the payload and copy the project id. This is obviously going to change.

Now you need to actually configure it before using it:

```terraform
provider "metal" {
    organization = "your-organization@github"
    project = "f8e67080-ba68-41d2-ad44-59dc65a09a33" # replace this with your uuid.
}
```

Now you are ready to go!
