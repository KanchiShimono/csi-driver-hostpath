package hostpath

import (
	"sync"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc"
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

func (s *nonBlockingGRPCServer) serve(endpoint string, ids csi.IdentityServer, cs csi.ControllerServer, ns csi.NodeServer) {
}
