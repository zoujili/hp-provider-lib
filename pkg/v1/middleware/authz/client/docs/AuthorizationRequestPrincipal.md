# AuthorizationRequestPrincipal

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | [optional] [default to null]
**Roles** | **[]string** | Array of role ids. Array of permissions ids. It could be empty. In this case, Authz service will look up roles from DB. | [optional] [default to null]
**Scopes** | **[]string** |  | [optional] [default to null]
**Permissions** | **[]string** | Array of permissions ids. It could be empty. In this case, Authz service will look up permissions from DB. | [optional] [default to null]
**Type_** | **string** | Type of principal. E.g. \&quot;user\&quot;, \&quot;client\&quot;. | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


