terraform {
  required_providers {
    netskopesdwan = {
      source = "virtualmo/netskopesdwan"
    }
  }
}

provider "netskopesdwan" {
  base_url  = "https://example.netskopesdwan.local"
  api_token = "replace-me"
  timeout   = "30s"
  insecure  = false
}

# TODO: Add datasource examples once the API fields and lookup arguments are confirmed.
