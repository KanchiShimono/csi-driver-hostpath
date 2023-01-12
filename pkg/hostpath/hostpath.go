package hostpath

import "k8s.io/mount-utils"

type hostPath struct {
	config Config
}

type Config struct {
	DriverName string
	Endpoint   string
	NodeID     string
}

func NewHostPathDriver(cfg Config) *hostPath {
	return &hostPath{config: cfg}
}

func (hp *hostPath) Run() error {
	ids := &IdentityServer{hp: hp}
	cs := &ControllerServer{hp: hp}
	ns := &NodeServer{hp: hp, mounter: mount.New("")}
	s := NewNonBlockingGRPCServer()

	s.Start(hp.config.Endpoint, ids, cs, ns)
	s.Wait()

	return nil
}
