module github.com/DaiYuANg/arcgo/clientx

go 1.26.1

replace github.com/DaiYuANg/arcgo/collectionx => ../collectionx

replace github.com/DaiYuANg/arcgo/observabilityx => ../observabilityx

require (
	github.com/DaiYuANg/arcgo/collectionx v0.0.0-00010101000000-000000000000
	github.com/DaiYuANg/arcgo/observabilityx v0.0.0-00010101000000-000000000000
	github.com/samber/lo v1.53.0
	github.com/samber/mo v1.16.0
	resty.dev/v3 v3.0.0-beta.6
)

require (
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/text v0.34.0 // indirect
)
