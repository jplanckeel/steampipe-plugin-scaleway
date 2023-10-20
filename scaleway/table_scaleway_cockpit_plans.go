package scaleway

import (
	"context"

	cockpit "github.com/scaleway/scaleway-sdk-go/api/cockpit/v1beta1"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableScalewayCockpitPlans(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "scaleway_cockpit_plans",
		Description: "List of all pricing plans available for Cockpit in Scaleway.",
		List: &plugin.ListConfig{
			Hydrate:    listCockpitPlans,
			KeyColumns: []*plugin.KeyColumn{},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "Name of a given pricing plan.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "plan_id",
				Description: "A unique identifier of a given pricing plan.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "retention_metrics_interval_days",
				Description: "Retention for metrics in days.",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("RetentionMetricsIntervalDays"),
			},
			{
				Name:        "retention_logs_interval_days",
				Description: "Retention for logs in days.",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("RetentionLogsIntervalDays"),
			},
			{
				Name:        "sample_ingestion_price",
				Description: "Ingestion price for 1 million samples in cents.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "logs_ingestion_price",
				Description: "Ingestion price for 1 GB of logs in cents.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "retention_price",
				Description: "Retention price in euros per month.",
				Type:        proto.ColumnType_INT,
			},
		},
	}
}

type cockpitPlanInfo = struct {
	cockpit.Plan
	RetentionMetricsIntervalDays float64
	RetentionLogsIntervalDays    float64
}

//// LIST FUNCTION

func listCockpitPlans(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create client
	client, err := getSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("scaleway_cockpit_plans.listCockpitPlans", "connection_error", err)
		return nil, err
	}

	// Create SDK objects for Scaleway Cockpit product
	cockpitApi := cockpit.NewAPI(client)

	req := &cockpit.ListPlansRequest{
		Page: scw.Int32Ptr(1),
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
		resp, err := cockpitApi.ListPlans(req)
		if err != nil {
			plugin.Logger(ctx).Error("scaleway_cockpit_plans.listCockpitPlans", "query_error", err)
			//Break if cockpit does not exist in project
			break
		}

		for _, plan := range resp.Plans {
			d.StreamListItem(ctx, &cockpitPlanInfo{
				*plan,
				(plan.RetentionMetricsInterval.ToTimeDuration().Hours() / 24),
				(plan.RetentionLogsInterval.ToTimeDuration().Hours() / 24),
			})

			// Increase the resource count by 1
			count++

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if resp.TotalCount == uint64(count) {
			break
		}
		req.Page = scw.Int32Ptr(*req.Page + 1)
	}

	return nil, nil
}
