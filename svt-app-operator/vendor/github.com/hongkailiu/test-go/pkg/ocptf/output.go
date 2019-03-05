package ocptf

import (
	"fmt"
)

const (
	// UnderlineMetaKey is underline meta key, ie, _meta
	UnderlineMetaKey = "_meta"
	// HostVarsKey is host vars key, ie, hostvars
	HostVarsKey = "hostvars"
)

// ListOutput represents output for listing inventory content
type ListOutput struct {
	GroupMap map[string]interface{}
}

// HostOutput represents output for listing inventory content for a host
type HostOutput struct {
	VarMap map[string]interface{}
}

// Group represent group
type Group struct {
	Name     string                 `json:"-"`
	Hosts    []string               `json:"hosts"`
	Vars     map[string]interface{} `json:"vars"`
	Children []string               `json:"children"`
}

// Host represents host
type Host struct {
	Name     string
	ID       string `json:"-"`
	PublicIP string `json:"-"`
	VarMap   map[string]interface{}
}

// GetListOutput gets a ListOutput
func GetListOutput(groups []Group, hosts []Host) (*ListOutput, error) {
	groupMap := make(map[string]interface{})
	for _, g := range groups {
		if g.Name == UnderlineMetaKey {
			return nil, fmt.Errorf("key (%s) is reserved", UnderlineMetaKey)
		}
		if groupMap[g.Name] != nil {
			return nil, fmt.Errorf("key (%s) exists already in groupMap with value (%v)", g.Name, groupMap[g.Name])
		}
		groupMap[g.Name] = g
	}
	HostVars := make(map[string]map[string]interface{})
	hostVarsMap := make(map[string]interface{})
	for _, h := range hosts {
		if hostVarsMap[h.Name] != nil {
			return nil, fmt.Errorf("key (%s) exists already in hostVarsMap with value (%v)", h.Name, hostVarsMap[h.Name])
		}
		hostVarsMap[h.Name] = h.VarMap
	}
	HostVars[HostVarsKey] = hostVarsMap
	groupMap[UnderlineMetaKey] = HostVars
	listOutput := ListOutput{groupMap}
	return &listOutput, nil

}

// GetHostOutput gets a HostOutput
func GetHostOutput(name string, hosts []Host) *HostOutput {
	for _, h := range hosts {
		if name == h.Name {
			return &HostOutput{h.VarMap}
		}
	}
	return &HostOutput{map[string]interface{}{}}

}
