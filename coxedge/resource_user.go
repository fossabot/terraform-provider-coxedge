package coxedge

import (
	"context"
	"coxedge/terraform-provider/coxedge/apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: getUserSchema(),
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//Get the API Client
	coxEdgeClient := m.(apiclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	//Convert resource data to API Object
	newUser := convertResourceDataToUserCreateAPIObject(d)

	//Call the API
	createdUser, err := coxEdgeClient.CreateUser(newUser)
	if err != nil {
		return diag.FromErr(err)
	}

	//Save the ID
	d.SetId(createdUser.Id)

	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//Get the API Client
	coxEdgeClient := m.(apiclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	//Get the resource ID
	resourceId := d.Id()

	//Get the resource
	user, err := coxEdgeClient.GetUser(resourceId)
	if err != nil {
		return diag.FromErr(err)
	}

	convertUserAPIObjectToResourceData(d, user)

	//Update state
	resourceUserRead(ctx, d, m)

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//Get the API Client
	coxEdgeClient := m.(apiclient.Client)

	//Get the resource ID
	resourceId := d.Id()

	//Convert resource data to API object
	updatedUser := convertResourceDataToUserCreateAPIObject(d)

	//Call the API
	_, err := coxEdgeClient.UpdateUser(resourceId, updatedUser)
	if err != nil {
		return diag.FromErr(err)
	}

	//Set last_updated
	d.Set("last_updated", time.Now().Format(time.RFC850))

	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	//Get the API Client
	coxEdgeClient := m.(apiclient.Client)

	//Get the resource ID
	resourceId := d.Id()

	//Delete the User
	err := coxEdgeClient.DeleteUser(resourceId)
	if err != nil {
		return diag.FromErr(err)
	}

	// From Docs: d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}

func convertResourceDataToUserCreateAPIObject(d *schema.ResourceData) apiclient.UserCreateRequest {
	//Get roles
	var roleIds []apiclient.IdOnlyHelper
	for _, val := range d.Get("roles").(map[string]interface{}) {
		newRoleId := apiclient.IdOnlyHelper{
			Id: val.(string),
		}
		roleIds = append(roleIds, newRoleId)
	}

	//Create update user struct
	updatedUser := apiclient.UserCreateRequest{
		UserName:  d.Get("userName").(string),
		FirstName: d.Get("firstName").(string),
		LastName:  d.Get("lastName").(string),
		Email:     d.Get("email").(string),
		OrganizationId: apiclient.IdOnlyHelper{
			Id: d.Get("organizationId").(string),
		},
		Roles: roleIds,
	}

	return updatedUser
}

func convertUserAPIObjectToResourceData(d *schema.ResourceData, user *apiclient.User) {
	//Store the data
	d.Set("id", user.Id)
	d.Set("organizationId", user.Organization.Id)
	d.Set("userName", user.UserName)
	d.Set("firstName", user.FirstName)
	d.Set("lastName", user.LastName)

	roles := make([]interface{}, len(user.Roles), len(user.Roles))
	for i, role := range user.Roles {
		item := make(map[string]interface{})
		item["id"] = role.Id
		roles[i] = item
	}
	d.Set("roles", roles)
}
