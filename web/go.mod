module github.com/Sn0rt/secret2es/web

go 1.23

require (
	github.com/Sn0rt/secret2es v0.0.0
	github.com/external-secrets/external-secrets v0.10.2
)

replace (
 github.com/Sn0rt/secret2es/pkg/converter => ../pkg/converter
 github.com/Sn0rt/secret2es => ../
)
