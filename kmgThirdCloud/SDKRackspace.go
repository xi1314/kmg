package kmgThirdCloud

import (
	"fmt"
	"github.com/rackspace/gophercloud"
	OpenStackImages "github.com/rackspace/gophercloud/openstack/compute/v2/images"
	OpenStackServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/flavors"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
	"time"
)

type RackspaceInstanceStatus string

const (
	RackspaceInstanceStatusACTIVE  RackspaceInstanceStatus = "ACTIVE"
	RackspaceInstanceStatusBUILD   RackspaceInstanceStatus = "BUILD"
	RackspaceInstanceStatusDELETED RackspaceInstanceStatus = "DELETED"
	RackspaceInstanceStatusERROR   RackspaceInstanceStatus = "ERROR"
	RackspaceInstanceStatusUNKNOWN RackspaceInstanceStatus = "UNKNOWN"
)

type RackspaceSDK struct {
	Username     string
	APIKey       string
	p            *gophercloud.ServiceClient
	SSHKeyName   string
	Region       string
	InstanceName string
	ImageName    string
	FlavorName   string
}

func NewRackspaceSDK(username, apiKey, SSHKeyName string) *RackspaceSDK {
	if SSHKeyName == "" {
		panic("Empty SSHKeyName of instance don't allow!")
	}
	sdk := &RackspaceSDK{
		Username:     username,
		APIKey:       apiKey,
		SSHKeyName:   SSHKeyName,
		Region:       "HKG",
		InstanceName: "kmg-RackspaceSDK-auto-build",
		ImageName:    "Ubuntu 14.04 LTS (Trusty Tahr) (PVHVM)",
		FlavorName:   "4 GB General Purpose v1",
	}
	ao := gophercloud.AuthOptions{
		Username: sdk.Username,
		APIKey:   sdk.APIKey,
	}
	provider, err := rackspace.AuthenticatedClient(ao)
	handleErr(err)
	serviceClient, err := rackspace.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: sdk.Region,
	})
	handleErr(err)
	sdk.p = serviceClient
	return sdk
}

func (sdk *RackspaceSDK) CreateInstance() (ip string) {
	opt := servers.CreateOpts{
		Name:       sdk.InstanceName,
		ImageName:  sdk.ImageName,
		FlavorName: sdk.FlavorName,
		KeyPair:    sdk.SSHKeyName,
	}
	result := servers.Create(sdk.p, opt)
	one, err := result.Extract()
	handleErr(err)
	for {
		rackspaceInstance, err := servers.Get(sdk.p, one.ID).Extract()
		handleErr(err)
		if rackspaceInstance.Status == string(RackspaceInstanceStatusACTIVE) {
			one = rackspaceInstance
			break
		}
		if rackspaceInstance.Status == string(RackspaceInstanceStatusBUILD) {
			time.Sleep(time.Second)
		}
		if rackspaceInstance.Status == string(RackspaceInstanceStatusUNKNOWN) || rackspaceInstance.Status == string(RackspaceInstanceStatusERROR) {
			panic("Rackspace CreateInstance got ERROR/UNKNOWN " + rackspaceInstance.AccessIPv4 + " " + rackspaceInstance.ID)
		}
	}
	return one.AccessIPv4
}

func (sdk *RackspaceSDK) RenameInstanceByIp(name, ip string) {
	instance, exist := sdk.ListAllInstance()[ip]
	if exist {
		servers.Update(sdk.p, instance.Id, OpenStackServers.UpdateOpts{
			Name: name,
		})
	}
}

func (sdk *RackspaceSDK) DeleteInstance(ip string) {
	instance, exist := sdk.ListAllInstance()[ip]
	if exist {
		servers.Delete(sdk.p, instance.Id)
	}
	for {
		_, exist := sdk.ListAllInstance()[ip]
		if !exist {
			break
		}
		time.Sleep(time.Second)
	}
}

func (sdk *RackspaceSDK) ListAllInstance() (ipInstanceMap map[string]Instance) {
	ipInstanceMap = map[string]Instance{}
	err := servers.List(sdk.p, nil).EachPage(func(p pagination.Page) (quit bool, err error) {
		serverSlice, err := servers.ExtractServers(p)
		for _, server := range serverSlice {
			if server.Status != string(RackspaceInstanceStatusACTIVE) {
				fmt.Println(server.AccessIPv4, server.Status)
				continue
			}
			ipInstanceMap[server.AccessIPv4] = Instance{
				Ip:   server.AccessIPv4,
				Id:   server.ID,
				Name: server.Name,
			}
		}
		return true, err
	})
	handleErr(err)
	return
}

func (sdk *RackspaceSDK) GetAllFlavor() []flavors.Flavor {
	out := []flavors.Flavor{}
	err := flavors.ListDetail(sdk.p, nil).EachPage(func(p pagination.Page) (quit bool, err error) {
		v, err := flavors.ExtractFlavors(p)
		out = append(out, v...)
		return true, err
	})
	handleErr(err)
	return out
}

func (sdk *RackspaceSDK) GetAllImage() []OpenStackImages.Image {
	out := []OpenStackImages.Image{}
	err := OpenStackImages.ListDetail(sdk.p, nil).EachPage(func(p pagination.Page) (quit bool, err error) {
		v, err := OpenStackImages.ExtractImages(p)
		out = append(out, v...)
		return true, err
	})
	handleErr(err)
	return out
}

func (sdk *RackspaceSDK) PrintAllImage() {
	all := sdk.GetAllImage()
	for _, i := range all {
		fmt.Println("ID", i.ID, "Name", i.Name)
	}
}

func (sdk *RackspaceSDK) PrintAllFlavor() {
	all := sdk.GetAllFlavor()
	for _, i := range all {
		fmt.Println("ID", i.ID, "Name", i.Name)
	}
}
