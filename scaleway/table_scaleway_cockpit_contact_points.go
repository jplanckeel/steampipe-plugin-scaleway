package scaleway

import (
	"context"

	"github.com/scaleway/scaleway-sdk-go/api/cockpit/v1beta1"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableScalewayCockpitContactPoints(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:              "scaleway_cockpit_contact_points",
		Description:       "A Compute Instance bare metal is a physical server in Scaleway.",
		GetMatrixItemFunc: BuildZoneList,
		List: &plugin.ListConfig{
			Hydrate:    listCockpitContactPoints,
			KeyColumns: []*plugin.KeyColumn{},
		},
		Columns: []*plugin.Column{
			{
				Name:        "totalcount",
				Description: "Count of all contact points created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "contact_points",
				Description: "Array of contact points",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ContactPoints"),
			},
			{
				Name:        "has_additional_receivers",
				Description: "Specifies whether the contact point has other receivers than the default receiver.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "has_additional_contact_points",
				Description: "Specifies whether there are unmanaged contact points.",
				Type:        proto.ColumnType_BOOL,
			},

			// Scaleway standard columns
			{
				Name:        "project",
				Description: "The ID of the project where the server resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "organization",
				Description: "The ID of the organization where the server resides.",
				Type:        proto.ColumnType_STRING,
			},
		},
	}
}

//// LIST FUNCTION

func listCockpitContactPoints(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_cockpit_contact_points.listCockpitContactPoints", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Baremetal product
	cockpitApi := cockpit.NewAPI(client)

	req := &cockpit.ListContactPointsRequest{
		Page:      scw.Int32Ptr(1),
		ProjectID: "",
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
			return nil, err
		}

		for _, baremetal := range resp.ContactPoints {
			d.StreamListItem(ctx, baremetal)

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
