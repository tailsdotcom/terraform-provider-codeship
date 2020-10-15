# CodeShip Provider

A provider to setup CI projects on https://codeship.com/.

At the moment it only handles CodeShip Pro (PRs welcome).

## Example Usage

```hcl
resource "codeship_project" "main" {
  repo = github_repository.main.html_url
}
```

## Argument Reference

Can be provided via `CODESHIP_*` environment variables.

* `organization` - (Required) The name of your organization account in CodeShip. 
* `username` - (Required).
* `password` - (Required).
