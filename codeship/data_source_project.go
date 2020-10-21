package codeship

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/codeship/codeship-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectRead,
		Schema: map[string]*schema.Schema{
			"repo": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"aes_key": &schema.Schema{
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if m == nil {
		return diag.Errorf("Codeship authentication.")
	}
	c := m.(*codeship.Organization)
	projectList, resp, err := c.ListProjects(ctx, codeship.PerPage(50))
	if err != nil {
		return diag.FromErr(err)
	}
	for {
		for _, project := range projectList.Projects {
			if project.RepositoryURL == d.Get("repo").(string) {
				d.SetId(project.UUID)
				d.Set("aes_key", project.AesKey)
				return nil
			}
		}
		if resp.IsLastPage() || resp.Next == "" {
			break
		}
		next, err := resp.NextPage()
		if err != nil {
			return diag.FromErr(err)
		}
		projectList, resp, err = c.ListProjects(ctx, codeship.Page(next), codeship.PerPage(50))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return diag.Errorf("Project not found: %s", d.Get("repo").(string))
}
