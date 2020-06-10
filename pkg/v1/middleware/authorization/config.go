package authorization

// Config .
type Config struct {
	AuthzServiceHost   string
	UserGetter         IUserGetter
	OrganizationGetter IOrganizationGetter
	Skipper            Skipper
}

// Skipper the skipper func type
type Skipper func(string) bool

// ConfigFunc .
type ConfigFunc func(*Config)

// ConfigWithAuthzServiceHost set config with the host of authz service
func ConfigWithAuthzServiceHost(host string) ConfigFunc {
	return func(c *Config) {
		c.AuthzServiceHost = host
	}
}

// ConfigWithSkipper set skipper for authorization.
func ConfigWithSkipper(skipper Skipper) ConfigFunc {
	return func(c *Config) {
		c.Skipper = skipper
	}
}

// ConfigWithUserGetter set UserGetter for authorization.
func ConfigWithUserGetter(userGetter IUserGetter) ConfigFunc {
	return func(c *Config) {
		c.UserGetter = userGetter
	}
}

// ConfigWithDefaultUserGetter set default UserGetter
func ConfigWithDefaultUserGetter() ConfigFunc {
	return func(c *Config) {
		c.UserGetter = NewDefaultUserGetter()
	}
}

// ConfigWithOrganizationGetter set Organization Getter for authorization
func ConfigWithOrganizationGetter(organizationGetter IOrganizationGetter) ConfigFunc {
	return func(c *Config) {
		c.OrganizationGetter = organizationGetter
	}
}

// ConfigWithDefaultOrganizationGetter set default Organization Getter for authorization
func ConfigWithDefaultOrganizationGetter(organizationTags []string, defaultOrganizationID string) ConfigFunc {
	return func(c *Config) {
		c.OrganizationGetter = NewDefaultOrganizationGetter(organizationTags, defaultOrganizationID)
	}
}
