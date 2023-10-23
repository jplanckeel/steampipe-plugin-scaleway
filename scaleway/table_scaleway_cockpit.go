package scaleway

import (
	"context"

	"github.com/scaleway/scaleway-sdk-go/api/account/v2"
	cockpit "github.com/scaleway/scaleway-sdk-go/api/cockpit/v1beta1"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableScalewayCockpit(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "scaleway_cockpit",
		Description: "List of cockpit in Scaleway.",
		List: &plugin.ListConfig{
			ParentHydrate: listProjects,
			Hydrate:       getCockpits,
			KeyColumns:    []*plugin.KeyColumn{},
		},
		Columns: []*plugin.Column{
			{
				Name:        "project_id",
				Description: "A unique identifier of the cockpit token.",
				Type:        proto.ColumnType_STRING,
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
				Name:        "endpoints",
				Description: "Token permissions.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "status",
				Description: "The ID of the Project.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "managed_alerts_enabled",
				Description: "The ID of the Project.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "plan",
				Description: "The ID of the Project.",
				Type:        proto.ColumnType_JSON,
			},
		},
	}
}

//// LIST FUNCTION

func getCockpits(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

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

	req := &cockpit.GetCockpitRequest{
		ProjectID: projectData.ID,
	}

	resp, err := cockpitApi.GetCockpit(req)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_cockpit_contact_points.listCockpitTokens", "query_error", err)
		return nil, nil
	}

	return resp, nil
}
