package codeship

import (
	"context"

	"github.com/codeship/codeship-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"organization": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CODESHIP_ORGANIZATION", nil),
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CODESHIP_USERNAME", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CODESHIP_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"codeship_project":   resourceProject(),
			"codeship_encrypted": resourceEncrypted(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"codeship_project": dataSourceProject(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	organization := d.Get("organization").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	auth := codeship.NewBasicAuth(username, password)
	client, err := codeship.New(auth)
	if err != nil {
		return nil, []diag.Diagnostic{{
			Severity: diag.Warning,
			Summary:  "Codeship authentication.",
			Detail:   err.Error(),
		}}
	}
	org, err := client.Organization(ctx, organization)
	if err != nil {
		return nil, []diag.Diagnostic{{
			Severity: diag.Warning,
			Summary:  "Codeship can not find organization.",
			Detail:   err.Error(),
		}}
	}
	return org, nil
}
