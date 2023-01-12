package main

import (
	"flag"

	"github.com/KanchiShimono/csi-driver-hostpath/pkg/hostpath"
	"k8s.io/klog/v2"
)

func init() {
	flag.Set("logtostderr", "true")
}

func main() {
	cfg := hostpath.Config{}

	flag.StringVar(&cfg.DriverName, "drivername", "hostpath.csi.kanchi.github.io", "name of the driver")
	flag.StringVar(&cfg.Endpoint, "endpoint", "unix://tmp/csi.sock", "CSI endpoint")
	flag.StringVar(&cfg.NodeID, "nodeid", "", "node id")

	driver := hostpath.NewHostPathDriver(cfg)

	if err := driver.Run(); err != nil {
		klog.Fatalf("failed to run driver: %s", err)
	}
}
