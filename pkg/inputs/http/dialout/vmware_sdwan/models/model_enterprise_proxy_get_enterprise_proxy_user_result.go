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

type EnterpriseProxyGetEnterpriseProxyUserResult struct {
	Id int32 `json:"id,omitempty"`
	Created time.Time `json:"created,omitempty"`
	UserType string `json:"userType,omitempty"`
	Username string `json:"username,omitempty"`
	Domain string `json:"domain,omitempty"`
	Password string `json:"password,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName string `json:"lastName,omitempty"`
	OfficePhone string `json:"officePhone,omitempty"`
	MobilePhone string `json:"mobilePhone,omitempty"`
	IsNative bool `json:"isNative,omitempty"`
	IsActive bool `json:"isActive,omitempty"`
	IsLocked bool `json:"isLocked,omitempty"`
	DisableSecondFactor bool `json:"disableSecondFactor,omitempty"`
	Email string `json:"email,omitempty"`
	LastLogin time.Time `json:"lastLogin,omitempty"`
	Modified time.Time `json:"modified,omitempty"`
	RoleId int32 `json:"roleId,omitempty"`
	RoleName string `json:"roleName,omitempty"`
	EnterpriseProxyId int32 `json:"enterpriseProxyId,omitempty"`
	NetworkId int32 `json:"networkId,omitempty"`
}
