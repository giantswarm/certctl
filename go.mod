module github.com/giantswarm/certctl/v2

go 1.14

require (
	github.com/giantswarm/apprclient/v2 v2.0.0-20200807082146-02053a5c7c4d
	github.com/giantswarm/backoff v0.2.0
	github.com/giantswarm/go-uuid v0.0.0-20141202165402-ed3ca8a15a93
	github.com/giantswarm/helmclient/v2 v2.0.0-20200807083927-a727a3bb1283
	github.com/giantswarm/k8sclient/v4 v4.0.0-20200806115259-2d3b230ace59
	github.com/giantswarm/microerror v0.2.1
	github.com/giantswarm/micrologger v0.3.1
	github.com/giantswarm/vaultrole v0.2.0
	github.com/hashicorp/vault/api v1.0.4
	github.com/juju/errgo v0.0.0-20140925100237-08cceb5d0b53
	github.com/spf13/afero v1.3.3
	github.com/spf13/cobra v1.0.0
	k8s.io/api v0.18.5
	k8s.io/apimachinery v0.18.5
)

replace (
	github.com/ghodss/yaml => github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/lib/pq => github.com/lib/pq v1.3.0
	github.com/mailru/easyjson => github.com/mailru/easyjson v0.7.0
	github.com/mattn/go-colorable => github.com/mattn/go-colorable v0.0.9
	github.com/mattn/go-isatty => github.com/mattn/go-isatty v0.0.9
	github.com/mattn/go-runewidth => github.com/mattn/go-runewidth v0.0.4
)
