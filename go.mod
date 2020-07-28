module github.com/giantswarm/certctl

go 1.13

require (
	github.com/giantswarm/apprclient v0.2.1-0.20200724085653-63c7eb430dcf
	github.com/giantswarm/backoff v0.2.0
	github.com/giantswarm/go-uuid v0.0.0-20141202165402-ed3ca8a15a93
	github.com/giantswarm/helmclient v1.0.6-0.20200724131413-ea0311052b6e
	github.com/giantswarm/k8sclient/v3 v3.1.3-0.20200724085258-345602646ea8
	github.com/giantswarm/microerror v0.2.0
	github.com/giantswarm/micrologger v0.3.1
	github.com/giantswarm/vaultrole v0.2.0
	github.com/hashicorp/vault/api v1.0.4
	github.com/juju/errgo v0.0.0-20140925100237-08cceb5d0b53
	github.com/spf13/afero v1.3.2
	github.com/spf13/cobra v1.0.0
	k8s.io/api v0.18.5
	k8s.io/apimachinery v0.18.5
)

replace (
	github.com/ghodss/yaml => github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/mailru/easyjson => github.com/mailru/easyjson v0.7.0
	github.com/mattn/go-colorable => github.com/mattn/go-colorable v0.0.9
	github.com/mattn/go-isatty => github.com/mattn/go-isatty v0.0.9
)
