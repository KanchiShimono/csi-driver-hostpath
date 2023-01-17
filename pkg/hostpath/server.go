package hostpath

import (
	"sync"

	"github.com/KanchiShimono/csi-driver-hostpath/internal/endpoint"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/kubernetes-csi/csi-lib-utils/protosanitizer"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"k8s.io/klog/v2"
)

func NewNonBlockingGRPCServer() *nonBlockingGRPCServer {
	return &nonBlockingGRPCServer{}
}

type nonBlockingGRPCServer struct {
	wg      sync.WaitGroup
	server  *grpc.Server
	cleanup func()
}

func (s *nonBlockingGRPCServer) Start(endpoint string, ids csi.IdentityServer, cs csi.ControllerServer, ns csi.NodeServer) {
	s.wg.Add(1)
	go s.serve(endpoint, ids, cs, ns)
}

func (s *nonBlockingGRPCServer) Wait() {
	s.wg.Wait()
}

func (s *nonBlockingGRPCServer) Stop() {
	s.server.GracefulStop()
	s.cleanup()
}

func (s *nonBlockingGRPCServer) ForceStop() {
	s.server.Stop()
	s.cleanup()
}

func (s *nonBlockingGRPCServer) serve(ep string, ids csi.IdentityServer, cs csi.ControllerServer, ns csi.NodeServer) {
	listener, cleanup, err := endpoint.Listen(ep)
	if err != nil {
		klog.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(logGRPC),
	}
	server := grpc.NewServer(opts...)
	s.server = server
	s.cleanup = cleanup

	if ids != nil {
		csi.RegisterIdentityServer(server, ids)
	}
	if cs != nil {
		csi.RegisterControllerServer(server, cs)
	}
	if ns != nil {
		csi.RegisterNodeServer(server, ns)
	}

	klog.Infof("Listening for connections on address: %#v", listener.Addr())
	server.Serve(listener)
}

func logGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	level := klog.Level(getLogLevel(info.FullMethod))
	klog.V(level).Infof("GRPC call: %s", info.FullMethod)
	klog.V(level).Infof("GRPC request: %s", protosanitizer.StripSecrets(req))

	resp, err := handler(ctx, req)
	if err != nil {
		klog.Errorf("GRPC error: %v", err)
	} else {
		klog.V(level).Infof("GRPC response: %s", protosanitizer.StripSecrets(resp))
	}

	return resp, err
}

func getLogLevel(method string) int32 {
	if method == "/csi.v1.Identity/Probe" ||
		method == "/csi.v1.Node/NodeGetCapabilities" ||
		method == "/csi.v1.Node/NodeGetVolumeStats" {
		return 8
	}
	return 2
}
