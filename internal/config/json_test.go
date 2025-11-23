package config

import (
	"github.com/spitfy/urlshortener/internal/config/db"
	handlerConf "github.com/spitfy/urlshortener/internal/handler/config"
	storageConf "github.com/spitfy/urlshortener/internal/repository/config"
	serviceConf "github.com/spitfy/urlshortener/internal/service/config"
	"testing"
)

func TestApplyJSONConfig(t *testing.T) {
	tests := []struct {
		name           string
		initialConfig  Config
		jsonConfig     JSONConfig
		expectedConfig Config
	}{
		{
			name: "apply all fields when default",
			initialConfig: Config{
				Handlers: handlerConf.Config{
					ServerAddr:  ":8080",
					EnableHTTPS: false,
				},
				Service: serviceConf.Config{
					ServerURL: "http://localhost:8080",
				},
				FileStorage: storageConf.Config{
					FileStoragePath: "",
				},
				DB: db.Config{
					DatabaseDsn: "",
				},
			},
			jsonConfig: JSONConfig{
				ServerAddress:   "localhost:8081",
				BaseURL:         "http://myapp.com",
				FileStoragePath: "/tmp/storage.json",
				DatabaseDSN:     "postgres://...",
				EnableHTTPS:     true,
			},
			expectedConfig: Config{
				Handlers: handlerConf.Config{
					ServerAddr:  "localhost:8081",
					EnableHTTPS: true,
				},
				Service: serviceConf.Config{
					ServerURL: "http://myapp.com",
				},
				FileStorage: storageConf.Config{
					FileStoragePath: "/tmp/storage.json",
				},
				DB: db.Config{
					DatabaseDsn: "postgres://...",
				},
			},
		},
		{
			name: "do not override if already set",
			initialConfig: Config{
				Handlers: handlerConf.Config{
					ServerAddr:  ":9090",
					EnableHTTPS: true,
				},
				Service: serviceConf.Config{
					ServerURL: "http://custom.com",
				},
				FileStorage: storageConf.Config{
					FileStoragePath: "/custom/path",
				},
				DB: db.Config{
					DatabaseDsn: "sqlite://...",
				},
			},
			jsonConfig: JSONConfig{
				ServerAddress:   "localhost:8081",
				BaseURL:         "http://myapp.com",
				FileStoragePath: "/tmp/storage.json",
				DatabaseDSN:     "postgres://...",
				EnableHTTPS:     false,
			},
			expectedConfig: Config{
				Handlers: handlerConf.Config{
					ServerAddr:  ":9090",
					EnableHTTPS: true,
				},
				Service: serviceConf.Config{
					ServerURL: "http://custom.com",
				},
				FileStorage: storageConf.Config{
					FileStoragePath: "/custom/path",
				},
				DB: db.Config{
					DatabaseDsn: "sqlite://...",
				},
			},
		},
		{
			name: "do not apply empty string values",
			initialConfig: Config{
				Handlers: handlerConf.Config{
					ServerAddr:  ":8080",
					EnableHTTPS: false,
				},
				Service: serviceConf.Config{
					ServerURL: "http://localhost:8080",
				},
				FileStorage: storageConf.Config{
					FileStoragePath: "",
				},
				DB: db.Config{
					DatabaseDsn: "",
				},
			},
			jsonConfig: JSONConfig{
				ServerAddress:   "",
				BaseURL:         "",
				FileStoragePath: "",
				DatabaseDSN:     "",
				EnableHTTPS:     true,
			},
			expectedConfig: Config{
				Handlers: handlerConf.Config{
					ServerAddr:  ":8080",
					EnableHTTPS: true, // bool applies even if default
				},
				Service: serviceConf.Config{
					ServerURL: "http://localhost:8080",
				},
				FileStorage: storageConf.Config{
					FileStoragePath: "",
				},
				DB: db.Config{
					DatabaseDsn: "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := tt.initialConfig
			applyJSONConfig(&conf, tt.jsonConfig)

			if conf.Handlers.ServerAddr != tt.expectedConfig.Handlers.ServerAddr {
				t.Errorf("Handlers.ServerAddr: expected %s, got %s", tt.expectedConfig.Handlers.ServerAddr, conf.Handlers.ServerAddr)
			}

			if conf.Handlers.EnableHTTPS != tt.expectedConfig.Handlers.EnableHTTPS {
				t.Errorf("Handlers.EnableHTTPS: expected %v, got %v", tt.expectedConfig.Handlers.EnableHTTPS, conf.Handlers.EnableHTTPS)
			}

			if conf.Service.ServerURL != tt.expectedConfig.Service.ServerURL {
				t.Errorf("Service.ServerURL: expected %s, got %s", tt.expectedConfig.Service.ServerURL, conf.Service.ServerURL)
			}

			if conf.FileStorage.FileStoragePath != tt.expectedConfig.FileStorage.FileStoragePath {
				t.Errorf("FileStorage.FileStoragePath: expected %s, got %s", tt.expectedConfig.FileStorage.FileStoragePath, conf.FileStorage.FileStoragePath)
			}

			if conf.DB.DatabaseDsn != tt.expectedConfig.DB.DatabaseDsn {
				t.Errorf("DB.DatabaseDsn: expected %s, got %s", tt.expectedConfig.DB.DatabaseDsn, conf.DB.DatabaseDsn)
			}
		})
	}
}

func Test_setJSONStringValue(t *testing.T) {
	type args struct {
		confValue    string
		defaultValue string
		jsonValue    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should not override if not default",
			args: args{
				confValue:    ":9090",
				defaultValue: ":8080",
				jsonValue:    "localhost:8081",
			},
			want: ":9090",
		},
		{
			name: "should override if default and json value is not empty",
			args: args{
				confValue:    ":8080",
				defaultValue: ":8080",
				jsonValue:    "localhost:8081",
			},
			want: "localhost:8081",
		},
		{
			name: "should not override if json value is empty",
			args: args{
				confValue:    ":8085",
				defaultValue: ":8080",
				jsonValue:    "",
			},
			want: ":8085",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setJSONStringValue(&tt.args.confValue, tt.args.defaultValue, tt.args.jsonValue)
			if tt.args.confValue != tt.want {
				t.Errorf("Expected %s, got %s", tt.want, tt.args.confValue)
			}
		})
	}
}

func Test_setJSONBoolValue(t *testing.T) {
	type args struct {
		confValue    bool
		defaultValue bool
		jsonValue    bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should not override if not default",
			args: args{
				confValue:    true,
				defaultValue: false,
				jsonValue:    false,
			},
			want: true,
		},
		{
			name: "should override if default, regardless of json value",
			args: args{
				confValue:    false,
				defaultValue: false,
				jsonValue:    true,
			},
			want: true,
		},
		{
			name: "should override if default and json value is false",
			args: args{
				confValue:    false,
				defaultValue: false,
				jsonValue:    false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setJSONBoolValue(&tt.args.confValue, tt.args.defaultValue, tt.args.jsonValue)
			if tt.args.confValue != tt.want {
				t.Errorf("Expected %v, got %v", tt.want, tt.args.confValue)
			}
		})
	}
}
