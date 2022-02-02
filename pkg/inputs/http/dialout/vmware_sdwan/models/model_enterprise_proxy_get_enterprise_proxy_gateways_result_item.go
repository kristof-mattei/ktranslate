/*
 * VeloCloud Orchestrator API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 3.3.2
 * Contact: support@velocloud.net
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models
import (
	"time"
)

type EnterpriseProxyGetEnterpriseProxyGatewaysResultItem struct {
	ActivationKey string `json:"activationKey,omitempty"`
	ActivationState string `json:"activationState,omitempty"`
	ActivationTime time.Time `json:"activationTime,omitempty"`
	BuildNumber string `json:"buildNumber,omitempty"`
	Certificates []GatewayCertificate `json:"certificates,omitempty"`
	ConnectedEdges int32 `json:"connectedEdges,omitempty"`
	ConnectedEdgeList []GatewayPoolGatewayConnectedEdgeList `json:"connectedEdgeList,omitempty"`
	Created time.Time `json:"created,omitempty"`
	DataCenters []map[string]interface{} `json:"dataCenters,omitempty"`
	Description string `json:"description,omitempty"`
	DeviceId string `json:"deviceId,omitempty"`
	DnsName string `json:"dnsName,omitempty"`
	EndpointPkiMode string `json:"endpointPkiMode,omitempty"`
	EnterpriseAssociations []GatewayEnterpriseAssoc `json:"enterpriseAssociations,omitempty"`
	EnterpriseAssociationCount map[string]int32 `json:"enterpriseAssociationCount,omitempty"`
	Enterprises []Enterprise `json:"enterprises,omitempty"`
	EnterpriseProxyId int32 `json:"enterpriseProxyId,omitempty"`
	GatewayState string `json:"gatewayState,omitempty"`
	HandOffDetail GatewayHandoffDetail `json:"handOffDetail,omitempty"`
	HandOffEdges []GatewayHandoffEdge `json:"handOffEdges,omitempty"`
	Id int32 `json:"id,omitempty"`
	IpAddress string `json:"ipAddress,omitempty"`
	IpsecGatewayDetail GatewayPoolGatewayIpsecGatewayDetail `json:"ipsecGatewayDetail,omitempty"`
	IsLoadBalanced bool `json:"isLoadBalanced,omitempty"`
	LastContact string `json:"lastContact,omitempty"`
	LogicalId string `json:"logicalId,omitempty"`
	Modified time.Time `json:"modified,omitempty"`
	Name string `json:"name,omitempty"`
	NetworkId int32 `json:"networkId,omitempty"`
	Pools []GatewayGatewayPool `json:"pools,omitempty"`
	PrivateIpAddress string `json:"privateIpAddress,omitempty"`
	Roles []GatewayRole `json:"roles,omitempty"`
	ServiceState string `json:"serviceState,omitempty"`
	ServiceUpSince string `json:"serviceUpSince,omitempty"`
	Site GatewaySite `json:"site,omitempty"`
	SiteId int32 `json:"siteId,omitempty"`
	SoftwareVersion string `json:"softwareVersion,omitempty"`
	SystemUpSince string `json:"systemUpSince,omitempty"`
	Utilization float32 `json:"utilization,omitempty"`
	UtilizationDetail GatewayPoolGatewayUtilizationDetail `json:"utilizationDetail,omitempty"`
}
