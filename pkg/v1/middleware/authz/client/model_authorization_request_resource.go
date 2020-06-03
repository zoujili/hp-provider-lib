/*
 * Authorization Service
 *
 * Authorization Service manages RBAC and ABAC policies
 *
 * API version: 0.0.1
 * Contact: xi.he@hp.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package client

type AuthorizationRequestResource struct {
	// resource id
	Id string `json:"id,omitempty"`
	// resource name defined in metadata. E.g. company, store...
	Name string `json:"name,omitempty"`
}
