package scaleway

import (
	"context"

	"github.com/scaleway/scaleway-sdk-go/api/account/v2"
	cockpit "github.com/scaleway/scaleway-sdk-go/api/cockpit/v1beta1"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableScalewayCockpitContactPoints(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "scaleway_cockpit_contact_points",
		Description: "Contact points define who is notified when an alert fires in Scaleway.",
		List: &plugin.ListConfig{
			ParentHydrate: listProjects,
			Hydrate:       listCockpitContactPoints,
			KeyColumns:    []*plugin.KeyColumn{},
		},
		Columns: []*plugin.Column{
			{
				Name:        "email",
				Description: "Contact points created.",
				Type:        proto.ColumnType_STRING,
			},

			// Scaleway standard columns
			{
				Name:        "project",
				Description: "The ID of the project where the server resides.",
				Type:        proto.ColumnType_STRING,
			},
		},
	}
}

type cockpitContactInfo = struct {
	Email   string
	Project string
}

//// LIST FUNCTION

func listCockpitContactPoints(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_cockpit_contact_points.listCockpitContactPoints", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Cockpit product
	cockpitApi := cockpit.NewAPI(client)

	// Get Project details
	projectData := h.Item.(*account.Project)

	req := &cockpit.ListContactPointsRequest{
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
		resp, err := cockpitApi.ListContactPoints(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_cockpit_contact_points.listCockpitContactPoints", "query_error", err)
			//Break if cockpit does not exist in project
			break
		}

		for _, contactPoint := range resp.ContactPoints {
			d.StreamListItem(ctx, cockpitContactInfo{contactPoint.Email.To, projectData.ID})

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
