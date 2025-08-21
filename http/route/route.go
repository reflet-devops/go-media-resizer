package route

const (
	HealthCheckPingRoute = "/health/ping"
	CgiExtraResizeRoute  = "/cdn-cgi/image/:options/:source"
	MetricsRoute         = "/metrics"
)

var MandatoryRoutes = []string{
	HealthCheckPingRoute,
}

var CgiExtraRoutes = []string{
	CgiExtraResizeRoute,
}
