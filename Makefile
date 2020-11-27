terraform-provider-codeship: *.go */*.go
	go build .

install: terraform-provider-codeship
	mkdir -p ~/.terraform.d/plugins/tails.com/tailsdotcom/codeship/0.1.1/linux_amd64
	cp $+ ~/.terraform.d/plugins/tails.com/tailsdotcom/codeship/0.1.1/linux_amd64

init: install
	terraform init

apply: init
	terraform apply
