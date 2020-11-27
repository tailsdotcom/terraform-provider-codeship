package codeship

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
				ForceNew: true,
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
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"notification_rule": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"notifier": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"target": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "all",
						},
						"branch": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"key": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"url": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"room": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"build_statuses": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"branch_match": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "exact",
						},
					},
				},
			},
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if m == nil {
		return diag.Errorf("Codeship authentication.")
	}
	c := m.(*codeship.Organization)
	project, _, err := c.CreateProject(ctx, codeship.ProjectCreateRequest{
		RepositoryURL: d.Get("repo").(string),
		Type:          codeship.ProjectTypePro,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return readProject(d, project)
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if m == nil {
		return diag.Errorf("Codeship authentication.")
	}
	c := m.(*codeship.Organization)
	project, _, err := c.GetProject(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return readProject(d, project)
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if m == nil {
		return diag.Errorf("Codeship authentication.")
	}
	c := m.(*codeship.Organization)
	project, _, err := c.UpdateProject(ctx, d.Id(), codeship.ProjectUpdateRequest{
		TeamIDs:           expandTeamIDs(d.Get("team_ids").(*schema.Set)),
		NotificationRules: expandNotificationRules(d.Get("notification_rule").(*schema.Set)),
		Type:              codeship.ProjectTypePro,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return readProject(d, project)
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Errorf("Codeship do not provide a delete API")
}

func expandTeamIDs(set *schema.Set) []int {
	result := make([]int, 0, set.Len())
	for _, i := range set.List() {
		result = append(result, i.(int))
	}
	return result
}

func expandNotificationRules(set *schema.Set) []codeship.NotificationRule {
	result := make([]codeship.NotificationRule, 0, set.Len())
	for _, i := range set.List() {
		ii := i.(map[string]interface{})
		l := codeship.NotificationRule{
			Notifier:      ii["notifier"].(string),
			Target:        ii["target"].(string),
			Branch:        ii["branch"].(string),
			BuildStatuses: expandBuildStatuses(ii["build_statuses"].([]interface{})),
			Options: codeship.NotificationOptions{
				Key:  ii["key"].(string),
				URL:  ii["url"].(string),
				Room: ii["room"].(string),
			},
			BranchMatch: ii["branch_match"].(string),
		}
		result = append(result, l)
	}
	return result
}

func expandBuildStatuses(list []interface{}) []string {
	result := make([]string, 0, len(list))
	for _, i := range list {
		result = append(result, i.(string))
	}
	return result
}
