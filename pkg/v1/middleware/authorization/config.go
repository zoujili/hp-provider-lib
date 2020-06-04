package authorization

type Config struct {
	AuthzServiceAddr   string
	UserGetter         IUserGetter
	OrganizationGetter IOrganizationGetter
	Skipper            Skipper
}

type Skipper func(string) bool

type ConfigFunc func(*Config)

func ConfigWithAuthzClientAddr(addr string) ConfigFunc {
	return func(c *Config) {
		c.AuthzServiceAddr = addr
	}
}

func ConfigWithSkipper(skipper Skipper) ConfigFunc {
	return func(c *Config) {
		c.Skipper = skipper
	}
}

func ConfigWithUserGetter(userGetter IUserGetter) ConfigFunc {
	return func(c *Config) {
		c.UserGetter = userGetter
	}
}

func ConfigWithDefaultUserGetter() ConfigFunc {
	return func(c *Config) {
		c.UserGetter = NewDefaultUserGetter()
	}
}

func ConfigWithOrganizationGetter(organizationGetter IOrganizationGetter) ConfigFunc {
	return func(c *Config) {
		c.OrganizationGetter = organizationGetter
	}
}

func ConfigWithDefaultOrganizationGetter(organizationTags []string) ConfigFunc {
	return func(c *Config) {
		c.OrganizationGetter = NewDefaultOrganizationGetter(organizationTags)
	}
}
