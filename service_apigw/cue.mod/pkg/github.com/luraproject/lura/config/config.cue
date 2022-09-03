package config

import "go.dev/time"

// HTTPMethod allows known HTTP methods
#HTTPMethod: "GET" | "POST" | "OPTIONS" | "DELETE" | "PATCH" | "PUT"

// ServiceConfig defines the lura service
#ServiceConfig: {
	// Krakend schema version
	"$id": "https://www.krakend.io/schema/v3.json"
	
	// name of the service
	name: string

	// set of endpoint definitions
	endpoints: [...null | #EndpointConfig] @go(,[]*EndpointConfig)

	// default timeout
	timeout?: time.#Duration

	// default TTL for GET
	cache_ttl?: time.#Duration

	// default set of hosts
	host: [...string] @go(,[]string)

	// port to bind the lura service
	port: int

	// version code of the configuration
	version: int

	// output_encoding defines the default encoding strategy to use for the endpoint responses
	output_encoding: string | *"no-op"

	// Extra configuration for customized behaviour
	extra_config: #ExtraConfig

	// read_timeout is the maximum duration for reading the entire
	// request, including the body.
	//
	// Because ReadTimeout does not let Handlers make per-request
	// decisions on each request body's acceptable deadline or
	// upload rate, most users will prefer to use
	// ReadHeaderTimeout. It is valid to use them both.
	read_timeout?: time.#Duration

	// write_timeout is the maximum duration before timing out
	// writes of the response. It is reset whenever a new
	// request's header is read. Like ReadTimeout, it does not
	// let Handlers make decisions on a per-request basis.
	write_timeout?: time.#Duration

	// idle_timeout is the maximum amount of time to wait for the
	// next request when keep-alives are enabled. If IdleTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, ReadHeaderTimeout is used.
	idle_timeout?: time.#Duration

	// read_header_timeout is the amount of time allowed to read
	// request headers. The connection's read deadline is reset
	// after reading the headers and the Handler can decide what
	// is considered too slow for the body.
	read_header_timeout?: time.#Duration

	// disable_keep_alives, if true, prevents re-use of TCP connections
	// between different HTTP requests.
	disable_keep_alives: bool | *false

	// disable_compression, if true, prevents the Transport from
	// requesting compression with an "Accept-Encoding: gzip"
	// request header when the Request contains no existing
	// Accept-Encoding value. If the Transport requests gzip on
	// its own and gets a gzipped response, it's transparently
	// decoded in the Response.Body. However, if the user
	// explicitly requested gzip it is not automatically
	// uncompressed.
	disable_compression: bool | *false

	// MaxIdleConns controls the maximum number of idle (keep-alive)
	// connections across all hosts. Zero means no limit.
	max_idle_conns: int | *0

	// MaxIdleConnsPerHost, if non-zero, controls the maximum idle
	// (keep-alive) connections to keep per-host. If zero,
	// DefaultMaxIdleConnsPerHost is used.
	max_idle_conns_per_host: int | *0

	// IdleConnTimeout is the maximum amount of time an idle
	// (keep-alive) connection will remain idle before closing
	// itself.
	// Zero means no limit.
	idle_conn_timeout?: time.#Duration

	// ResponseHeaderTimeout, if non-zero, specifies the amount of
	// time to wait for a server's response headers after fully
	// writing the request (including its body, if any). This
	// time does not include the time to read the response body.
	response_header_timeout?: time.#Duration

	// ExpectContinueTimeout, if non-zero, specifies the amount of
	// time to wait for a server's first response headers after fully
	// writing the request headers if the request has an
	// "Expect: 100-continue" header. Zero means no timeout and
	// causes the body to be sent immediately, without
	// waiting for the server to approve.
	// This time does not include the time to send the request header.
	expect_continue_timeout?: time.#Duration

	// DialerTimeout is the maximum amount of time a dial will wait for
	// a connect to complete. If Deadline is also set, it may fail
	// earlier.
	//
	// The default is no timeout.
	//
	// When using TCP and dialing a host name with multiple IP
	// addresses, the timeout may be divided between them.
	//
	// With or without a timeout, the operating system may impose
	// its own earlier timeout. For instance, TCP timeouts are
	// often around 3 minutes.
	dialer_timeout?: time.#Duration

	// DialerFallbackDelay specifies the length of time to wait before
	// spawning a fallback connection, when DualStack is enabled.
	// If zero, a default delay of 300ms is used.
	dialer_fallback_delay: time.#Duration | *(300 * time.#Millisecond)

	// DialerKeepAlive specifies the keep-alive period for an active
	// network connection.
	// If zero, keep-alives are not enabled. Network protocols
	// that do not support keep-alives ignore this field.
	dialer_keep_alive?: string

	// DisableStrictREST flags if the REST enforcement is disabled
	disable_strict_rest: bool | *false

	// Plugin defines the configuration for the plugin loader
	plugin?: null | #Plugin @go(,*Plugin)

	// TLS defines the configuration params for enabling TLS (HTTPS & HTTP/2) at
	// the router layer
	tls?: null | #TLS @go(,*TLS)

	// run lura in debug mode
	debug: bool | *false

	// SequentialStart flags if the agents should be started sequentially
	// before starting the router
	sequential_start: bool | *true
}

// EndpointConfig defines the configuration of a single endpoint to be exposed
// by the lura service
#EndpointConfig: {
	// url pattern to be registered and exposed to the world
	endpoint: string

	// HTTP method of the endpoint (GET, POST, PUT, etc)
	method: #HTTPMethod

	// set of definitions of the backends to be linked to this endpoint
	backend: [...null | #Backend] @go(,[]*Backend)

	// number of concurrent calls this endpoint must send to the backends
	concurrent_calls?: int

	// timeout of this endpoint
	timeout?: time.#Duration

	// duration of the cache header
	cache_ttl?: time.#Duration

	// list of query string params to be extracted from the URI
	input_query_strings: [...string] @go(,[]string)

	// Endpoint Extra configuration for customized behaviour
	extra_config: #ExtraConfig

	// InputHeaders defines the list of headers to pass to the backends
	input_headers: [...string] @go(,[]string)

	// OutputEncoding defines the encoding strategy to use for the endpoint responses
	output_encoding: string | *"no-op"
}

// Backend defines how lura should connect to the backend service (the API resource to consume)
// and how it should process the received response
#Backend: {
	// Group defines the name of the property the response should be moved to. If empty, the response is
	// not changed
	group?: string

	// Method defines the HTTP method of the request to send to the backend
	method: #HTTPMethod

	// Host is a set of hosts of the API
	host?: [...string] @go(,[]string)

	// HostSanitizationDisabled can be set to false if the hostname should be sanitized
	host_sanitization_disabled?: bool

	// URLPattern is the URL pattern to use to locate the resource to be consumed
	url_pattern: string

	// Deprecated: use DenyList
	// Blacklist is a set of response fields to remove. If empty, the filter id not used
	blacklist?: [...string] @go(,[]string)

	// Deprecated: use AllowList
	// Whitelist is a set of response fields to allow. If empty, the filter id not used
	whitelist?: [...string] @go(,[]string)

	// AllowList is a set of response fields to allow. If empty, the filter id not used
	allow_list?: [...string] @go(,[]string)

	// DenyList is a set of response fields to remove. If empty, the filter id not used
	deny_list?: [...string] @go(,[]string)

	// map of response fields to be renamed and their new names
	mapping?: {[string]: string} @go(,map[string]string)

	// the encoding format
	encoding: string | *"no-op"

	// the response to process is a collection
	is_collection: bool | *false

	// name of the field to extract to the root. If empty, the formater will do nothing
	target?: string

	// name of the service discovery driver to use
	sd?: string

	// list of keys to be replaced in the URLPattern
	url_keys?: [...string] @go(,[]string)

	// number of concurrent calls this endpoint must send to the API
	concurrent_calls?: int

	// timeout of this backend
	timeout?: time.#Duration

	// Backend Extra configuration for customized behaviours
	extra_config?: #ExtraConfig
}

// Plugin contains the config required by the plugin module
#Plugin: {
	folder:  string
	pattern: string
}

// TLS defines the configuration params for enabling TLS (HTTPS & HTTP/2) at the router layer
#TLS: {
	is_disabled: bool
	public_key:  string
	private_key: string
	min_version: string
	max_version: string
	curve_preferences: [...uint16] @go(,[]uint16)
	prefer_server_cipher_suites: bool
	cipher_suites: [...uint16] @go(,[]uint16)
	enable_mtls: bool
}

// ExtraConfig is a type to store extra configurations for customized behaviours
#ExtraConfig: {...}
