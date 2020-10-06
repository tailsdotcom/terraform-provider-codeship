package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/codeship/codeship-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"repo": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"aes_key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*codeship.Organization)
	project, _, err := c.CreateProject(ctx, codeship.ProjectCreateRequest{
		RepositoryURL: d.Get("repo").(string),
		Type:          codeship.ProjectTypePro,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(project.UUID)
	d.Set("name", project.Name)
	d.Set("aes_key", project.AesKey)
	return nil
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*codeship.Organization)
	project, _, err := c.GetProject(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(project.UUID)
	err = d.Set("repo", project.RepositoryURL)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("aes_key", project.AesKey)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceProjectRead(ctx, d, m)
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Errorf("Codeship do not provide a delete API")
}
