package authorization_test

//
//type mockUserGetter struct {
//}
//
//func (m *mockUserGetter) UserID(_ context.Context) string {
//	return "29ab5fff-c81d-44f4-82b5-6d619d453f02"
//}
//
//func getInterceptor(t *testing.T) *authorization.Interceptor {
//	return authorization.NewInterceptor(
//		authorization.ConfigWithDefaultOrganizationGetter([]string{"project_id"}),
//		authorization.ConfigWithAuthzServiceHost("hpbp-dev.hpbp.io"),
//		authorization.ConfigWithUserGetter(&mockUserGetter{}),
//	)
//}
//
//func TestInterceptor_UnaryServerInterceptor(t *testing.T) {
//	t.Skip()
//	interceptor := getInterceptor(t)
//	_, err := interceptor.UnaryServerInterceptor()(
//		context.Background(),
//		&gaiav1.GetApplicationRequest{
//			ProjectId:     "dd8cae95-190b-4407-bb90-b9c7bc23df9a",
//			ApplicationId: "jacky_branch_org2-2",
//		},
//		&grpc.UnaryServerInfo{
//			FullMethod: "http://company-service/api/v1/commpanys/1234",
//		},
//		func(ctx context.Context, req interface{}) (interface{}, error) {
//			return nil, nil
//		},
//	)
//	assert.NoError(t, err)
//}
