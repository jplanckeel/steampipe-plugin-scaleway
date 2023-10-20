package scaleway

import (
	"context"

	"github.com/scaleway/scaleway-sdk-go/api/account/v2"
	cockpit "github.com/scaleway/scaleway-sdk-go/api/cockpit/v1beta1"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableScalewayCockpitTokens(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "scaleway_cockpit_tokens",
		Description: "List of cockpit token in Scaleway.",
		List: &plugin.ListConfig{
			ParentHydrate: listProjects,
			Hydrate:       listCockpitPlans,
			KeyColumns:    []*plugin.KeyColumn{},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "Name of the token.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "token_id",
				Description: "A unique identifier of the cockpit token.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "created_at",
				Description: "The date and time of the token's creation.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_at",
				Description: "The date and time of the token's last update.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "scopes",
				Description: "Token permissions.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "project",
				Description: "The ID of the Project.",
				Type:        proto.ColumnType_STRING,
			},
		},
	}
}

//// LIST FUNCTION

func listCockpitTokens(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_cockpit_contact_points.listCockpitTokens", "connection_error", err)
		return nil, err
	}

	// Get Project details
	projectData := h.Item.(*account.Project)

	// Create SDK objects for Scaleway Cockpit product
	cockpitApi := cockpit.NewAPI(client)

	req := &cockpit.ListTokensRequest{
		Page:      scw.Int32Ptr(1),
		ProjectID: projectData.ID,
	}

	// Retrieve the list of contact points
	maxResult := int64(100)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < maxResult {
			maxResult = *limit
		}
	}
	req.PageSize = scw.Uint32Ptr(uint32(maxResult))

	var count int
	for {
		resp, err := cockpitApi.ListTokens(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_cockpit_contact_points.listCockpitTokens", "query_error", err)
			//Break if cockpit does not exist in project
			break
		}

		for _, token := range resp.Tokens {
			d.StreamListItem(ctx, token)

			// Increase the resource count by 1
			count++

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if resp.TotalCount == uint32(count) {
			break
		}
		req.Page = scw.Int32Ptr(*req.Page + 1)
	}

	return nil, nil
}
