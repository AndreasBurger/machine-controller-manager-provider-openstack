package v1alpha1

import (
	"slices"
	unsafe "unsafe"

	conversion "k8s.io/apimachinery/pkg/conversion"

	openstack "github.com/gardener/machine-controller-manager-provider-openstack/pkg/apis/openstack"
)

func idOrName(id string, name string) string {
	if id != "" {
		return id
	}
	return name
}

func Convert_v1alpha1_MachineProviderConfigSpec_To_openstack_MachineProviderConfigSpec(in *MachineProviderConfigSpec, out *openstack.MachineProviderConfigSpec, s conversion.Scope) error {
	out.ImageID = in.ImageID
	out.ImageName = in.ImageName
	out.Region = in.Region
	out.AvailabilityZone = in.AvailabilityZone
	out.FlavorName = in.FlavorName
	out.KeyName = in.KeyName
	out.SecurityGroups = *(*[]string)(unsafe.Pointer(&in.SecurityGroups))
	out.Tags = *(*map[string]string)(unsafe.Pointer(&in.Tags))

	out.RootDiskSize = in.RootDiskSize
	out.RootDiskType = (*string)(unsafe.Pointer(in.RootDiskType))
	out.UseConfigDrive = (*bool)(unsafe.Pointer(in.UseConfigDrive))
	out.ServerGroupID = (*string)(unsafe.Pointer(in.ServerGroupID))

	pod_network_cidr := in.PodNetworkCidr
	subnet := ""
	if in.SubnetID != nil {
		subnet = *in.SubnetID
	}

	// Either NetworkID or Networks is set, _not_ both (see validation)
	// TODO: Is this still validated beforehand? If not, can it be?
	if in.NetworkID != "" {
		out.Network = openstack.OpenStackNetwork{
			NetworkID: in.NetworkID,
			SubnetID:  subnet,
			Cidr:      pod_network_cidr,
		}

	} else {
		// Networks ought to be set
		if nc := len(in.Networks); nc == 0 {
			out.Network = openstack.OpenStackNetwork{}
		} else if nc == 1 {
			// Only one network. This has to be the main network.
			n := in.Networks[0]
			id := idOrName(n.Id, n.Name)

			out.Network = openstack.OpenStackNetwork{
				NetworkID: id,
				SubnetID:  subnet,
			}
		} else {
			// Multiple networks. Try to determine main one (via PodNetwork attr), otherwise
			// treat first one encountered as "main".
			out.AdditionalNetworks = make([]openstack.OpenStackNetwork, nc-1)
			has_pod_network := slices.ContainsFunc(in.Networks, func(n OpenStackNetwork) bool { return n.PodNetwork })
			for i, n := range in.Networks {
				id := idOrName(n.Id, n.Name)

				nw := openstack.OpenStackNetwork{
					NetworkID: id,
					SubnetID:  subnet,
				}

				if n.PodNetwork {
					nw.Cidr = pod_network_cidr
				}

				if has_pod_network && n.PodNetwork || !has_pod_network && i == 0 {
					nw.Cidr = pod_network_cidr // TODO: redundant in most cases.
					out.Network = nw
				} else {
					out.AdditionalNetworks = append(out.AdditionalNetworks, nw)
				}
			}
		}
	}
	return nil
}

func Convert_openstack_MachineProviderConfigSpec_To_v1alpha1_MachineProviderConfigSpec(in *openstack.MachineProviderConfigSpec, out *MachineProviderConfigSpec, s conversion.Scope) error {
	out.ImageID = in.ImageID
	out.ImageName = in.ImageName
	out.Region = in.Region
	out.AvailabilityZone = in.AvailabilityZone
	out.FlavorName = in.FlavorName
	out.KeyName = in.KeyName
	out.SecurityGroups = *(*[]string)(unsafe.Pointer(&in.SecurityGroups))
	out.Tags = *(*map[string]string)(unsafe.Pointer(&in.Tags))
	out.PodNetworkCidr = in.Network.Cidr
	out.RootDiskSize = in.RootDiskSize
	out.RootDiskType = (*string)(unsafe.Pointer(in.RootDiskType))
	out.UseConfigDrive = (*bool)(unsafe.Pointer(in.UseConfigDrive))
	out.ServerGroupID = (*string)(unsafe.Pointer(in.ServerGroupID))

	if in.Network.SubnetID != "" {
		out.SubnetID = &in.Network.SubnetID
	} else {
		out.SubnetID = nil
	}

	if nc := len(in.AdditionalNetworks); nc == 0 {
		out.Networks = nil
		out.NetworkID = in.Network.NetworkID
	} else {
		out.Networks = make([]OpenStackNetwork, nc+1)
		// TODO: PodNetwork ??
		out.Networks = append(out.Networks, OpenStackNetwork{Id: in.Network.NetworkID, Name: "", PodNetwork: true})
		for _, n := range in.AdditionalNetworks {
			//  TODO: Strictly speaking, this could cause issues for other users of the API as
			// they need to be able to handle a name being presented as ID. Check.
			out.Networks = append(out.Networks, OpenStackNetwork{Id: n.NetworkID, Name: "", PodNetwork: false})
		}
	}

	return nil
}

func Convert_openstack_OpenStackNetwork_To_v1alpha1_OpenStackNetwork(in *openstack.OpenStackNetwork, out *OpenStackNetwork, s conversion.Scope) error {
	return nil
}

func Convert_v1alpha1_OpenStackNetwork_To_openstack_OpenStackNetwork(in *OpenStackNetwork, out *openstack.OpenStackNetwork, s conversion.Scope) error {
	return nil
}
