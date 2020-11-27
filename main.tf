terraform {
  required_providers {
    codeship = {
      source = "tails.com/tailsdotcom/codeship"
    }
  }
  required_version = ">= 0.13"
}

resource "codeship_project" "test" {
  repo = "https://github.com/tailsdotcom/terraform-provider-codeship"
}
