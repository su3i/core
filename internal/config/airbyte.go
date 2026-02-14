package config

type AirbyteConfig struct {
	AirbyteEndpoint string `required:"true"`
	AirbyteClientId string `required:"true"`
	AirbyteClientSecret string `required:"true"`
	AirbyteWorkspaceId string `required:"true"`
}