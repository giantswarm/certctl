module github.com/giantswarm/certctl/v2

go 1.14

require (
	github.com/asaskevich/govalidator v0.0.0-20200108200545-475eaeb16496 // indirect
	github.com/giantswarm/appcatalog v0.3.1
	github.com/giantswarm/apptest v0.7.1
	github.com/giantswarm/backoff v0.2.0
	github.com/giantswarm/go-uuid v0.0.0-20141202165402-ed3ca8a15a93
	github.com/giantswarm/k8sclient/v4 v4.0.0
	github.com/giantswarm/microerror v0.2.1
	github.com/giantswarm/micrologger v0.3.4
	github.com/giantswarm/vaultrole v0.2.0
	github.com/hashicorp/vault/api v1.0.4
	github.com/juju/errgo v0.0.0-20140925100237-08cceb5d0b53
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.5.1 // indirect
	golang.org/x/crypto v0.0.0-20200414173820-0848c9571904 // indirect
	k8s.io/api v0.18.9
	k8s.io/apimachinery v0.18.9
	sigs.k8s.io/yaml v1.2.0
)

replace (
	// Apply security fix not present in 3.3.10
	github.com/coreos/etcd v3.3.10+incompatible => github.com/coreos/etcd v3.3.24+incompatible
	github.com/ghodss/yaml => github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	// Apply security fix not present in v1.4.0.
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2
	github.com/lib/pq => github.com/lib/pq v1.3.0
	github.com/mailru/easyjson => github.com/mailru/easyjson v0.7.0
	github.com/mattn/go-colorable => github.com/mattn/go-colorable v0.0.9
	github.com/mattn/go-isatty => github.com/mattn/go-isatty v0.0.9
	github.com/mattn/go-runewidth => github.com/mattn/go-runewidth v0.0.4
	sigs.k8s.io/cluster-api => github.com/giantswarm/cluster-api v0.3.10-gs
)
