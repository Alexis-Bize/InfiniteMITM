package InfiniteMITMDomainsModule

type domains struct {
	Root      string
	Authoring string
	Discovery string
	HaloStats string
}

var HaloWaypointSVCDomains = domains{
	Root: "svc.halowaypoint.com:443",
	Authoring: "authoring-infiniteugc.svc.halowaypoint.com:443",
	Discovery: "discovery-infiniteugc.svc.halowaypoint.com:443",
	HaloStats: "halostats.svc.halowaypoint.com:443",
}
