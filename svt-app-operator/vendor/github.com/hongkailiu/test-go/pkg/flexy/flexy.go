package flexy

import (
	"errors"
	"fmt"
	"time"
)

const (
	// InstanceCountLimit limits the instances created by flexy
	InstanceCountLimit = 23
	// RoleCountLimit limits the instances created by flexy for each roles
	RoleCountLimit = 6
)

// Start starts flexy job
func Start(cp CloudProvider, config OCPClusterConfig, inputPath string, outputFolder string) error {
	if len(config.OCPRoles) > RoleCountLimit {
		return fmt.Errorf("RoleCountLimit is %d: too many roles: %d", RoleCountLimit, len(config.OCPRoles))
	}

	var masterGroup, etcdGroup, nodeGroup, lbGroup, glusterFSGroup []Host
	ocVars := map[string]string{}
	for k, v := range config.OpenshiftAnsibleVar {
		ocVars[k] = v
	}

	switch config.CloudProvider {
	case CloudProviderAWS:
		ocVars["openshift_clusterid"] = config.KubernetesClusterValue
	case CloudProviderDryRunner:
		ocVars["openshift_clusterid"] = config.KubernetesClusterValue
		ocVars["openshift_hosted_registry_storage_gcs_keyfile"] = config.GCSKeyfile
	case CloudProviderGCE:
		ocVars["openshift_hosted_registry_storage_gcs_keyfile"] = config.GCSKeyfile
	}

	instanceCount := 0
	for _, role := range config.OCPRoles {
		if err := ValidateOCPRole(role.Name); err != nil {
			return err
		}
		if instanceCount > InstanceCountLimit {
			return fmt.Errorf("instances over %d", InstanceCountLimit)
		}

		if role.Size > 0 {
			for i := 1; i <= role.Size; i++ {
				name := fmt.Sprintf("%s-%s-%d", config.InstancePrefix, role.Name, i)
				if config.AllInOne {
					name = fmt.Sprintf("%s-%s", config.InstancePrefix, "all-in-one")
					if instanceCount == 1 {
						return errors.New("required more than 1 instance for all-in-one cluster")
					}
				}
				host := Host{}
				configParams := map[string]string{
					"name":                   name,
					"imageID":                config.ImageID,
					"kubernetesClusterValue": config.KubernetesClusterValue,
					"publicKeyFile":          config.PublicKeyFile,
				}
				err := cp.CreateAnInstance(role, configParams, &host)
				if err != nil {
					return err
				}
				instanceCount++

				err = cp.WaitUntilRunning(&host, 2*time.Minute)
				if err != nil {
					return err
				}

				host.OCMasterSchedulable = false
				switch role.Name {
				case OCPRoleMaster:
					host.OCNodeGroupName = "node-config-master"
					host.OCMasterSchedulable = true
				case OCPRoleInfra:
					host.OCNodeGroupName = "node-config-infra"
				case OCPRoleCompute:
					host.OCNodeGroupName = "node-config-compute"
				case OCPRoleGlusterFS:
					host.OCNodeGroupName = "node-config-compute"
				}

				if config.AllInOne {
					host.OCNodeGroupName = "node-config-all-in-one"
					ocVars["openshift_master_default_subdomain"] = fmt.Sprintf("apps.%s.xip.io", host.IPv4PublicIP)
				}
				if role.Name == OCPRoleInfra && len(ocVars["openshift_master_default_subdomain"]) == 0 {
					//"openshift_master_default_subdomain": "apps.someip.xip.io",
					ocVars["openshift_master_default_subdomain"] = fmt.Sprintf("apps.%s.xip.io", host.IPv4PublicIP)
				}

				switch role.Name {
				case OCPRoleMaster:
					masterGroup = append(masterGroup, host)
					nodeGroup = append(nodeGroup, host)
				case OCPRoleInfra:
					nodeGroup = append(nodeGroup, host)
				case OCPRoleCompute:
					nodeGroup = append(nodeGroup, host)
				case OCPRoleGlusterFS:
					nodeGroup = append(nodeGroup, host)
					glusterFSGroup = append(glusterFSGroup, host)
				}

			}

		}

	}
	if len(etcdGroup) == 0 {
		etcdGroup = masterGroup
	}
	inventoryConfig := Config{masterGroup, etcdGroup, nodeGroup, lbGroup, glusterFSGroup, ocVars}
	err := Generate(inputPath, inventoryConfig, outputFolder)
	if err != nil {
		return err
	}
	return nil
}
