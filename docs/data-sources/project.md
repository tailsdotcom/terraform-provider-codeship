# Project Data Source

Access an existing [project](https://documentation.codeship.com/general/projects/getting-started/).

## Example Usage

```hcl
data "codeship_project" "main" {
  repo = github_repository.main.html_url
}
```

## Argument Reference

* `repo` - (Required) The URL of the VCS repository the project is for.

## Attribute Reference

* `aes_key` - The AES key used to encrypt secrets for the project.
* `ssh_key` - The public SSH key CodeShip will use when deploying.
* `team_ids` - The Teams who have access to this project. As per the API these are opaque identifiers; AFAIK you have to import a resource to discover them.
* `notification_rule`:
  * `build_statuses` - List of statuses to notify for, from: `started`, `failed`, `success`, and `recovered`.
  * `branch` - The repository branch to notify for.
  * `branch_match` - How to `branch` is matched.
  * `target` - Whom to notify: `all`, `committer`.
  * `notifier` - The service to notify: `email`, `slack`, `hipchat`, `campfire`, `grove`, `flowdock`, or `webhook`.
  * `key` - A service defined value.
  * `url` - A service defined value.
  * `room` - A service defined value.
