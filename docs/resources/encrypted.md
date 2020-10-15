# Encrypted Resource

Provides a method of [encrypting](https://documentation.codeship.com/pro/builds-and-configuration/environment-variables/#encrypted-environment-variables) data with a project's AES key.



## Example Usage

```hcl
resource "codeship_encrypted" "aws" {
  aes_key  = codeship_project.main.aes_key
  content  = <<EOF
AWS_ACCESS_KEY_ID=${aws_iam_access_key.deployment.id}
AWS_SECRET_ACCESS_KEY=${aws_iam_access_key.deployment.secret}
AWS_REGION=${data.aws_region.current.name}
EOF
}

resource "github_repository_file" "aws" {
  repository = codeship_project.main.repo
  file       = "aws-deployment.env.encrypted"
  content    = codeship_encrypted.aws.encrypted_content
}
```

## Argument Reference

* `aes_key` - (Required) The secret to encrypt with.
* `content` - (Required) The plain-text content to encrypt.

## Attribute Reference

* `encrypted_content` - The resulting encrypted content.