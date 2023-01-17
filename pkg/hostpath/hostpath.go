package hostpath

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"k8s.io/mount-utils"
)

type hostPath struct {
	config  Config
	mounter mount.Interface
	vcam    []*csi.VolumeCapability_AccessMode
	cscap   []*csi.ControllerServiceCapability
	nscap   []*csi.NodeServiceCapability
}

type Config struct {
	DriverName    string
	Endpoint      string
	NodeID        string
	DataDir       string
	VendorVersion string
}

func NewHostPathDriver(cfg Config) (*hostPath, error) {
	if cfg.DriverName == "" {
		return nil, errors.New("no driver name provided")
	}
	if cfg.Endpoint == "" {
		return nil, errors.New("no node id provided")
	}
	if cfg.NodeID == "" {
		return nil, errors.New("no driver endpoint provided")
	}
	if cfg.DataDir == "" {
		return nil, errors.New("no data dir provided")
	}
	if cfg.VendorVersion == "" {
		return nil, errors.New("no vendor version provided")
	}

	hp := &hostPath{config: cfg, mounter: mount.New("")}
	hp.addVolumeCapabilityAccessModes(
		csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
		csi.VolumeCapability_AccessMode_SINGLE_NODE_MULTI_WRITER,
	)
	hp.addControllerServiceCapabilities(
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
		csi.ControllerServiceCapability_RPC_SINGLE_NODE_MULTI_WRITER,
	)
	hp.addNodeServiceCapabilities(
		csi.NodeServiceCapability_RPC_GET_VOLUME_STATS,
		csi.NodeServiceCapability_RPC_SINGLE_NODE_MULTI_WRITER,
	)

	return hp, nil
}

func NewIdentityServer(hp *hostPath) *IdentityServer {
	return &IdentityServer{hp: hp}
}

func NewControllerServer(hp *hostPath) *ControllerServer {
	return &ControllerServer{hp: hp}
}

func NewNodeServer(hp *hostPath) *NodeServer {
	return &NodeServer{hp: hp}
}

func (hp *hostPath) Run() error {
	ids := NewIdentityServer(hp)
	cs := NewControllerServer(hp)
	ns := NewNodeServer(hp)
	s := NewNonBlockingGRPCServer()

	s.Start(hp.config.Endpoint, ids, cs, ns)
	s.Wait()

	return nil
}

func (hp *hostPath) addVolumeCapabilityAccessModes(ml ...csi.VolumeCapability_AccessMode_Mode) {
	var vcam []*csi.VolumeCapability_AccessMode
	for _, m := range ml {
		vcam = append(vcam, &csi.VolumeCapability_AccessMode{
			Mode: m,
		})
	}
	hp.vcam = vcam
}

func (hp *hostPath) addControllerServiceCapabilities(cl ...csi.ControllerServiceCapability_RPC_Type) {
	var csc []*csi.ControllerServiceCapability
	for _, c := range cl {
		csc = append(csc, &csi.ControllerServiceCapability{
			Type: &csi.ControllerServiceCapability_Rpc{
				Rpc: &csi.ControllerServiceCapability_RPC{
					Type: c,
				},
			},
		})
	}
	hp.cscap = csc
}

func (hp *hostPath) addNodeServiceCapabilities(nl ...csi.NodeServiceCapability_RPC_Type) {
	var nsc []*csi.NodeServiceCapability
	for _, n := range nl {
		nsc = append(nsc, &csi.NodeServiceCapability{
			Type: &csi.NodeServiceCapability_Rpc{
				Rpc: &csi.NodeServiceCapability_RPC{
					Type: n,
				},
			},
		})
	}
	hp.nscap = nsc
}

func (hp *hostPath) getVolumePath(volumeID string) string {
	return filepath.Join(hp.config.DataDir, volumeID)
}

func (hp *hostPath) createVolume(volumeID string) error {
	path := hp.getVolumePath(volumeID)
	return os.MkdirAll(path, 0o755)
}

func (hp *hostPath) deleteVolume(volumeID string) error {
	path := hp.getVolumePath(volumeID)
	return os.RemoveAll(path)
}

func (hp *hostPath) validateVolumeCapabilities(caps []*csi.VolumeCapability) error {
	if caps == nil {
		return errors.New("volume capabilities is nil")
	}
	for _, cap := range caps {
		if err := hp.validateVolumeCapability(cap); err != nil {
			return err
		}
	}
	return nil
}

func (hp *hostPath) validateVolumeCapability(cap *csi.VolumeCapability) error {
	if err := hp.validateAccessMode(cap.GetAccessMode()); err != nil {
		return err
	}
	if cap.GetMount() == nil {
		return errors.New("only mount access type is supported")
	}
	return nil
}

func (hp *hostPath) validateAccessMode(am *csi.VolumeCapability_AccessMode) error {
	if am == nil {
		return errors.New("access mode is nil")
	}
	for _, vac := range hp.vcam {
		if am.GetMode() == vac.GetMode() {
			return nil
		}
	}
	return fmt.Errorf("%s access mode is not supported", am.GetMode())
}
