module github.com/giantswarm/certctl

go 1.13

require (
	github.com/DATA-DOG/go-sqlmock v1.4.1 // indirect
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.4.2 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/chai2010/gettext-go v0.0.0-20170215093142-bf70f2a70fb1 // indirect
	github.com/emicklei/go-restful v2.11.1+incompatible // indirect
	github.com/fortytw2/leaktest v1.3.0 // indirect
	github.com/frankban/quicktest v1.4.1 // indirect
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32 // indirect
	github.com/giantswarm/apiextensions v0.0.0-20191209114846-a4fd7939e26e // indirect
	github.com/giantswarm/apprclient v0.0.0-20191209123802-955b7e89e6e2 // indirect
	github.com/giantswarm/backoff v0.0.0-20190913091243-4dd491125192
	github.com/giantswarm/e2e-harness v0.1.1-0.20191209145255-b2048d8645c1
	github.com/giantswarm/e2etemplates v0.0.0-20191115190313-115e721c4fee
	github.com/giantswarm/errors v0.0.0-20191119093322-2640113be13f // indirect
	github.com/giantswarm/go-uuid v0.0.0-20141202165402-ed3ca8a15a93
	github.com/giantswarm/helmclient v0.0.0-20191209124624-d3c47349776d
	github.com/giantswarm/k8sclient v0.0.0-20191209120459-6cb127468cd6
	github.com/giantswarm/k8sportforward v0.0.0-20191209165148-21368288d82d // indirect
	github.com/giantswarm/microerror v0.1.1-0.20200205143715-01b76f66cae6
	github.com/giantswarm/micrologger v0.0.0-20200205144836-079154bcae45
	github.com/go-openapi/spec v0.19.4 // indirect
	github.com/go-test/deep v1.0.2 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/groupcache v0.0.0-20191027212112-611e8accdfc9 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/hashicorp/go-retryablehttp v0.6.4 // indirect
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/hashicorp/vault/api v1.0.5-0.20191108163347-bdd38fca2cff
	github.com/hashicorp/vault/sdk v0.1.14-0.20191112033314-390e96e22eb2 // indirect
	github.com/huandu/xstrings v1.2.1 // indirect
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/jmoiron/sqlx v1.2.0 // indirect
	github.com/juju/errgo v0.0.0-20140925100237-08cceb5d0b53
	github.com/mitchellh/reflectwalk v1.0.1 // indirect
	github.com/pierrec/lz4 v2.4.0+incompatible // indirect
	github.com/prometheus/client_model v0.1.0 // indirect
	github.com/prometheus/common v0.7.0 // indirect
	github.com/prometheus/procfs v0.0.8 // indirect
	github.com/rubenv/sql-migrate v0.0.0-20191121092708-da1cb182f00e // indirect
	github.com/spf13/cobra v0.0.5
	go.uber.org/atomic v1.4.0 // indirect
	golang.org/x/crypto v0.0.0-20191206172530-e9b2fee46413 // indirect
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553 // indirect
	golang.org/x/oauth2 v0.0.0-20191202225959-858c2ad4c8b6 // indirect
	golang.org/x/sys v0.0.0-20191206220618-eeba5f6aabab // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543 // indirect
	google.golang.org/genproto v0.0.0-20191206224255-0243a4be9c8f // indirect
	google.golang.org/grpc v1.25.1 // indirect
	gopkg.in/square/go-jose.v2 v2.4.0 // indirect
	k8s.io/apimachinery v0.16.6
	k8s.io/helm v2.16.1+incompatible // indirect
	k8s.io/kubernetes v1.16.6 // indirect
	sigs.k8s.io/controller-runtime v0.4.0 // indirect
)

replace (
	// All of that is because helm has an import to k8s.io/kubernetes which
	// uses relative paths to those.
	k8s.io/api v0.0.0 => k8s.io/api v0.16.6
	k8s.io/apiextensions-apiserver v0.0.0 => k8s.io/apiextensions-apiserver v0.16.6
	k8s.io/apimachinery v0.0.0 => k8s.io/apimachinery v0.16.6
	k8s.io/apiserver v0.0.0 => k8s.io/apiserver v0.16.6
	k8s.io/cli-runtime v0.0.0 => k8s.io/cli-runtime v0.16.6
	k8s.io/client-go v0.0.0 => k8s.io/client-go v0.16.6
	k8s.io/cloud-provider v0.0.0 => k8s.io/cloud-provider v0.16.6
	k8s.io/cluster-bootstrap v0.0.0 => k8s.io/cluster-bootstrap v0.16.6
	k8s.io/code-generator v0.0.0 => k8s.io/code-generator v0.16.6
	k8s.io/component-base v0.0.0 => k8s.io/component-base v0.16.6
	k8s.io/cri-api v0.0.0 => k8s.io/cri-api v0.16.6
	k8s.io/csi-translation-lib v0.0.0 => k8s.io/csi-translation-lib v0.16.6
	k8s.io/kube-aggregator v0.0.0 => k8s.io/kube-aggregator v0.16.6
	k8s.io/kube-controller-manager v0.0.0 => k8s.io/kube-controller-manager v0.16.6
	k8s.io/kube-proxy v0.0.0 => k8s.io/kube-proxy v0.16.6
	k8s.io/kube-scheduler v0.0.0 => k8s.io/kube-scheduler v0.16.6
	k8s.io/kubectl v0.0.0 => k8s.io/kubectl v0.16.6
	k8s.io/kubelet v0.0.0 => k8s.io/kubelet v0.16.6
	k8s.io/legacy-cloud-providers v0.0.0 => k8s.io/legacy-cloud-providers v0.16.6
	k8s.io/metrics v0.0.0 => k8s.io/metrics v0.16.6
	k8s.io/sample-apiserver v0.0.0 => k8s.io/sample-apiserver v0.16.6
)
