package hostpath

import (
	"context"
	"os"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
)

type ControllerServer struct {
	hp *hostPath
}

func (cs *ControllerServer) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	name := req.GetName()
	caps := req.GetVolumeCapabilities()
	capacity := req.GetCapacityRange().GetRequiredBytes()
	params := req.GetParameters()

	if len(name) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Name missing in request")
	}
	if len(caps) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume capabilities missing in request")
	}
	if err := cs.hp.validateVolumeCapabilities(caps); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid volume capabilities: %s", err)
	}

	uniqID, err := uuid.NewUUID()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate volume ID: %s", err)
	}

	volumeID := uniqID.String()

	klog.V(2).Infof("CreateVolume: volumeID(%v)", volumeID)
	if err := cs.hp.createVolume(volumeID); err != nil && !os.IsExist(err) {
		return nil, status.Errorf(codes.Internal, "Failed to create volume: %s", err)
	}

	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			CapacityBytes: capacity,
			VolumeId:      volumeID,
			VolumeContext: params,
		},
	}, nil
}

func (cs *ControllerServer) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	volumeID := req.GetVolumeId()

	if len(volumeID) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}

	klog.V(2).Infof("DeleteVolume: volumeID(%v)", volumeID)
	if err := cs.hp.deleteVolume(volumeID); err != nil && !os.IsNotExist(err) {
		return nil, status.Errorf(codes.Internal, "Failed to delete volume: %s", err)
	}

	return &csi.DeleteVolumeResponse{}, nil
}

func (cs *ControllerServer) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	volumeID := req.GetVolumeId()
	caps := req.GetVolumeCapabilities()

	if len(volumeID) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}
	if len(caps) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume capabilities missing in request")
	}
	if err := cs.hp.validateVolumeCapabilities(caps); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid volume capabilities: %s", err)
	}

	path := cs.hp.getVolumePath(volumeID)
	if _, err := os.Lstat(path); os.IsNotExist(err) {
		return nil, status.Errorf(codes.NotFound, "Volume is not found %q: %s", path, err)
	}

	return &csi.ValidateVolumeCapabilitiesResponse{
		Confirmed: &csi.ValidateVolumeCapabilitiesResponse_Confirmed{
			VolumeContext:      req.GetVolumeContext(),
			VolumeCapabilities: caps,
			Parameters:         req.GetParameters(),
		},
	}, nil
}

func (cs *ControllerServer) ControllerGetCapabilities(ctx context.Context, req *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	return &csi.ControllerGetCapabilitiesResponse{
		Capabilities: cs.hp.cscap,
	}, nil
}

func (cs *ControllerServer) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (cs *ControllerServer) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (cs *ControllerServer) ListVolumes(ctx context.Context, req *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (cs *ControllerServer) GetCapacity(ctx context.Context, req *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (cs *ControllerServer) CreateSnapshot(ctx context.Context, req *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (cs *ControllerServer) DeleteSnapshot(ctx context.Context, req *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (cs *ControllerServer) ListSnapshots(ctx context.Context, req *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (cs *ControllerServer) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (cs *ControllerServer) ControllerGetVolume(ctx context.Context, req *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}
