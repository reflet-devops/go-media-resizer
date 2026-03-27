package route

const (
	HealthCheckPingRoute = "/health/ping"
	CgiExtraResizeRoute  = "/cdn-cgi/image/:options/:source"
	MetricsRoute         = "/metrics"

	ProjectIdHeader  = "X-Project-Id"
	CacheTagHeader   = "Cache-Tag"
	DebugInfoHeader  = "X-Debug-Info"
)

var MandatoryRoutes = []string{
	HealthCheckPingRoute,
}

var CgiExtraRoutes = []string{
	CgiExtraResizeRoute,
}
