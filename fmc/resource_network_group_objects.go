package fmc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var network_group_type string = "NetworkGroup"

func resourceNetworkGroupObjects() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkGroupObjectsCreate,
		ReadContext:   resourceNetworkGroupObjectsRead,
		UpdateContext: resourceNetworkGroupObjectsUpdate,
		DeleteContext: resourceNetworkGroupObjectsDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"objects": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"literals": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceNetworkGroupObjectsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	// Warning or errors can be collected in a slice type
	// var diags diag.Diagnostics
	var diags diag.Diagnostics

	var objs []NetworkGroupObjectObjects
	var lits []NetworkGroupObjectLiterals

	if inputObjs, ok := d.GetOk("objects"); ok {
		for _, obj := range inputObjs.([]interface{}) {
			obji := obj.(map[string]interface{})
			objs = append(objs, NetworkGroupObjectObjects{
				ID:   obji["id"].(string),
				Name: obji["name"].(string),
				Type: obji["type"].(string),
			})
		}
	}

	if inputLits, ok := d.GetOk("literals"); ok {
		for _, lit := range inputLits.([]interface{}) {
			liti := lit.(map[string]interface{})
			lits = append(lits, NetworkGroupObjectLiterals{
				Value: liti["value"].(string),
				Type:  liti["type"].(string),
			})
		}
	}

	res, err := c.CreateNetworkGroupObject(ctx, &NetworkGroupObject{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Type:        network_group_type,
		Objects:     objs,
		Literals:    lits,
	})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to create network group object",
			Detail:   err.Error(),
		})
		return diags
	}
	d.SetId(res.ID)
	return resourceNetworkGroupObjectsRead(ctx, d, m)
}

func resourceNetworkGroupObjectsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	item, err := c.GetNetworkGroupObject(ctx, id)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to read network group object",
			Detail:   err.Error(),
		})
		return diags
	}
	if err := d.Set("name", item.Name); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to read network group object",
			Detail:   err.Error(),
		})
		return diags
	}

	if err := d.Set("description", item.Description); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to read network group object",
			Detail:   err.Error(),
		})
		return diags
	}

	if err := d.Set("type", item.Type); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to read network group object",
			Detail:   err.Error(),
		})
		return diags
	}

	objects := make([]interface{}, len(item.Objects))
	for _, obj := range item.Objects {
		obji := make(map[string]interface{})
		obji["id"] = obj.ID
		obji["name"] = obj.Name
		obji["type"] = obj.Type
		objects = append(objects, obji)
	}

	if err := d.Set("objects", objects); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to read network group object",
			Detail:   err.Error(),
		})
		return diags
	}

	literals := make([]interface{}, len(item.Literals))
	for _, lit := range item.Literals {
		liti := make(map[string]interface{})
		liti["value"] = lit.Value
		liti["type"] = lit.Type
		literals = append(literals, liti)
	}

	if err := d.Set("literals", literals); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to read network group object",
			Detail:   err.Error(),
		})
		return diags
	}

	return diags
}

func resourceNetworkGroupObjectsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	var diags diag.Diagnostics
	id := d.Id()
	if d.HasChanges("name", "description", "objects", "literals") {
		var objs []NetworkGroupObjectObjects
		var lits []NetworkGroupObjectLiterals

		if inputObjs, ok := d.GetOk("objects"); ok {
			for _, obj := range inputObjs.([]interface{}) {
				obji := obj.(map[string]interface{})
				objs = append(objs, NetworkGroupObjectObjects{
					ID:   obji["id"].(string),
					Name: obji["name"].(string),
					Type: obji["type"].(string),
				})
			}
		}

		if inputLits, ok := d.GetOk("literals"); ok {
			for _, lit := range inputLits.([]interface{}) {
				liti := lit.(map[string]interface{})
				lits = append(lits, NetworkGroupObjectLiterals{
					Value: liti["value"].(string),
					Type:  liti["type"].(string),
				})
			}
		}
		_, err := c.UpdateNetworkGroupObject(ctx, id, &NetworkGroupObjectUpdateInput{
			ID:          d.Id(),
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Type:        network_group_type,
			Objects:     objs,
			Literals:    lits,
		})

		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "unable to update network group object",
				Detail:   err.Error(),
			})
			return diags
		}
	}
	return resourceNetworkGroupObjectsRead(ctx, d, m)
}

func resourceNetworkGroupObjectsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()

	err := c.DeleteNetworkGroupObject(ctx, id)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to delete network group object",
			Detail:   err.Error(),
		})
		return diags
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}