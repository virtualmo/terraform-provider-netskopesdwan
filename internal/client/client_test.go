package client

import "testing"

func TestDecodeGatewayListResponse(t *testing.T) {
	t.Parallel()

	body := []byte(`{
		"page_info": {
			"end_cursor": "cursor-2",
			"has_next": true,
			"total_count": 2
		},
		"data": [
			{
				"id": "gw-1",
				"name": "Gateway One",
				"created_at": "2026-03-25T10:00:00Z",
				"modified_at": "2026-03-25T11:00:00Z",
				"config_updates_enabled": true,
				"managed": false
			}
		]
	}`)

	got, err := decodeGatewayListResponse(body)
	if err != nil {
		t.Fatalf("decodeGatewayListResponse() error = %v", err)
	}

	if got.PageInfo.EndCursor != "cursor-2" {
		t.Fatalf("EndCursor = %q, want %q", got.PageInfo.EndCursor, "cursor-2")
	}
	if !got.PageInfo.HasNext {
		t.Fatal("HasNext = false, want true")
	}
	if got.PageInfo.TotalCount != 2 {
		t.Fatalf("TotalCount = %d, want %d", got.PageInfo.TotalCount, 2)
	}
	if len(got.Data) != 1 {
		t.Fatalf("len(Data) = %d, want %d", len(got.Data), 1)
	}

	item := got.Data[0]
	if item.ID != "gw-1" {
		t.Fatalf("ID = %q, want %q", item.ID, "gw-1")
	}
	if item.Name != "Gateway One" {
		t.Fatalf("Name = %q, want %q", item.Name, "Gateway One")
	}
	if item.CreatedAt != "2026-03-25T10:00:00Z" {
		t.Fatalf("CreatedAt = %q, want %q", item.CreatedAt, "2026-03-25T10:00:00Z")
	}
	if item.ModifiedAt != "2026-03-25T11:00:00Z" {
		t.Fatalf("ModifiedAt = %q, want %q", item.ModifiedAt, "2026-03-25T11:00:00Z")
	}
	if !item.ConfigUpdatesEnabled {
		t.Fatal("ConfigUpdatesEnabled = false, want true")
	}
	if item.Managed {
		t.Fatal("Managed = true, want false")
	}
}

func TestDecodeGatewayListResponseEmptyDataDefaultsToEmptySlice(t *testing.T) {
	t.Parallel()

	body := []byte(`{
		"page_info": {
			"end_cursor": "",
			"has_next": false,
			"total_count": 0
		}
	}`)

	got, err := decodeGatewayListResponse(body)
	if err != nil {
		t.Fatalf("decodeGatewayListResponse() error = %v", err)
	}

	if got.Data == nil {
		t.Fatal("Data = nil, want empty slice")
	}
	if len(got.Data) != 0 {
		t.Fatalf("len(Data) = %d, want 0", len(got.Data))
	}
}

func TestDecodeGatewayResponseBareObject(t *testing.T) {
	t.Parallel()

	body := []byte(`{
		"id": "gw-9",
		"name": "Gateway Nine",
		"created_at": "2026-03-24T09:00:00Z",
		"modified_at": "2026-03-25T09:30:00Z",
		"config_updates_enabled": false,
		"managed": true
	}`)

	got, err := decodeGatewayResponse(body)
	if err != nil {
		t.Fatalf("decodeGatewayResponse() error = %v", err)
	}

	if got.ID != "gw-9" {
		t.Fatalf("ID = %q, want %q", got.ID, "gw-9")
	}
	if got.Name != "Gateway Nine" {
		t.Fatalf("Name = %q, want %q", got.Name, "Gateway Nine")
	}
	if got.CreatedAt != "2026-03-24T09:00:00Z" {
		t.Fatalf("CreatedAt = %q, want %q", got.CreatedAt, "2026-03-24T09:00:00Z")
	}
	if got.ModifiedAt != "2026-03-25T09:30:00Z" {
		t.Fatalf("ModifiedAt = %q, want %q", got.ModifiedAt, "2026-03-25T09:30:00Z")
	}
	if got.ConfigUpdatesEnabled {
		t.Fatal("ConfigUpdatesEnabled = true, want false")
	}
	if !got.Managed {
		t.Fatal("Managed = false, want true")
	}
}
