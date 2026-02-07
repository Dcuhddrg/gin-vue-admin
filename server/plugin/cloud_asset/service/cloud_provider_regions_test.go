package service

import (
	"testing"
)

func TestCloudProviderService_GetRegions(t *testing.T) {
	s := &CloudProviderService{}

	tests := []struct {
		name      string
		provider  string
		accessKey string
		secretKey string
		wantErr   bool
	}{
		{
			name:      "Unsupported Provider",
			provider:  "unknown",
			accessKey: "test",
			secretKey: "test",
			wantErr:   true,
		},
		{
			name:      "AWS Not Supported",
			provider:  "aws",
			accessKey: "test",
			secretKey: "test",
			wantErr:   true,
		},
		{
			name:      "Aliyun Invalid Key",
			provider:  "aliyun",
			accessKey: "invalid",
			secretKey: "invalid",
			wantErr:   true,
		},
		{
			name:      "Tencent Invalid Key",
			provider:  "tencent",
			accessKey: "invalid",
			secretKey: "invalid",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.GetRegions(tt.provider, tt.accessKey, tt.secretKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("CloudProviderService.GetRegions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
