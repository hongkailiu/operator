package ocptf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	// VERSION of the ocptf cmd
	VERSION = "0.0.1"
)

// LoadTFStateFile loads terraform tf state file
func LoadTFStateFile(path string) (*TFState, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	tfState := TFState{}
	err = json.Unmarshal(bytes, &tfState)
	if err != nil {
		return nil, err
	}
	return &tfState, nil
}

// TFState is terraform state file content
type TFState struct {
	Modules []Module
}

// Module is a module in terraform state file
type Module struct {
	Resources map[string]Resource
}

// Resource is a resource in terraform state file
type Resource struct {
	Type    string
	Primary Primary
}

// Primary is a primary in terraform state file
type Primary struct {
	ID         string
	Attributes map[string]string
}

func load(path string) ([]Group, []Host, error) {
	tfState, err := LoadTFStateFile(path)
	if err != nil {
		return nil, nil, err
	}
	if len(tfState.Modules) != 1 {
		return nil, nil, fmt.Errorf("len(tfState.Modules) is %d", len(tfState.Modules))
	}

	osev3Group := Group{Name: "OSEv3", Vars: map[string]interface{}{}, Children: []string{"masters", "nodes", "etcd"}, Hosts: []string{}}
	mastersGroup := Group{Name: "masters", Vars: map[string]interface{}{}, Children: []string{}}
	nodesGroup := Group{Name: "nodes", Vars: map[string]interface{}{}, Children: []string{}}
	etcdGroup := Group{Name: "etcd", Vars: map[string]interface{}{}, Children: []string{}}
	glusterGroup := Group{Name: "glusterfs", Vars: map[string]interface{}{}, Children: []string{}}
	var hosts []Host
	var firstNodeHost *Host
	var crio = false
	if strings.ToLower(os.Getenv("crio")) == "true" {
		crio = true
	}

	for k, v := range tfState.Modules[0].Resources {
		log.WithFields(log.Fields{"k": k, "v": v}).Debug("resource")

		if v.Type != "aws_instance" {
			continue
		}
		if osev3Group.Vars["openshift_cloudprovider_kind"] == nil {
			osev3Group.Vars["openshift_cloudprovider_kind"] = "aws"
		}
		h := Host{Name: v.Primary.Attributes["public_dns"], ID: v.Primary.ID, PublicIP: v.Primary.Attributes["public_ip"]}
		h.VarMap = map[string]interface{}{"openshift_public_hostname": h.Name}

		if !strings.Contains(k, "master") && !strings.Contains(k, "infra") &&
			!strings.Contains(k, "node") && !strings.Contains(k, "worker") &&
			!strings.Contains(k, "etcd") && !strings.Contains(k, "gluster") {
			return nil, nil, fmt.Errorf("malformed instance name,: %s", k)
		}

		if strings.Contains(k, "master") {
			mastersGroup.Hosts = append(mastersGroup.Hosts, h.Name)
			nodesGroup.Hosts = append(nodesGroup.Hosts, h.Name)
			h.VarMap["openshift_node_group_name"] = getOpenshiftNodeGroupName("node-config-master", crio)
			h.VarMap["openshift_schedulable"] = true
		}
		if strings.Contains(k, "infra") {
			nodesGroup.Hosts = append(nodesGroup.Hosts, h.Name)
			h.VarMap["openshift_node_group_name"] = getOpenshiftNodeGroupName("node-config-infra", crio)
			if osev3Group.Vars["openshift_master_default_subdomain"] == nil {
				osev3Group.Vars["openshift_master_default_subdomain"] = fmt.Sprintf("apps.%s.xip.io", h.PublicIP)
			}
		}
		if strings.Contains(k, "worker") || strings.Contains(k, "node") {
			nodesGroup.Hosts = append(nodesGroup.Hosts, h.Name)
			h.VarMap["openshift_node_group_name"] = getOpenshiftNodeGroupName("node-config-compute", crio)
			//h.VarMap["k001"] = 23
			//h.VarMap["k002"] = 23.23
		}
		if strings.Contains(k, "gluster") {
			nodesGroup.Hosts = append(nodesGroup.Hosts, h.Name)
			glusterGroup.Hosts = append(glusterGroup.Hosts, h.Name)
			h.VarMap["openshift_node_group_name"] = getOpenshiftNodeGroupName("node-config-compute", crio)
			h.VarMap["glusterfs_devices"] = `'["/dev/nvme2n1"]'`
		}

		hosts = append(hosts, h)
		if firstNodeHost == nil && len(nodesGroup.Hosts) == 1 {
			firstNodeHost = &h
		}
	}

	if len(nodesGroup.Hosts) == 1 {
		log.WithFields(log.Fields{"firstNodeHost": firstNodeHost}).Debug("===")
		firstNodeHost.VarMap["openshift_node_group_name"] = getOpenshiftNodeGroupName("node-config-all-in-one", crio)
		osev3Group.Vars["openshift_master_default_subdomain"] = fmt.Sprintf("apps.%s.xip.io", firstNodeHost.PublicIP)
	}

	if len(etcdGroup.Hosts) == 0 {
		etcdGroup.Hosts = mastersGroup.Hosts
	}

	log.WithFields(log.Fields{"hosts": hosts}).Debug("hosts")
	if strings.ToLower(os.Getenv("install_ocp_gluster")) != "false" && len(glusterGroup.Hosts) != 0 {
		osev3Group.Children = append(osev3Group.Children, "glusterfs")
	}
	groups := []Group{osev3Group, mastersGroup, nodesGroup, etcdGroup, glusterGroup}
	if strings.ToLower(os.Getenv("install_ocp_gluster")) == "false" || len(glusterGroup.Hosts) == 0 {
		groups = []Group{osev3Group, mastersGroup, nodesGroup, etcdGroup}
	}

	return groups, hosts, nil
}

func getOpenshiftNodeGroupName(name string, crio bool) string {
	if crio {
		return fmt.Sprintf("%s-crio", name)
	}
	return name
}

// DoList answers `--list` flag
func DoList(path string, dynamic bool) error {
	log.WithFields(log.Fields{"path": path, "dynamic": dynamic}).Debug("DoList")
	groups, hosts, err := load(path)
	if err != nil {
		return err
	}
	listOutput, err := GetListOutput(groups, hosts)
	if err != nil {
		return err
	}

	if dynamic {
		bytes, err := json.MarshalIndent(listOutput.GroupMap, "", "  ")
		if err != nil {
			return err
		}
		jsonString := fmt.Sprintf(string(bytes))
		log.WithFields(log.Fields{"jsonString": jsonString}).Debug("")
		fmt.Println(jsonString)

	} else {
		return doListStatic(listOutput)
	}

	return nil
}

func doListStatic(listOutput *ListOutput) error {
	hostVarsMap, err := getHostVarsMap(listOutput)
	if err != nil {
		return fmt.Errorf("error occurred when getHostVarsMap(&listOutput)")
	}
	var b bytes.Buffer
	b.WriteString("###ocptf generated inventory###\n")
	fmt.Println()
	for k, v := range listOutput.GroupMap {
		log.WithFields(log.Fields{"k": k, "v": v}).Debug("DoList: listOutput.GroupMap")
		if k != UnderlineMetaKey {
			b.WriteString(fmt.Sprintf("\n[%s]\n", k))
			g, ok := v.(Group)
			if !ok {
				return fmt.Errorf("wrong format for Group: %+v", v)
			}
			for _, h := range g.Hosts {
				b.WriteString(fmt.Sprintf("%s", h))
				hValue := hostVarsMap[h]
				varMap, ok := hValue.(map[string]interface{})
				if !ok {
					return fmt.Errorf("wrong format for VarMap: %+v", hValue)
				}
				for varK, varV := range varMap {
					log.WithFields(log.Fields{"k": k, "varK": varK, "varV": varV, "reflect.TypeOf(varV)": reflect.TypeOf(varV)}).Debug("DoList: varMap")
					if k == "etcd" || k == "masters" {
						continue
					}
					if k == "nodes" && varK == "glusterfs_devices" {
						continue
					}
					if k == "glusterfs" && varK != "glusterfs_devices" {
						continue
					}
					switch varV.(type) {
					case int:
						b.WriteString(fmt.Sprintf(" %s=%d", varK, varV))
					case float64:
						b.WriteString(fmt.Sprintf(" %s=%s", varK, strconv.FormatFloat(varV.(float64), 'f', -1, 64)))
					case string:
						b.WriteString(fmt.Sprintf(" %s=%s", varK, varV))
					case bool:
						b.WriteString(fmt.Sprintf(" %s=%s", varK, strconv.FormatBool(varV.(bool))))
					default:
						return fmt.Errorf("unknown type for varV: %+v", varV)
					}
				}
				b.WriteString("\n")
			}
			if len(g.Children) != 0 {
				b.WriteString(fmt.Sprintf("\n[%s:children]\n", k))
				for _, c := range g.Children {
					b.WriteString(fmt.Sprintf("%s\n", c))
				}
			}

			if len(g.Vars) != 0 {
				b.WriteString(fmt.Sprintf("\n[%s:vars]\n", k))
				for varK, varV := range g.Vars {
					b.WriteString(fmt.Sprintf("%s=%s\n", varK, varV))
				}
			}

		}
	}
	fmt.Println(fmt.Sprintf("%s\n", b.String()))
	return nil
}

func getHostVarsMap(listOutput *ListOutput) (map[string]interface{}, error) {
	for k, v := range listOutput.GroupMap {
		if k == UnderlineMetaKey {
			HostVars, ok := v.(map[string]map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("wrong format for HostVars: %+v", v)
			}
			return HostVars[HostVarsKey], nil
		}
	}
	return nil, nil
}

// DoHost answers `--host` flag
func DoHost(path, name string, dynamic bool) error {
	_, hosts, err := load(path)
	if err != nil {
		return err
	}
	hostOutput := GetHostOutput(name, hosts)

	if dynamic {
		bytes, err := json.MarshalIndent(hostOutput.VarMap, "", "  ")
		if err != nil {
			return err
		}
		jsonString := fmt.Sprintf(string(bytes))
		fmt.Println(jsonString)

	} else {
		return fmt.Errorf("TODO: static inv output")
	}

	return nil
}
