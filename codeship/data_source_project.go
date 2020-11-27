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
			"ssh_key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"team_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"notification_rule": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"notifier": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"target": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"branch": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"key": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"room": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"build_statuses": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"branch_match": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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
				return readProject(d, project)
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

func readProject(d *schema.ResourceData, project codeship.Project) diag.Diagnostics {
	d.SetId(project.UUID)
	err := d.Set("repo", project.RepositoryURL)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("aes_key", project.AesKey)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("ssh_key", project.SSHKey)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("team_ids", project.TeamIDs)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("notification_rule", flattenNotificationRules(project.NotificationRules))
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func flattenNotificationRules(list []codeship.NotificationRule) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"notifier":       i.Notifier,
			"target":         i.Target,
			"branch":         i.Branch,
			"build_statuses": i.BuildStatuses,
			"key":            i.Options.Key,
			"url":            i.Options.URL,
			"room":           i.Options.Room,
			"branch_match":   i.BranchMatch,
		}
		result = append(result, l)
	}
	return result
}
