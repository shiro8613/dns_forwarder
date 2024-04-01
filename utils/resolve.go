package utils

import "net"

type Resolve struct {
	v4 net.IPAddr
	v6 net.IPAddr
}

func ResolveIP(name string) (*Resolve, error) {
	resolve := new(Resolve)

	addr4, err := net.ResolveIPAddr("ip4", name)
	if err != nil {
		return nil, err
	}

	resolve.v4 = *addr4

	addr6, err := net.ResolveIPAddr("ip6", name)
	if err != nil {
		return nil, err
	}

	resolve.v6 = *addr6

	return resolve, nil
}

func (r *Resolve) GetIPv4() net.IPAddr {
	return r.v4
}

func (r *Resolve) GetIPv6() net.IPAddr {
	return r.v6
}
