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

type DeviceSettingsRefs struct {
	DeviceSettingsauthentication map[string]interface{} `json:"deviceSettings:authentication,omitempty"`
	DeviceSettingscssprovider map[string]interface{} `json:"deviceSettings:css:provider,omitempty"`
	DeviceSettingscsssite map[string]interface{} `json:"deviceSettings:css:site,omitempty"`
	DeviceSettingsdnsbackupProvider map[string]interface{} `json:"deviceSettings:dns:backupProvider,omitempty"`
	DeviceSettingsdnsprimaryProvider map[string]interface{} `json:"deviceSettings:dns:primaryProvider,omitempty"`
	DeviceSettingsdnsprivateProviders map[string]interface{} `json:"deviceSettings:dns:privateProviders,omitempty"`
	DeviceSettingshandOffGatewaysgateways map[string]interface{} `json:"deviceSettings:handOffGateways:gateways,omitempty"`
	DeviceSettingslanallocation map[string]interface{} `json:"deviceSettings:lan:allocation,omitempty"`
	DeviceSettingssecurityVnflicense map[string]interface{} `json:"deviceSettings:securityVnf:license,omitempty"`
	DeviceSettingssecurityVnfservice map[string]interface{} `json:"deviceSettings:securityVnf:service,omitempty"`
	DeviceSettingssegment map[string]interface{} `json:"deviceSettings:segment,omitempty"`
	DeviceSettingssegmentnetflowCollectors map[string]interface{} `json:"deviceSettings:segment:netflowCollectors,omitempty"`
	DeviceSettingssegmentnetflowFilters map[string]interface{} `json:"deviceSettings:segment:netflowFilters,omitempty"`
	DeviceSettingstacacs map[string]interface{} `json:"deviceSettings:tacacs,omitempty"`
	DeviceSettingsvnfsedge map[string]interface{} `json:"deviceSettings:vnfs:edge,omitempty"`
	DeviceSettingsvnfsvnfImage map[string]interface{} `json:"deviceSettings:vnfs:vnfImage,omitempty"`
	DeviceSettingsvpndataCenter map[string]interface{} `json:"deviceSettings:vpn:dataCenter,omitempty"`
	DeviceSettingsvpnedgeHub map[string]interface{} `json:"deviceSettings:vpn:edgeHub,omitempty"`
	DeviceSettingsvpnedgeHubCluster map[string]interface{} `json:"deviceSettings:vpn:edgeHubCluster,omitempty"`
	DeviceSettingswebProxyprovider map[string]interface{} `json:"deviceSettings:webProxy:provider,omitempty"`
}
