package InfiniteMITMDomainsModule

type domains struct {
	Root      string
	Authoring string
	Discovery string
}

var HaloWaypointSVCDomains = domains{
	Root: "svc.halowaypoint.com:443",
	Authoring: "authoring-infiniteugc.svc.halowaypoint.com:443",
	Discovery: "discovery-infiniteugc.svc.halowaypoint.com:443",
}
