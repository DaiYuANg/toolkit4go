module github.com/DaiYuANg/arcgo/examples/rbac_backend

go 1.26.1

replace github.com/DaiYuANg/arcgo/collectionx => ../collectionx

replace github.com/DaiYuANg/arcgo/logx => ../logx

replace github.com/DaiYuANg/arcgo/observabilityx => ../observabilityx

replace github.com/DaiYuANg/arcgo/pkg => ../pkg

replace github.com/DaiYuANg/arcgo/httpx => ../httpx

replace github.com/DaiYuANg/arcgo/bunx => ../bunx

require (
	github.com/samber/lo v1.53.0
	github.com/samber/mo v1.16.0
	github.com/stretchr/testify v1.11.1
	github.com/uptrace/bun v1.2.18
	golang.org/x/crypto v0.48.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/puzpuzpuz/xsync/v3 v3.5.1 // indirect
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
