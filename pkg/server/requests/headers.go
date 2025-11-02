package requests

const (
	HeaderAccessControlAllowOrigin = "Access-Control-Allow-Origin"

	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"

	HeaderAccessControlAllowHeaders = "Access-Control-Allow-Headers"

	HeaderAccessControlAllowMethods = "Access-Control-Allow-Methods"
)

const (
	ValueAllowOriginWildcard = "*"

	ValueAllowCredentialsTrue = "true"

	ValueAllowHeadersDefault = "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"

	ValueAllowMethodsDefault = "POST, OPTIONS, GET, PUT, DELETE"
)
