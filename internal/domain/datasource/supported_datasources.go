package datasource

var SupportedDatasources = []map[string]interface{}{
	// ─────────────────────────────────────────────
	// Postgres
	// ─────────────────────────────────────────────
	{
		"name": "Postgres",
		"logo": "",
		"desc": "PostgreSQL is a powerful, open-source object-relational database system.",
		"sourceType": "postgres",
		"form": []map[string]interface{}{
			{
				"title":       "host",
				"desc":        "Hostname of the database.",
				"type":        "string",
				"required":    true,
				"placeholder": "localhost",
			},
			{
				"title":       "port",
				"desc":        "Port of the database.",
				"type":        "integer",
				"required":    true,
				"default":     5432,
				"min":         0,
				"max":         65536,
				"placeholder": "5432",
			},
			{
				"title":       "database",
				"desc":        "Name of the database.",
				"type":        "string",
				"required":    true,
				"placeholder": "test",
			},
			{
				"title":       "schemas",
				"desc":        "The list of schemas (case sensitive) to sync from. Defaults to public.",
				"type":        "array",
				"items":       "string",
				"default":     []string{"public"},
				"minLength":   0,
			},
			{
				"title":       "username",
				"desc":        "Username to access the database.",
				"type":        "string",
				"required":    true,
				"placeholder": "test",
			},
			{
				"title":       "password",
				"desc":        "Password associated with the username.",
				"type":        "string",
				"secret":      true,
				"placeholder": "test",
			},
			{
				"title": "jdbc_url_params",
				"desc":  "Additional properties to pass to the JDBC URL string when connecting to the database formatted as 'key=value' pairs separated by the symbol '&'. (Eg. key1=value1&key2=value2&key3=value3).",
				"type":  "string",
			},
			{
				"title": "ssl_mode",
				"desc":  "SSL connection modes. Read more in the docs.",
				"type":  "object",
				"oneOf": []map[string]interface{}{
					{"title": "disable", "value": "disable"},
					{"title": "allow", "value": "allow"},
					{"title": "prefer", "value": "prefer"},
					{"title": "require", "value": "require"},
					{
						"title": "verify-ca",
						"value": "verify-ca",
						"fields": []map[string]interface{}{
							{"title": "ca_certificate", "desc": "CA certificate", "type": "string", "required": true, "multiline": true},
							{"title": "client_certificate", "desc": "Client certificate", "type": "string", "multiline": true},
							{"title": "client_key", "desc": "Client key", "type": "string", "secret": true, "multiline": true},
							{"title": "client_key_password", "desc": "Client key password", "type": "string", "secret": true},
						},
					},
					{
						"title": "verify-full",
						"value": "verify-full",
						"fields": []map[string]interface{}{
							{"title": "ca_certificate", "desc": "CA certificate", "type": "string", "required": true, "multiline": true},
							{"title": "client_certificate", "desc": "Client certificate", "type": "string", "multiline": true},
							{"title": "client_key", "desc": "Client key", "type": "string", "secret": true, "multiline": true},
							{"title": "client_key_password", "desc": "Client key password", "type": "string", "secret": true},
						},
					},
				},
			},
			{
				"title": "replication_method",
				"desc":  "Configures how data is extracted from the database.",
				"type":  "object",
				"oneOf": []map[string]interface{}{
					{
						"title": "Read Changes using Write-Ahead Log (CDC)",
						"value": "CDC",
						"fields": []map[string]interface{}{
							{"title": "plugin", "desc": "A logical decoding plugin installed on the PostgreSQL server.", "type": "string", "default": "pgoutput", "enum": []string{"pgoutput"}},
							{"title": "replication_slot", "desc": "A plugin logical replication slot.", "type": "string", "required": true},
							{"title": "publication", "desc": "A Postgres publication used for consuming changes.", "type": "string", "required": true},
							{"title": "initial_waiting_seconds", "desc": "The amount of time the connector will wait when it launches to determine if there is new data to sync.", "type": "integer", "default": 300, "min": 120, "max": 1200},
							{"title": "lsn_commit_behaviour", "desc": "Determines when Airbyte should flush the LSN.", "type": "string", "default": "After loading Data in the destination", "enum": []string{"After loading Data in the destination", "While reading Data"}},
						},
					},
					{
						"title": "Detect Changes with Xmin System Column",
						"value": "Xmin",
					},
					{
						"title": "Scan Changes with User Defined Cursor",
						"value": "Standard",
					},
				},
			},
			{
				"title": "tunnel_method",
				"desc":  "Whether to initiate an SSH tunnel before connecting to the database, and if so, which kind of authentication to use.",
				"type":  "object",
				"oneOf": []map[string]interface{}{
					{
						"title": "No Tunnel",
						"value": "NO_TUNNEL",
					},
					{
						"title": "SSH Key Authentication",
						"value": "SSH_KEY_AUTH",
						"fields": []map[string]interface{}{
							{"title": "tunnel_host", "desc": "Hostname of the jump server host that allows inbound SSH tunnel.", "type": "string", "required": true},
							{"title": "tunnel_port", "desc": "Port on the proxy/jump server that accepts inbound SSH connections.", "type": "integer", "required": true, "default": 22},
							{"title": "tunnel_user", "desc": "OS-level username for logging into the jump server host.", "type": "string", "required": true},
							{"title": "ssh_key", "desc": "OS-level user account ssh key credentials in RSA PEM format.", "type": "string", "required": true, "secret": true, "multiline": true},
						},
					},
					{
						"title": "Password Authentication",
						"value": "SSH_PASSWORD_AUTH",
						"fields": []map[string]interface{}{
							{"title": "tunnel_host", "desc": "Hostname of the jump server host that allows inbound SSH tunnel.", "type": "string", "required": true},
							{"title": "tunnel_port", "desc": "Port on the proxy/jump server that accepts inbound SSH connections.", "type": "integer", "required": true, "default": 22},
							{"title": "tunnel_user", "desc": "OS-level username for logging into the jump server host.", "type": "string", "required": true},
							{"title": "tunnel_user_password", "desc": "OS-level password for logging into the jump server host.", "type": "string", "required": true, "secret": true},
						},
					},
				},
			},
			{
				"title":    "sourceType",
				"desc":     "Source type identifier.",
				"type":     "const",
				"value":    "postgres",
				"enum":     []string{"postgres"},
				"required": true,
				"hidden":   true,
			},
		},
	},

	// ─────────────────────────────────────────────
	// MySQL
	// ─────────────────────────────────────────────
	{
		"name":       "MySQL",
		"logo":       "",
		"desc":       "MySQL is the world's most popular open-source relational database management system.",
		"sourceType": "mysql",
		"form": []map[string]interface{}{
			{
				"title":       "host",
				"desc":        "The host name of the database.",
				"type":        "string",
				"required":    true,
				"placeholder": "localhost",
			},
			{
				"title":       "port",
				"desc":        "The port to connect to.",
				"type":        "integer",
				"required":    true,
				"default":     3306,
				"min":         0,
				"max":         65536,
				"placeholder": "3306",
			},
			{
				"title":       "database",
				"desc":        "The database name.",
				"type":        "string",
				"required":    true,
				"placeholder": "test",
			},
			{
				"title":       "username",
				"desc":        "The username which is used to access the database.",
				"type":        "string",
				"required":    true,
				"placeholder": "test",
			},
			{
				"title":       "password",
				"desc":        "The password associated with the username.",
				"type":        "string",
				"secret":      true,
				"placeholder": "test",
			},
			{
				"title": "jdbc_url_params",
				"desc":  "Additional properties to pass to the JDBC URL string when connecting to the database formatted as 'key=value' pairs separated by the symbol '&'. (example: key1=value1&key2=value2&key3=value3).",
				"type":  "string",
			},
			{
				"title":   "ssl",
				"desc":    "Encrypt data using SSL.",
				"type":    "boolean",
				"default": true,
			},
			{
				"title": "ssl_mode",
				"desc":  "SSL connection modes. Read more in the docs.",
				"type":  "object",
				"oneOf": []map[string]interface{}{
					{"title": "preferred", "value": "preferred"},
					{"title": "required", "value": "required"},
					{
						"title": "Verify CA",
						"value": "verify_ca",
						"fields": []map[string]interface{}{
							{"title": "ca_certificate", "desc": "CA certificate", "type": "string", "required": true, "multiline": true},
							{"title": "client_certificate", "desc": "Client certificate", "type": "string", "multiline": true},
							{"title": "client_key", "desc": "Client key", "type": "string", "secret": true, "multiline": true},
							{"title": "client_key_password", "desc": "Client key password", "type": "string", "secret": true},
						},
					},
					{
						"title": "Verify Identity",
						"value": "verify_identity",
						"fields": []map[string]interface{}{
							{"title": "ca_certificate", "desc": "CA certificate", "type": "string", "required": true, "multiline": true},
							{"title": "client_certificate", "desc": "Client certificate", "type": "string", "multiline": true},
							{"title": "client_key", "desc": "Client key", "type": "string", "secret": true, "multiline": true},
							{"title": "client_key_password", "desc": "Client key password", "type": "string", "secret": true},
						},
					},
				},
			},
			{
				"title":    "replication_method",
				"desc":     "Configures how data is extracted from the database.",
				"type":     "object",
				"required": true,
				"oneOf": []map[string]interface{}{
					{
						"title": "Read Changes using Binary Log (CDC)",
						"value": "CDC",
						"fields": []map[string]interface{}{
							{"title": "initial_waiting_seconds", "desc": "The amount of time the connector will wait when it launches to determine if there is new data to sync.", "type": "integer", "default": 300, "min": 120, "max": 1200},
							{"title": "server_time_zone", "desc": "Enter the configured MySQL server timezone. This should only be done if the configured timezone in your MySQL instance does not conform to IANNA standard.", "type": "string"},
							{"title": "invalid_cdc_cursor_position_behavior", "desc": "Determines how the connector behaves when it detects an invalid CDC cursor position.", "type": "string", "default": "Fail sync", "enum": []string{"Fail sync", "Re-sync data"}},
							{"title": "initial_load_timeout_hours", "desc": "The amount of time an initial load is allowed to continue for before catching up on CDC logs.", "type": "integer", "default": 8, "min": 4, "max": 24},
						},
					},
					{
						"title": "Scan Changes with User Defined Cursor",
						"value": "Standard",
					},
				},
			},
			{
				"title": "tunnel_method",
				"desc":  "Whether to initiate an SSH tunnel before connecting to the database, and if so, which kind of authentication to use.",
				"type":  "object",
				"oneOf": []map[string]interface{}{
					{
						"title": "No Tunnel",
						"value": "NO_TUNNEL",
					},
					{
						"title": "SSH Key Authentication",
						"value": "SSH_KEY_AUTH",
						"fields": []map[string]interface{}{
							{"title": "tunnel_host", "desc": "Hostname of the jump server host that allows inbound SSH tunnel.", "type": "string", "required": true},
							{"title": "tunnel_port", "desc": "Port on the proxy/jump server that accepts inbound SSH connections.", "type": "integer", "required": true, "default": 22},
							{"title": "tunnel_user", "desc": "OS-level username for logging into the jump server host.", "type": "string", "required": true},
							{"title": "ssh_key", "desc": "OS-level user account ssh key credentials in RSA PEM format.", "type": "string", "required": true, "secret": true, "multiline": true},
						},
					},
					{
						"title": "Password Authentication",
						"value": "SSH_PASSWORD_AUTH",
						"fields": []map[string]interface{}{
							{"title": "tunnel_host", "desc": "Hostname of the jump server host that allows inbound SSH tunnel.", "type": "string", "required": true},
							{"title": "tunnel_port", "desc": "Port on the proxy/jump server that accepts inbound SSH connections.", "type": "integer", "required": true, "default": 22},
							{"title": "tunnel_user", "desc": "OS-level username for logging into the jump server host.", "type": "string", "required": true},
							{"title": "tunnel_user_password", "desc": "OS-level password for logging into the jump server host.", "type": "string", "required": true, "secret": true},
						},
					},
				},
			},
			{
				"title":    "sourceType",
				"desc":     "Source type identifier.",
				"type":     "const",
				"value":    "mysql",
				"enum":     []string{"mysql"},
				"required": true,
				"hidden":   true,
			},
		},
	},

	// ─────────────────────────────────────────────
	// Microsoft SQL Server
	// ─────────────────────────────────────────────
	{
		"name":       "Microsoft SQL Server",
		"logo":       "",
		"desc":       "Microsoft SQL Server is a relational database management system developed by Microsoft.",
		"sourceType": "mssql",
		"form": []map[string]interface{}{
			{
				"title":       "host",
				"desc":        "The hostname of the database.",
				"type":        "string",
				"required":    true,
				"placeholder": "localhost",
			},
			{
				"title":       "port",
				"desc":        "The port of the database.",
				"type":        "integer",
				"required":    true,
				"default":     1433,
				"min":         0,
				"max":         65536,
				"placeholder": "1433",
			},
			{
				"title":       "database",
				"desc":        "The name of the database.",
				"type":        "string",
				"required":    true,
				"placeholder": "test",
			},
			{
				"title":     "schemas",
				"desc":      "The list of schemas to sync from. Defaults to dbo. Case sensitive.",
				"type":      "array",
				"items":     "string",
				"default":   []string{"dbo"},
				"minLength": 0,
			},
			{
				"title":       "username",
				"desc":        "The username which is used to access the database.",
				"type":        "string",
				"required":    true,
				"placeholder": "test",
			},
			{
				"title":       "password",
				"desc":        "The password associated with the username.",
				"type":        "string",
				"required":    true,
				"secret":      true,
				"placeholder": "test",
			},
			{
				"title": "jdbc_url_params",
				"desc":  "Additional properties to pass to the JDBC URL string when connecting to the database formatted as 'key=value' pairs separated by the symbol '&'. (example: key1=value1&key2=value2&key3=value3).",
				"type":  "string",
			},
			{
				"title": "ssl_method",
				"desc":  "The encryption method which is used when communicating with the database.",
				"type":  "object",
				"oneOf": []map[string]interface{}{
					{
						"title": "Unencrypted",
						"value": "unencrypted",
					},
					{
						"title": "Encrypted (trust server certificate)",
						"value": "encrypted_trust_server_certificate",
					},
					{
						"title": "Encrypted (verify certificate)",
						"value": "encrypted_verify_certificate",
						"fields": []map[string]interface{}{
							{"title": "ssl_certificate", "desc": "Certificate to use for SSL connection.", "type": "string", "required": true, "multiline": true},
							{"title": "host_name_in_certificate", "desc": "Host name in the certificate CN.", "type": "string"},
						},
					},
				},
			},
			{
				"title": "replication_method",
				"desc":  "Configures how data is extracted from the database.",
				"type":  "object",
				"oneOf": []map[string]interface{}{
					{
						"title": "Read Changes using Change Data Capture (CDC)",
						"value": "CDC",
						"fields": []map[string]interface{}{
							{"title": "data_to_sync", "desc": "Choose how data is synced to the destination.", "type": "string", "default": "Existing and New", "enum": []string{"Existing and New", "New Changes Only"}},
							{"title": "snapshot_isolation", "desc": "Choose between Snapshot and Read Committed isolation levels.", "type": "string", "default": "Snapshot", "enum": []string{"Snapshot", "Read Committed"}},
							{"title": "initial_waiting_seconds", "desc": "The amount of time the connector will wait when it launches to determine if there is new data to sync.", "type": "integer", "default": 300, "min": 0, "max": 1200},
							{"title": "invalid_cdc_cursor_position_behavior", "desc": "Determines how the connector behaves when it detects an invalid CDC cursor position.", "type": "string", "default": "Fail sync", "enum": []string{"Fail sync", "Re-sync data"}},
						},
					},
					{
						"title": "Scan Changes with User Defined Cursor",
						"value": "Standard",
					},
				},
			},
			{
				"title": "tunnel_method",
				"desc":  "Whether to initiate an SSH tunnel before connecting to the database, and if so, which kind of authentication to use.",
				"type":  "object",
				"oneOf": []map[string]interface{}{
					{
						"title": "No Tunnel",
						"value": "NO_TUNNEL",
					},
					{
						"title": "SSH Key Authentication",
						"value": "SSH_KEY_AUTH",
						"fields": []map[string]interface{}{
							{"title": "tunnel_host", "desc": "Hostname of the jump server host that allows inbound SSH tunnel.", "type": "string", "required": true},
							{"title": "tunnel_port", "desc": "Port on the proxy/jump server that accepts inbound SSH connections.", "type": "integer", "required": true, "default": 22},
							{"title": "tunnel_user", "desc": "OS-level username for logging into the jump server host.", "type": "string", "required": true},
							{"title": "ssh_key", "desc": "OS-level user account ssh key credentials in RSA PEM format.", "type": "string", "required": true, "secret": true, "multiline": true},
						},
					},
					{
						"title": "Password Authentication",
						"value": "SSH_PASSWORD_AUTH",
						"fields": []map[string]interface{}{
							{"title": "tunnel_host", "desc": "Hostname of the jump server host that allows inbound SSH tunnel.", "type": "string", "required": true},
							{"title": "tunnel_port", "desc": "Port on the proxy/jump server that accepts inbound SSH connections.", "type": "integer", "required": true, "default": 22},
							{"title": "tunnel_user", "desc": "OS-level username for logging into the jump server host.", "type": "string", "required": true},
							{"title": "tunnel_user_password", "desc": "OS-level password for logging into the jump server host.", "type": "string", "required": true, "secret": true},
						},
					},
				},
			},
			{
				"title":    "sourceType",
				"desc":     "Source type identifier.",
				"type":     "const",
				"value":    "mssql",
				"enum":     []string{"mssql"},
				"required": true,
				"hidden":   true,
			},
		},
	},
}
