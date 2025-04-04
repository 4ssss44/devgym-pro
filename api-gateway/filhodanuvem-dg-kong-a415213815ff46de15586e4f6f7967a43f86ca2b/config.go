package kong

import (
	"net/http"
	"slices"
	"time"

	"github.com/devgymbr/kong/routes"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Services             []Service `yaml:"services"`
	lastModificationTime time.Time
}

func (c *Config) ModifiedSince(t time.Time) bool {
	return c.lastModificationTime.After(t) || c.lastModificationTime.Equal(t)
}

func (c *Config) Refresh(data []byte, modTime time.Time) error {
	err := yaml.Unmarshal(data, &c)
	if err != nil {
		return err
	}
	c.lastModificationTime = modTime

	for i := range c.Services {
		service := &c.Services[i]
		for j := range c.Services[i].Routes {
			service.Routes[j].PathRegexp, err = routes.Parse(service.Routes[j].Paths[0])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Config) FindServiceRoute(r *http.Request) (*Service, *Route) {
	for _, service := range c.Services {
		for _, route := range service.Routes {
			if route.PathRegexp == nil {
				continue
			}

			if slices.Index[[]string, string](route.Methods, r.Method) != -1 &&
				route.PathRegexp.MatchString(r.URL.Path) {

				return &service, &route
			}
		}
	}

	return nil, nil
}
