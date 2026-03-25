package provider

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	providerschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/virtualmo/terraform-provider-netskopesdwan/internal/client"
	"github.com/virtualmo/terraform-provider-netskopesdwan/internal/datasources"
)

var _ provider.Provider = &netskopeSDWANProvider{}

type netskopeSDWANProvider struct{}

type providerModel struct {
	BaseURL  types.String `tfsdk:"base_url"`
	APIToken types.String `tfsdk:"api_token"`
	Timeout  types.Int64  `tfsdk:"timeout"`
	Insecure types.Bool   `tfsdk:"insecure"`
}

func New() provider.Provider {
	return &netskopeSDWANProvider{}
}

func (p *netskopeSDWANProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "netskopesdwan"
}

func (p *netskopeSDWANProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = providerschema.Schema{
		Attributes: map[string]providerschema.Attribute{
			"base_url": providerschema.StringAttribute{
				Optional:    true,
				Description: "Base URL for the Netskope SD-WAN API. Can also be set with NETSKOPESDWAN_BASE_URL.",
			},
			"api_token": providerschema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "API token for the Netskope SD-WAN API. Can also be set with NETSKOPESDWAN_API_TOKEN.",
			},
			"timeout": providerschema.Int64Attribute{
				Optional:    true,
				Description: "Optional HTTP timeout in seconds. Can also be set with NETSKOPESDWAN_TIMEOUT.",
			},
			"insecure": providerschema.BoolAttribute{
				Optional:    true,
				Description: "Optional placeholder for insecure TLS behavior. Can also be set with NETSKOPESDWAN_INSECURE.",
			},
		},
	}
}

func (p *netskopeSDWANProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config providerModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	checkUnknownString(resp, config.BaseURL, path.Root("base_url"))
	checkUnknownString(resp, config.APIToken, path.Root("api_token"))
	checkUnknownInt64(resp, config.Timeout, path.Root("timeout"))
	checkUnknownBool(resp, config.Insecure, path.Root("insecure"))
	if resp.Diagnostics.HasError() {
		return
	}

	baseURL := firstNonEmpty(config.BaseURL.ValueString(), os.Getenv("NETSKOPESDWAN_BASE_URL"))
	apiToken := firstNonEmpty(config.APIToken.ValueString(), os.Getenv("NETSKOPESDWAN_API_TOKEN"))

	insecure := false
	if !config.Insecure.IsNull() {
		insecure = config.Insecure.ValueBool()
	} else if envValue := os.Getenv("NETSKOPESDWAN_INSECURE"); envValue != "" {
		parsed, err := strconv.ParseBool(envValue)
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("insecure"),
				"Invalid insecure value",
				fmt.Sprintf("Expected NETSKOPESDWAN_INSECURE to be a boolean, got %q: %s", envValue, err),
			)
			return
		}
		insecure = parsed
	}

	if baseURL == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("base_url"),
			"Missing base_url",
			"Set base_url in the provider configuration or NETSKOPESDWAN_BASE_URL in the environment.",
		)
	}

	if apiToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Missing api_token",
			"Set api_token in the provider configuration or NETSKOPESDWAN_API_TOKEN in the environment.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	timeoutSeconds := int64(30)
	if !config.Timeout.IsNull() {
		timeoutSeconds = config.Timeout.ValueInt64()
	} else if envValue := os.Getenv("NETSKOPESDWAN_TIMEOUT"); envValue != "" {
		parsedTimeout, err := strconv.ParseInt(envValue, 10, 64)
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("timeout"),
				"Invalid timeout value",
				fmt.Sprintf("Expected NETSKOPESDWAN_TIMEOUT to be an integer number of seconds, got %q: %s", envValue, err),
			)
			return
		}
		timeoutSeconds = parsedTimeout
	}

	if timeoutSeconds <= 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("timeout"),
			"Invalid timeout value",
			"Timeout must be greater than 0 seconds.",
		)
		return
	}

	timeout := time.Duration(timeoutSeconds) * time.Second
	apiClient := client.New(baseURL, apiToken, timeout, insecure)
	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient
}

func (p *netskopeSDWANProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewGatewaysDataSource,
		datasources.NewGatewayDataSource,
	}
}

func (p *netskopeSDWANProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}

	return ""
}

func checkUnknownString(resp *provider.ConfigureResponse, value types.String, attributePath path.Path) {
	if value.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			attributePath,
			"Unknown provider attribute",
			"This provider attribute must be known before configuring the provider.",
		)
	}
}

func checkUnknownBool(resp *provider.ConfigureResponse, value types.Bool, attributePath path.Path) {
	if value.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			attributePath,
			"Unknown provider attribute",
			"This provider attribute must be known before configuring the provider.",
		)
	}
}

func checkUnknownInt64(resp *provider.ConfigureResponse, value types.Int64, attributePath path.Path) {
	if value.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			attributePath,
			"Unknown provider attribute",
			"This provider attribute must be known before configuring the provider.",
		)
	}
}
