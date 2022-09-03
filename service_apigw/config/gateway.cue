package config

import (
	lura "github.com/luraproject/lura/config"
	"go.dev/time"
)

// set type
#KrakendConf: lura.#ServiceConfig

// define type with value
#KrakendConf: {
	name:            "service"
	version:         3
	output_encoding: "json"
	timeout:         3 * time.#Second
	port:            8080
	cache_ttl:       1 * time.#Hour
	host:            #default_hosts
	debug:           #debug_mode
	plugin: {
		pattern: ".so"
		folder:  "./dist/plugins/"
	}
	extra_config: close(#Extra)
	endpoints:    #AllEndpoints + []
}

// this line below used
// to generate.sh flat final JSON:
#KrakendConf
