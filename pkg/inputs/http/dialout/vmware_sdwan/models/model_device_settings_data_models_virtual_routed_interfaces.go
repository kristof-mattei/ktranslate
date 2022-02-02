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

type DeviceSettingsDataModelsVirtualRoutedInterfaces struct {
	Name               string                                    `json:"name,omitempty"`
	Disabled           bool                                      `json:"disabled,omitempty"`
	Addressing         DeviceSettingsDataModelsVirtualAddressing `json:"addressing,omitempty"`
	WanOverlay         string                                    `json:"wanOverlay,omitempty"`
	NatDirect          bool                                      `json:"natDirect,omitempty"`
	Ospf               DeviceSettingsDataModelsVirtualOspf       `json:"ospf,omitempty"`
	VlanId             int32                                     `json:"vlanId,omitempty"`
	L2                 EdgeDeviceSettingsDataL2                  `json:"l2,omitempty"`
	UnderlayAccounting bool                                      `json:"underlayAccounting,omitempty"`
	Trusted            bool                                      `json:"trusted,omitempty"`
	Rpf                string                                    `json:"rpf,omitempty"`
}
