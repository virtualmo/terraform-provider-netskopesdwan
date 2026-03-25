package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/virtualmo/terraform-provider-netskopesdwan/internal/client"
)

var (
	_ datasource.DataSource              = &gatewaysDataSource{}
	_ datasource.DataSourceWithConfigure = &gatewaysDataSource{}
)

const gatewaysCollectionID = "gateways"

type gatewaysDataSource struct {
	client *client.Client
}

type gatewaysDataSourceModel struct {
	ID    types.String        `tfsdk:"id"`
	Items []gatewayModelValue `tfsdk:"items"`
}

type gatewayModelValue struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	CreatedAt            types.String `tfsdk:"created_at"`
	ModifiedAt           types.String `tfsdk:"modified_at"`
	ConfigUpdatesEnabled types.Bool   `tfsdk:"config_updates_enabled"`
	Managed              types.Bool   `tfsdk:"managed"`
}

func NewGatewaysDataSource() datasource.DataSource {
	return &gatewaysDataSource{}
}

func (d *gatewaysDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gateways"
}

func (d *gatewaysDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		Description: "Reads the Netskope SD-WAN gateways collection.",
		Attributes: map[string]datasourceschema.Attribute{
			"id": datasourceschema.StringAttribute{
				Computed:    true,
				Description: "Synthetic identifier for the gateways collection.",
			},
			"items": datasourceschema.ListNestedAttribute{
				Computed:    true,
				Description: "Gateways returned by GET /v2/gateways.",
				NestedObject: datasourceschema.NestedAttributeObject{
					Attributes: map[string]datasourceschema.Attribute{
						"id": datasourceschema.StringAttribute{
							Computed:    true,
							Description: "Gateway identifier.",
						},
						"name": datasourceschema.StringAttribute{
							Computed:    true,
							Description: "Gateway name.",
						},
						"created_at": datasourceschema.StringAttribute{
							Computed:    true,
							Description: "Gateway creation timestamp returned by the API.",
						},
						"modified_at": datasourceschema.StringAttribute{
							Computed:    true,
							Description: "Gateway modification timestamp returned by the API.",
						},
						"config_updates_enabled": datasourceschema.BoolAttribute{
							Computed:    true,
							Description: "Whether config updates are enabled for the gateway.",
						},
						"managed": datasourceschema.BoolAttribute{
							Computed:    true,
							Description: "Whether the gateway is managed.",
						},
					},
				},
			},
		},
	}
}

func (d *gatewaysDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	apiClient, ok := req.ProviderData.(*client.Client)
	if !ok {
		return
	}

	d.client = apiClient
}

func (d *gatewaysDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured provider client",
			"The provider client was not configured before reading netskopesdwan_gateways.",
		)
		return
	}

	result, err := d.client.ListGateways(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read GET /v2/gateways",
			"The Netskope SD-WAN gateways list request failed: "+err.Error(),
		)
		return
	}

	items := make([]gatewayModelValue, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, gatewayModelValue{
			ID:                   types.StringValue(item.ID),
			Name:                 types.StringValue(item.Name),
			CreatedAt:            types.StringValue(item.CreatedAt),
			ModifiedAt:           types.StringValue(item.ModifiedAt),
			ConfigUpdatesEnabled: types.BoolValue(item.ConfigUpdatesEnabled),
			Managed:              types.BoolValue(item.Managed),
		})
	}

	// TODO: Add pagination-aware Terraform behavior if users need access to
	// cursors or multi-page reads beyond the current single-call list behavior.
	state := gatewaysDataSourceModel{
		ID:    types.StringValue(gatewaysCollectionID),
		Items: items,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
