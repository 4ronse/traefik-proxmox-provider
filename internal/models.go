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
	// Early return for empty description
	if pc.Description == "" {
		return make(map[string]string)
	}

	m := make(map[string]string)
	lines := strings.Split(pc.Description, "\n")

	// Precompile the regex for better performance
	urlRegex := regexp.MustCompile(`^(traefik\.http\.services\.[^=]+\.loadbalancer\.server\.url)=(.*)$`)

	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue // Skip empty lines
		}
		
		// Special handling for traefik URL
		if strings.Contains(l, "traefik.http.services") && strings.Contains(l, "loadbalancer.server.url") {
			matches := urlRegex.FindStringSubmatch(l)
			if len(matches) == 3 {
				key := matches[1]
				value := matches[2]
				
				// Fix HTTP URL format with more comprehensive replacements
				value = strings.ReplaceAll(value, "http = //", "http://")
				value = strings.ReplaceAll(value, "https = //", "https://")
				value = strings.ReplaceAll(value, " = ", ":")
				
				m[key] = value
				continue
			}
		}
		
		// Extract key-value using a more robust approach
		var k, v string
		
		// Check for "=" first as it appears to be the primary delimiter
		if idx := strings.Index(l, "="); idx > 0 {
			k = strings.TrimSpace(strings.Trim(l[:idx], "\""))
			v = strings.TrimSpace(strings.Trim(l[idx+1:], "\""))
		} else if idx := strings.Index(l, ":"); idx > 0 {
			// Fall back to ":" delimiter
			k = strings.TrimSpace(strings.Trim(l[:idx], "\""))
			v = strings.TrimSpace(strings.Trim(l[idx+1:], "\""))
		}
		
		// Store only traefik-related key-value pairs
		if k != "" && strings.HasPrefix(k, "traefik") {
			m[k] = v
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
