package route

const (
	HealthCheckPingRoute = "/health/ping"
	CgiExtraResizeRoute  = "/cdn-cgi/image/:options/:source"
	MetricsRoute         = "/metrics"

	ProjectIdHeader = "X-Project-Id"
	CacheTagHeader  = "Cache-Tag"
)

var MandatoryRoutes = []string{
	HealthCheckPingRoute,
}

var CgiExtraRoutes = []string{
	CgiExtraResizeRoute,
}
