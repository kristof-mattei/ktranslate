package meraki

import (
	"context"
	"strings"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"

	apiclient "github.com/ddexterpark/dashboard-api-golang/client"
	"github.com/ddexterpark/dashboard-api-golang/client/devices"
	"github.com/ddexterpark/dashboard-api-golang/client/networks"
	"github.com/ddexterpark/dashboard-api-golang/client/organizations"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigFastest

type MerakiClient struct {
	log      logger.ContextL
	jchfChan chan []*kt.JCHF
	conf     *kt.SnmpDeviceConfig
	metrics  *kt.SnmpDeviceMetric
	client   *apiclient.MerakiAPIGolang
	auth     runtime.ClientAuthInfoWriter
	orgs     []*organizations.GetOrganizationsOKBodyItems0
	networks []networkDesc
}

type networkDesc struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	org  *organizations.GetOrganizationsOKBodyItems0
}

func NewMerakiClient(jchfChan chan []*kt.JCHF, conf *kt.SnmpDeviceConfig, metrics *kt.SnmpDeviceMetric, log logger.ContextL) (*MerakiClient, error) {
	c := MerakiClient{
		log:      log,
		jchfChan: jchfChan,
		conf:     conf,
		metrics:  metrics,
		orgs:     []*organizations.GetOrganizationsOKBodyItems0{},
		networks: []networkDesc{},
		auth:     httptransport.APIKeyAuth("X-Cisco-Meraki-API-Key", "header", conf.Ext.MerakiConfig.ApiKey),
	}

	host := conf.Ext.MerakiConfig.Host
	if host == "" {
		host = apiclient.DefaultHost
	}

	trans := apiclient.DefaultTransportConfig().WithHost(host)
	client := apiclient.NewHTTPClientWithConfig(nil, trans)
	c.log.Infof("Verifying %s connectivity", c.GetName())

	// First, list out all of the organizations present.
	params := organizations.NewGetOrganizationsParams()
	prod, err := client.Organizations.GetOrganizations(params, c.auth)
	if err != nil {
		return nil, err
	}

	for _, org := range prod.GetPayload() {
		lorg := org
		// Now list the networks for this org.
		params := organizations.NewGetOrganizationNetworksParams()
		params.SetOrganizationID(org.ID)
		prod, err := client.Organizations.GetOrganizationNetworks(params, c.auth)
		if err != nil {
			return nil, err
		}

		b, err := json.Marshal(prod.GetPayload())
		if err != nil {
			return nil, err
		}
		var networks []networkDesc
		err = json.Unmarshal(b, &networks)
		if err != nil {
			return nil, err
		}

		for _, network := range networks {
			network.org = lorg
			c.log.Infof("Adding network %s to list to track", network.Name)
			c.networks = append(c.networks, network)
		}

		if len(networks) > 0 { // Only add this org in to track if it has some networks.
			c.log.Infof("Adding organization %s to list to track", org.Name)
			c.orgs = append(c.orgs, lorg)
		}
	}

	// If we get this far, we have a list of things to look at.
	c.log.Infof("%s connected to API for %s with %d organization(s) and %d network(s).", c.GetName(), conf.DeviceName, len(c.orgs), len(c.networks))
	c.client = client

	return &c, nil
}

func (c *MerakiClient) GetName() string {
	return "Meraki API"
}

func (c *MerakiClient) Run(ctx context.Context, dur time.Duration) {
	poll := time.NewTicker(dur)
	defer poll.Stop()
	c.log.Infof("Running Every %v", dur)

	for {
		select {

		// Track the counters here, to convert from raw counters to differences
		case _ = <-poll.C:
			/* // Don't call this now.
				if res, err := c.getNetworkClients(); err != nil {
				c.log.Infof("Meraki cannot get Network Client Info: %v", err)
			} else if len(res) > 0 {
				c.jchfChan <- res
			}
			*/

			if res, err := c.getDeviceClients(dur); err != nil {
				c.log.Infof("Meraki cannot get Device Client Info: %v", err)
			} else if len(res) > 0 {
				c.jchfChan <- res
			}

		case <-ctx.Done():
			c.log.Infof("Meraki Poll Done")
			return
		}
	}
}

type client struct {
	Usage            map[string]float64 `json:"usage"`
	ID               string             `json:"id"`
	Description      string             `json:"description"`
	Mac              string             `json:"mac"`
	IP               string             `json:"ip"`
	User             string             `json:"user"`
	Vlan             int                `json:"vlan"`
	NamedVlan        string             `json:"namedVlan"`
	IPv6             string             `json:"ip6"`
	Manufacturer     string             `json:"manufacturer"`
	DeviceType       string             `json:"deviceTypePrediction"`
	RecentDeviceName string             `json:"recentDeviceName"`
	Status           string             `json:"status"`
	MdnsName         string             `json:"mdnsName"`
	DhcpHostname     string             `json:"dhcpHostname"`
	network          string
	device           networkDevice
}

func (c *MerakiClient) getNetworkClients() ([]*kt.JCHF, error) {
	clientSet := []client{}
	for _, network := range c.networks {
		params := networks.NewGetNetworkClientsParams()
		params.SetNetworkID(network.ID)

		prod, err := c.client.Networks.GetNetworkClients(params, c.auth)
		if err != nil {
			return nil, err
		}

		b, err := json.Marshal(prod.GetPayload())
		if err != nil {
			return nil, err
		}

		var clients []client
		err = json.Unmarshal(b, &clients)
		if err != nil {
			return nil, err
		}
		for _, client := range clients {
			client.network = network.Name
			clientSet = append(clientSet, client) // Toss these all in together
		}
	}

	return c.parseClients(clientSet)
}

type networkDevice struct {
	Name      string   `json:"name"`
	Serial    string   `json:"serial"`
	Mac       string   `json:"mac"`
	Model     string   `json:"model"`
	Notes     string   `json:"notes"`
	LanIP     string   `json:"lanIp"`
	Tags      []string `json:"tags"`
	NetworkID string   `json:"networkId"`
	Firmware  string   `json:"firmware"`
	network   networkDesc
}

func (c *MerakiClient) getNetworkDevices() (map[string][]networkDevice, error) {
	deviceSet := map[string][]networkDevice{}
	for _, network := range c.networks {
		params := networks.NewGetNetworkDevicesParams()
		params.SetNetworkID(network.ID)

		prod, err := c.client.Networks.GetNetworkDevices(params, c.auth)
		if err != nil {
			return nil, err
		}

		b, err := json.Marshal(prod.GetPayload())
		if err != nil {
			return nil, err
		}

		var devices []networkDevice
		err = json.Unmarshal(b, &devices)
		if err != nil {
			return nil, err
		}

		for i, _ := range devices {
			if devices[i].Name == "" {
				devices[i].Name = devices[i].Serial
			}
			devices[i].network = network
		}

		if len(devices) > 0 {
			deviceSet[network.Name] = devices
		}
	}

	return deviceSet, nil
}

func (c *MerakiClient) getDeviceClients(dur time.Duration) ([]*kt.JCHF, error) {
	networkDevs, err := c.getNetworkDevices()
	if err != nil {
		return nil, err
	}

	c.log.Infof("Got devices for %d networks", len(networkDevs))
	clientSet := []client{}
	durs := float32(3600) // Look back 1 hour in seconds, get all devices using APs in this range.
	for network, deviceSet := range networkDevs {
		c.log.Infof("Looking at %d devices for network %s", len(deviceSet), network)
		for _, device := range deviceSet {
			params := devices.NewGetDeviceClientsParams()
			params.SetSerial(device.Serial)
			params.SetTimespan(&durs)

			prod, err := c.client.Devices.GetDeviceClients(params, c.auth)
			if err != nil {
				return nil, err
			}

			b, err := json.Marshal(prod.GetPayload())
			if err != nil {
				return nil, err
			}

			var clients []client
			err = json.Unmarshal(b, &clients)
			if err != nil {
				return nil, err
			}
			for _, client := range clients {
				client.network = network
				client.device = device
				clientSet = append(clientSet, client) // Toss these all in together
			}
		}
	}

	return c.parseClients(clientSet)
}

func (c *MerakiClient) parseClients(cs []client) ([]*kt.JCHF, error) {
	res := make([]*kt.JCHF, 0)
	for _, client := range cs {
		dst := kt.NewJCHF()
		if client.IPv6 != "" {
			dst.DstAddr = client.IPv6
		} else {
			dst.DstAddr = client.IP
		}
		dst.CustomStr = map[string]string{
			"network":            client.network,
			"client_id":          client.ID,
			"description":        client.Description,
			"status":             client.Status,
			"vlan_name":          client.NamedVlan,
			"client_mac_addr":    client.Mac,
			"user":               client.User,
			"manufacturer":       client.Manufacturer,
			"device_type":        client.DeviceType,
			"recent_device_name": client.RecentDeviceName,
			"dhcp_hostname":      client.DhcpHostname,
			"mdsn_name":          client.MdnsName,
		}
		dst.CustomInt = map[string]int32{
			"vlan": int32(client.Vlan),
		}
		dst.CustomBigInt = map[string]int64{}
		dst.EventType = kt.KENTIK_EVENT_SNMP_DEV_METRIC
		//dst.Provider = c.conf.Provider // @TODO, pick a provider for this one.

		if client.device.Serial != "" {
			dst.DeviceName = client.device.Name // Here, device is this device's name.
			dst.SrcAddr = client.device.LanIP
			dst.CustomStr["device_serial"] = client.device.Serial
			dst.CustomStr["device_firmware"] = client.device.Firmware
			dst.CustomStr["device_mac_addr"] = client.device.Mac
			dst.CustomStr["device_tags"] = strings.Join(client.device.Tags, ",")
			dst.CustomStr["device_notes"] = client.device.Notes
			dst.CustomStr["device_model"] = client.device.Model
			dst.CustomStr["src_ip"] = client.device.LanIP
			if client.device.network.org != nil {
				dst.CustomStr["org_name"] = client.device.network.org.Name
				dst.CustomStr["org_id"] = client.device.network.org.ID
			}
		} else {
			dst.DeviceName = client.network // Here, device is this network's name.
			dst.SrcAddr = c.conf.DeviceIP
		}

		dst.Timestamp = time.Now().Unix()
		dst.CustomMetrics = map[string]kt.MetricInfo{}

		dst.CustomBigInt["Sent"] = int64(client.Usage["sent"] * 1000) // Unit is kilobytes, convert to bytes
		dst.CustomMetrics["Sent"] = kt.MetricInfo{Oid: "meraki", Mib: "meraki", Profile: "meraki.clients", Type: "meraki.clients"}

		dst.CustomBigInt["Recv"] = int64(client.Usage["recv"] * 1000) // Same, convert to bytes.
		dst.CustomMetrics["Recv"] = kt.MetricInfo{Oid: "meraki", Mib: "meraki", Profile: "meraki.clients", Type: "meraki.clients"}

		c.conf.SetUserTags(dst.CustomStr)
		res = append(res, dst)
	}

	return res, nil
}
