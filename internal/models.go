package internal

import (
	"strings"
)

type ParsedConfig struct {
	Description string `json:"description,omitempty"`
}

type ParsedAgentInterfaces struct {
	Result []struct {
		IPAddresses []IP `json:"ip-addresses"`
	} `json:"result"`
}

type Node struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"node,omitempty"`
	Status string `json:"status,omitempty"`
}

type NodeStatus struct {
	Node string `json:"node"`
}

type VirtualMachine struct {
	VMID   uint64 `json:"vmid"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Container struct {
	VMID   uint64 `json:"vmid"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Version struct {
	Release string `json:"release"`
}

type Service struct {
	ID     uint64
	Name   string
	IPs    []IP
	Config map[string]string
}

type IP struct {
	Address     string `json:"ip-address,omitempty"`
	AddressType string `json:"ip-address-type,omitempty"`
	Prefix      uint64 `json:"prefix,omitempty"`
}

func NewService(id uint64, name string, config map[string]string) Service {
	return Service{ID: id, Name: name, Config: config, IPs: make([]IP, 0)}
}

func (pc *ParsedConfig) GetTraefikMap() map[string]string {
	m := make(map[string]string)
	lines := strings.Split(pc.Description, "\n")
	for _, l := range lines {
		var k, v string
		if strings.Contains(l, ":") {
			parts := strings.SplitN(l, ":", 2)
			if len(parts) > 1 {
				k = strings.Replace(parts[0], "\"", "", -1)
				v = strings.Replace(parts[1], "\"", "", -1)
			}
		} else if strings.Contains(l, "=") {
			parts := strings.SplitN(l, "=", 2)
			if len(parts) > 1 {
				k = strings.Replace(parts[0], "\"", "", -1)
				v = strings.Replace(parts[1], "\"", "", -1)
			}
		}
		
		if k != "" && strings.HasPrefix(k, "traefik") {
			m[strings.TrimSpace(k)] = strings.TrimSpace(v)
		}
	}
	return m
}

func (pai *ParsedAgentInterfaces) GetIPs() []IP {
	ips := make([]IP, 0)
	for _, r := range pai.Result {
		ips = append(ips, r.IPAddresses...)
	}
	return ips
}
