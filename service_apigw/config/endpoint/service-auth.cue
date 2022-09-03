package config

_service_auth_host: #backend_hosts.Auth

#ServiceAuthGRPC: "plugin/http-client": {
	name: "grpc-gateway-service_auth"
	endpoints: [ _service_auth_host]
}

#ServiceAuthEndpoints: [
	{
		endpoint: "/register"
		method:   "POST"
		_all_headers_and_querystrings
		backend: [
			{
				url_pattern: "/v1/service_auth/register"
				method:      "POST"
				extra_config: {
					#ServiceAuthGRPC
				}
			},
		]
	},
]
