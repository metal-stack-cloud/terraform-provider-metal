provider "metal" {
  # All arguments are optional and can be omitted
  # The defaults are derived from the environment variables METAL_STACK_CLOUD_* or ~/.metal-stack-cloud/config.yaml
  api_url      = "https://api.metalstack.cloud" # default
  api_token    = "..."
  organization = "yourcompany@github"
  project      = "b94449a2-e105-42ca-9cff-9f7b11fec318" # your project id
}
