package config

_service_area_host: #backend_hosts.Area

#ServiceAreaGRPC: "plugin/http-client": {
	name: "grpc-gateway-service_area"
	endpoints: [ _service_area_host]
}

#ServiceAreaEndpoints: [
	{
		endpoint: "/get_area_info"
		method:   "GET"
		_all_headers_and_querystrings
		backend: [
			{
				url_pattern: "/v1/service_area/get_area_info"
				method:      "GET"
				extra_config: {
					#ServiceAreaGRPC
				}
			},
		]
	},
]
