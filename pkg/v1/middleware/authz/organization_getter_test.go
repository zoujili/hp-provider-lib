package authz_test

import (
	"testing"

	gaiav1 "github.azc.ext.hp.com/hp-business-platform/lib-hpbp-proto-go/gen/proto/go/hpbp/gaia/v1"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/middleware/authz"
	"github.com/stretchr/testify/assert"
)

func TestDefaultOrganizationGetter_ExternalOrganizationID(t *testing.T) {
	orgGetter := authz.NewDefaultOrganizationGetter([]string{"project_id", "application_id"})
	organizationID := orgGetter.ExternalOrganizationID(nil, &gaiav1.GetApplicationRequest{
		ProjectId:     "dd8cae95-190b-4407-bb90-b9c7bc23df9a ",
		ApplicationId: "0f089b78-16fe-4310-a56b-96548e845937",
	})
	assert.Equal(t, "0f089b78-16fe-4310-a56b-96548e845937", organizationID)
}
