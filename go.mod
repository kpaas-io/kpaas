module github.com/kpaas-io/kpaas

go 1.13

require (
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/containerd/containerd v1.3.1 // indirect
	github.com/coreos/etcd v3.3.17+incompatible
	github.com/coreos/go-semver v0.3.0
	github.com/docker/docker v1.4.2-0.20191127222017-3152f9436292
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/gin-gonic/gin v1.5.0
	github.com/go-openapi/jsonreference v0.19.3 // indirect
	github.com/go-openapi/spec v0.19.4 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.3.2
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/json-iterator/go v1.1.8 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mailru/easyjson v0.7.0 // indirect
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/pkg/sftp v1.10.1
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd
	github.com/sirupsen/logrus v1.4.2
	github.com/sony/sonyflake v1.0.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.4.0
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.3
	golang.org/x/crypto v0.0.0-20191122220453-ac88ee75c92c
	golang.org/x/net v0.0.0-20191116160921-f9c825593386
	golang.org/x/sys v0.0.0-20191120155948-bd437916bb0e // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	golang.org/x/tools v0.0.0-20191118051429-5a76f03bc7c3 // indirect
	google.golang.org/grpc v1.25.1
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/go-playground/validator.v9 v9.30.2 // indirect
	gopkg.in/yaml.v2 v2.2.5 // indirect
	k8s.io/client-go v0.0.0
	k8s.io/kubernetes v1.16.3
	gotest.tools v2.2.0+incompatible // indirect
)

replace (
	k8s.io/api => k8s.io/kubernetes/staging/src/k8s.io/api v0.0.0-20191207073709-d9a258903170
	k8s.io/apiextensions-apiserver => k8s.io/kubernetes/staging/src/k8s.io/apiextensions-apiserver v0.0.0-20191207073709-d9a258903170
	k8s.io/apimachinery => k8s.io/kubernetes/staging/src/k8s.io/apimachinery v0.0.0-20191207073709-d9a258903170
	k8s.io/apiserver => k8s.io/kubernetes/staging/src/k8s.io/apiserver v0.0.0-20191207073709-d9a258903170
	k8s.io/cli-runtime => k8s.io/kubernetes/staging/src/k8s.io/cli-runtime v0.0.0-20191207073709-d9a258903170
	k8s.io/client-go => k8s.io/kubernetes/staging/src/k8s.io/client-go v0.0.0-20191207073709-d9a258903170
	k8s.io/cloud-provider => k8s.io/kubernetes/staging/src/k8s.io/cloud-provider v0.0.0-20191207073709-d9a258903170
	k8s.io/cluster-bootstrap => k8s.io/kubernetes/staging/src/k8s.io/cluster-bootstrap v0.0.0-20191207073709-d9a258903170
	k8s.io/code-generator => k8s.io/kubernetes/staging/src/k8s.io/code-generator v0.0.0-20191207073709-d9a258903170
	k8s.io/component-base => k8s.io/kubernetes/staging/src/k8s.io/component-base v0.0.0-20191207073709-d9a258903170
	k8s.io/cri-api => k8s.io/kubernetes/staging/src/k8s.io/cri-api v0.0.0-20191207073709-d9a258903170
	k8s.io/csi-translation-lib => k8s.io/kubernetes/staging/src/k8s.io/csi-translation-lib v0.0.0-20191207073709-d9a258903170
	k8s.io/kube-aggregator => k8s.io/kubernetes/staging/src/k8s.io/kube-aggregator v0.0.0-20191207073709-d9a258903170
	k8s.io/kube-controller-manager => k8s.io/kubernetes/staging/src/k8s.io/kube-controller-manager v0.0.0-20191207073709-d9a258903170
	k8s.io/kube-proxy => k8s.io/kubernetes/staging/src/k8s.io/kube-proxy v0.0.0-20191207073709-d9a258903170
	k8s.io/kube-scheduler => k8s.io/kubernetes/staging/src/k8s.io/kube-scheduler v0.0.0-20191207073709-d9a258903170
	k8s.io/kubectl => k8s.io/kubernetes/staging/src/k8s.io/kubectl v0.0.0-20191207073709-d9a258903170
	k8s.io/kubelet => k8s.io/kubernetes/staging/src/k8s.io/kubelet v0.0.0-20191207073709-d9a258903170
	k8s.io/legacy-cloud-providers => k8s.io/kubernetes/staging/src/k8s.io/legacy-cloud-providers v0.0.0-20191207073709-d9a258903170
	k8s.io/metrics => k8s.io/kubernetes/staging/src/k8s.io/metrics v0.0.0-20191207073709-d9a258903170
	k8s.io/node-api => k8s.io/kubernetes/staging/src/k8s.io/node-api v0.0.0-20191207073709-d9a258903170
	k8s.io/sample-apiserver => k8s.io/kubernetes/staging/src/k8s.io/sample-apiserver v0.0.0-20191207073709-d9a258903170
	k8s.io/sample-cli-plugin => k8s.io/kubernetes/staging/src/k8s.io/sample-cli-plugin v0.0.0-20191207073709-d9a258903170
	k8s.io/sample-controller => k8s.io/kubernetes/staging/src/k8s.io/sample-controller v0.0.0-20191207073709-d9a258903170
)
