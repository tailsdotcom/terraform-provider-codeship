# Project Resource

Creates a new [project](https://documentation.codeship.com/general/projects/getting-started/).

## Example Usage

```hcl
resource "codeship_project" "main" {
  repo = github_repository.main.html_url
}
```

## Argument Reference

* `repo` - (Required) The URL of the VCS repository the project is for.

## Attribute Reference

* `aes_key` - The AES key used to encrypt secrets for the project.