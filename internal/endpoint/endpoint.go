package endpoint

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func Parse(ep string) (string, string, error) {
	if strings.HasPrefix(strings.ToLower(ep), "unix://") || strings.HasPrefix(strings.ToLower(ep), "tcp://") {
		s := strings.SplitN(ep, "://", 2)
		if s[1] != "" {
			return s[0], s[1], nil
		}
		return "", "", fmt.Errorf("Invalid endpint: %v", ep)
	}
	return "unix", ep, nil
}

func Listen(ep string) (net.Listener, func(), error) {
	proto, addr, err := Parse(ep)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {}

	if proto == "unix" {
		addr = "/" + addr
		if err := os.Remove(addr); err != nil && !os.IsNotExist(err) {
			return nil, nil, fmt.Errorf("%s: %q", addr, err)
		}

		cleanup = func() {
			os.Remove(addr)
		}
	}

	l, err := net.Listen(proto, addr)
	return l, cleanup, err
}
