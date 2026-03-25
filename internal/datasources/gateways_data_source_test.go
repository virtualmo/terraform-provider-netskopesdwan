package datasources

import (
	"context"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestGatewaysDataSourceSchemaIsConservative(t *testing.T) {
	t.Parallel()

	ds := NewGatewaysDataSource()
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

	gotAttributes := make([]string, 0, len(resp.Schema.Attributes))
	for name := range resp.Schema.Attributes {
		gotAttributes = append(gotAttributes, name)
	}
	slices.Sort(gotAttributes)

	wantAttributes := []string{"id", "items"}
	if !slices.Equal(gotAttributes, wantAttributes) {
		t.Fatalf("attributes = %v, want %v", gotAttributes, wantAttributes)
	}

	if gatewaysCollectionID != "gateways" {
		t.Fatalf("gatewaysCollectionID = %q, want %q", gatewaysCollectionID, "gateways")
	}
}
