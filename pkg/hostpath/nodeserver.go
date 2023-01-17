package hostpath

import (
	"context"
	"os"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/volume/util/fs"
)

type NodeServer struct {
	hp *hostPath
}

func (ns *NodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	volumeID := req.GetVolumeId()
	cap := req.GetVolumeCapability()
	source := ns.hp.getVolumePath(volumeID)
	target := req.GetTargetPath()

	if len(volumeID) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}
	if len(target) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Target path missing in request")
	}
	if cap == nil {
		return nil, status.Error(codes.InvalidArgument, "Volume capability missing in request")
	}
	if err := ns.hp.validateVolumeCapability(cap); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid volume capability: %s", err)
	}

	isMnt, err := ns.hp.mounter.IsMountPoint(target)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, status.Errorf(codes.Internal, "Check target path: %q", target)
		}
		if err := os.MkdirAll(target, 0o750); err != nil {
			return nil, status.Errorf(codes.Internal, "Create target path: %q", target)
		}
		isMnt = false
	}

	if isMnt {
		return &csi.NodePublishVolumeResponse{}, nil
	}

	fsType := cap.GetMount().GetFsType()
	mountOptions := []string{"bind"}
	mountOptions = append(mountOptions, cap.GetMount().GetMountFlags()...)
	if req.GetReadonly() {
		mountOptions = append(mountOptions, "ro")
	}

	klog.V(2).Infof("NodePublishVolume: volumeID(%v) source(%s) targetPath(%s) fsType(%s) mountflags(%v)", volumeID, source, target, fsType, mountOptions)
	if err := ns.hp.mounter.Mount(source, target, "", mountOptions); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to mount device: %q at %q: %v", source, target, err)
	}

	return &csi.NodePublishVolumeResponse{}, nil
}

func (ns *NodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	volumeID := req.GetVolumeId()
	target := req.GetTargetPath()

	if len(volumeID) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}
	if len(target) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Target path missing in request")
	}

	if isMnt, err := ns.hp.mounter.IsMountPoint(target); err != nil {
		if !os.IsNotExist(err) {
			return nil, status.Errorf(codes.Internal, "Check target path: %q", target)
		}
	} else if isMnt {
		klog.V(2).Infof("NodeUnpublishVolume: unmounting volume %s on %s", volumeID, target)
		if err := ns.hp.mounter.Unmount(target); err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to unmount target %q: %v", target, err)
		}
		klog.V(2).Infof("NodeUnpublishVolume: unmount volume %s on %s successfully", volumeID, target)
	}

	if err := os.RemoveAll(target); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to remove target %q: %v", target, err)
	}

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

func (ns *NodeServer) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	return &csi.NodeGetInfoResponse{
		NodeId: ns.hp.config.NodeID,
	}, nil
}

func (ns *NodeServer) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: ns.hp.nscap,
	}, nil
}

func (ns *NodeServer) NodeGetVolumeStats(ctx context.Context, req *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	volumeID := req.GetVolumeId()
	volumePath := req.GetVolumePath()

	if len(volumeID) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID not provided")
	}
	if len(volumePath) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume path not provided")
	}

	if _, err := os.Lstat(volumePath); err != nil {
		return nil, status.Errorf(codes.NotFound, "Could not get file information from %q: %+v", volumePath, err)
	}

	available, capacity, usage, inodes, inodesFree, inodesUsed, err := fs.Info(volumePath)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get volume stats %q: %v", volumePath, err)
	}

	return &csi.NodeGetVolumeStatsResponse{
		Usage: []*csi.VolumeUsage{
			{
				Available: available,
				Total:     capacity,
				Used:      usage,
				Unit:      csi.VolumeUsage_BYTES,
			},
			{
				Available: inodesFree,
				Total:     inodes,
				Used:      inodesUsed,
				Unit:      csi.VolumeUsage_INODES,
			},
		},
	}, nil
}

func (ns *NodeServer) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (ns *NodeServer) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (ns *NodeServer) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}
