module github.com/KanchiShimono/csi-driver-hostpath

go 1.19

require (
	github.com/container-storage-interface/spec v1.7.0
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/kubernetes-csi/csi-lib-utils v0.12.0
	golang.org/x/net v0.4.0
	google.golang.org/grpc v1.52.0
	k8s.io/klog/v2 v2.80.1
	k8s.io/kubernetes v1.26.0
	k8s.io/mount-utils v0.26.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2 // indirect
	github.com/moby/sys/mountinfo v0.6.2 // indirect
	github.com/prometheus/client_golang v1.14.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	google.golang.org/genproto v0.0.0-20221118155620-16455021b5e6 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/apimachinery v0.26.0 // indirect
	k8s.io/apiserver v0.0.0 // indirect
	k8s.io/component-base v0.26.0 // indirect
	k8s.io/utils v0.0.0-20221107191617-1a15be271d1d // indirect
)

replace k8s.io/mount-utils => k8s.io/mount-utils v0.26.0

replace k8s.io/apiserver => k8s.io/apiserver v0.26.0
