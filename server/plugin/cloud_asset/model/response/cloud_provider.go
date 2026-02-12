package response

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/cloud_asset/model"

type ProviderInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ProviderLimits struct {
	MaxSyncBatch int `json:"maxSyncBatch"`
}

type ProviderConfigResponse struct {
	ProviderInfo ProviderInfo     `json:"providerInfo"`
	Regions      []model.RegionVO `json:"regions"`
	Limits       ProviderLimits   `json:"limits"`
}
