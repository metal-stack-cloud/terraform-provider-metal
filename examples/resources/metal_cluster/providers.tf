terraform {
  required_providers {
    metal = {
      source = "metalstack.cloud/terraform/metal"
    }
  }
}

provider "metal" {
  # All arguments are optional and can be omitted
  # The defaults are derived from the environment variables METAL_STACK_CLOUD_* or ~/.metal-stack-cloud/config.yaml
  api_url      = "https://api.metalstack.cloud" # default
  organization = "x-cellent@github"
  project      = "f8e67080-ba68-41d2-ad44-59dc65a09a33" # your project id
}
