package httpx

func runEndpointHooks(server ServerRuntime, endpoint Endpoint, hooks []EndpointHooks, selectHook func(EndpointHooks) EndpointHookFunc) {
	for _, hook := range hooks {
		selected := selectHook(hook)
		if selected != nil {
			selected(server, endpoint)
		}
	}
}
