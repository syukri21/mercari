package config

import lura "github.com/luraproject/lura/config"

#Extra: {
	"auth/revoker": {
		n:         10000000
		p:         0.0000001
		hash_name: "optimal"
		ttl:       1500
		port:      1234
		token_keys: [
			"jti",
			"did",
		]
	}
	"telemetry/logging": {
		level:         "DEBUG"
		prefix:        "[KRAKEND]"
		syslog:        false
		stdout:        true
		format:        "custom"
		custom_format: "{\"timestamp\":\"%{time:2006-01-02T15:04:05.000+00:00}\", \"level\": \"%{level}\", \"message\": \"%{message}\", \"module\": \"%{module}\"}"
	}
	"security/cors": {
		allow_origins: [
			"*",
		]
		allow_methods: [
			"POST",
			"GET",
			"PUT",
			"DELETE",
			"PATCH",
			"OPTIONS",
		]
		allow_headers: [
			"Origin",
			"Authorization",
			"Content-Type",
			"Access-Control-Allow-Origin",
			"X-Securities-Request-PartnerID",
			"X-Securities-Request-Timestamp",
			"X-Securities-Signature",
			"X-Securities-ID",
			"Lang",
		]
		allow_credentials: false
		expose_headers: [
			"Content-Length",
		]
		max_age: "12h"
		debug:   false
	}
} & lura.#ExtraConfig
