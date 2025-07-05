package http

import (
	"crypto/tls"
	"testing"

	"github.com/crossplane-contrib/provider-http/apis/common"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildTLSConfig(t *testing.T) {
	logger := logging.NewNopLogger()
	client := &client{
		log: logger,
	}

	tests := []struct {
		name        string
		tlsConfig   *common.TLSConfig
		expectError bool
		validate    func(t *testing.T, config *tls.Config)
	}{
		{
			name:        "nil config returns empty config",
			tlsConfig:   nil,
			expectError: false,
			validate: func(t *testing.T, config *tls.Config) {
				assert.False(t, config.InsecureSkipVerify)
				assert.Nil(t, config.RootCAs)
				assert.Empty(t, config.Certificates)
			},
		},
		{
			name: "insecure skip verify",
			tlsConfig: &common.TLSConfig{
				InsecureSkipVerify: true,
			},
			expectError: false,
			validate: func(t *testing.T, config *tls.Config) {
				assert.True(t, config.InsecureSkipVerify)
			},
		},
		{
			name: "invalid CA certificate",
			tlsConfig: &common.TLSConfig{
				CAData: "invalid-cert-data",
			},
			expectError: true,
		},
		{
			name: "client cert without key should fail",
			tlsConfig: &common.TLSConfig{
				ClientCertData: "some-cert",
			},
			expectError: true,
		},
		{
			name: "client key without cert should fail",
			tlsConfig: &common.TLSConfig{
				ClientKeyData: "some-key",
			},
			expectError: true,
		},
		{
			name: "invalid client certificate",
			tlsConfig: &common.TLSConfig{
				ClientCertData: "invalid-cert",
				ClientKeyData:  "invalid-key",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := client.buildTLSConfig(tt.tlsConfig)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, config)

			if tt.validate != nil {
				tt.validate(t, config)
			}
		})
	}
}