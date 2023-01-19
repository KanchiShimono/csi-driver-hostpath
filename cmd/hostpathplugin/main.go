package main

import (
	"flag"

	"github.com/KanchiShimono/csi-driver-hostpath/pkg/hostpath"
	"k8s.io/klog/v2"
)

var version string

func init() {
	flag.Set("logtostderr", "true")
}

func main() {
	cfg := hostpath.Config{VendorVersion: version}
	klog.InitFlags(nil)
	flag.StringVar(&cfg.DriverName, "drivername", "hostpath.csi.kanchi.github.io", "name of the driver")
	flag.StringVar(&cfg.Endpoint, "endpoint", "unix://tmp/csi.sock", "CSI endpoint")
	flag.StringVar(&cfg.NodeID, "nodeid", "", "node id")
	flag.StringVar(&cfg.DataDir, "datadir", "/var/lib/hostpath.csi.kanchi.github.io/data", "data dir")
	flag.Parse()

	driver, err := hostpath.NewHostPathDriver(cfg)
	if err != nil {
		klog.Fatalf("Failed to initialize driver: %s", err)
	}

	if err := driver.Run(); err != nil {
		klog.Fatalf("Failed to run driver: %s", err)
	}
}
