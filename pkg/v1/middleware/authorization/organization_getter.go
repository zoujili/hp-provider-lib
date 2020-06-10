package authorization

import (
	"context"
	"reflect"
	"strings"
)

// DefaultOrganizationGetter .
type DefaultOrganizationGetter struct {
	organizationTags      []string // organizationTags is ordered, getter will get the organization_id from the req by tags
	organizationTagsMap   map[string]int
	defaultOrganizationID string
}

// NewDefaultOrganizationGetter I suggest support organizationTags with hard code
func NewDefaultOrganizationGetter(organizationTags []string, defaultOrganizationID string) *DefaultOrganizationGetter {
	g := &DefaultOrganizationGetter{
		organizationTags:      organizationTags,
		organizationTagsMap:   map[string]int{},
		defaultOrganizationID: defaultOrganizationID,
	}
	for _, tag := range organizationTags {
		g.organizationTagsMap[tag] = 1
	}
	return g
}

func getTag(tag string) string {
	return strings.Split(tag, ",")[0]
}

// ExternalOrganizationID DefaultOrganizationGetter will just scan the first class params of the req
func (g *DefaultOrganizationGetter) ExternalOrganizationID(_ context.Context, req interface{}) string {
	if req == nil {
		return ""
	}
	v := reflect.ValueOf(req)
	t := reflect.TypeOf(req)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}
	if v.Kind() != reflect.Struct {
		return ""
	}
	orgMap := map[string]string{}
	for i := 0; i < v.NumField(); i++ {
		key := getTag(t.Field(i).Tag.Get("json"))
		if _, ok := g.organizationTagsMap[key]; ok && v.Field(i).Kind() == reflect.String {
			orgMap[key] = v.Field(i).String()
		}
	}
	organizationID := ""
	for _, tag := range g.organizationTags {
		if orgMap[tag] != "" {
			organizationID = orgMap[tag]
		}
	}
	if organizationID == "" {
		organizationID = g.defaultOrganizationID
	}
	return organizationID
}
