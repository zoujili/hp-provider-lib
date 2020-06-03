# \AuthorizationApi

All URIs are relative to *http://hpbp-authz-service.hpbp/hpbp-authz/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AuthorizeRequest**](AuthorizationApi.md#AuthorizeRequest) | **Post** /authorizations | AuthorizeRequest


# **AuthorizeRequest**
> AuthorizationResult AuthorizeRequest(ctx, body, optional)
AuthorizeRequest

Authoriza request on both RBAC and ABAC.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**AuthorizationRequest**](AuthorizationRequest.md)|  | 
 **optional** | ***AuthorizationApiAuthorizeRequestOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a AuthorizationApiAuthorizeRequestOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **rbac** | **optional.String**| Specify if RBAC check is needed for this request. The value must be one of true or false. Default is true. | 
 **abac** | **optional.String**| Specify if ABAC check is needed for this request. The value must be one of true or false. Default is true. | 

### Return type

[**AuthorizationResult**](AuthorizationResult.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

