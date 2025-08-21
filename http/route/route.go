package route

const RouteHealthCheckPing = "/health/ping"
const RouteCgiExtraResize = "/cdn-cgi/image/:options/:source"

var MandatoryRoutes = []string{
	RouteHealthCheckPing,
}

var CgiExtraRoutes = []string{
	RouteCgiExtraResize,
}
