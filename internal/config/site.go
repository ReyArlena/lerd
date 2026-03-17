package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Site represents a single registered Lerd site.
type Site struct {
	Name        string `yaml:"name"`
	Domain      string `yaml:"domain"`
	Path        string `yaml:"path"`
	PHPVersion  string `yaml:"php_version"`
	NodeVersion string `yaml:"node_version"`
	Secured     bool   `yaml:"secured"`
}

// SiteRegistry holds all registered sites.
type SiteRegistry struct {
	Sites []Site `yaml:"sites"`
}

// LoadSites reads sites.yaml, returning an empty registry if the file does not exist.
func LoadSites() (*SiteRegistry, error) {
	data, err := os.ReadFile(SitesFile())
	if err != nil {
		if os.IsNotExist(err) {
			return &SiteRegistry{}, nil
		}
		return nil, err
	}

	var reg SiteRegistry
	if err := yaml.Unmarshal(data, &reg); err != nil {
		return nil, err
	}
	return &reg, nil
}

// SaveSites writes the registry to sites.yaml.
func SaveSites(reg *SiteRegistry) error {
	if err := os.MkdirAll(DataDir(), 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(reg)
	if err != nil {
		return err
	}
	return os.WriteFile(SitesFile(), data, 0644)
}

// AddSite appends or updates a site in the registry.
func AddSite(site Site) error {
	reg, err := LoadSites()
	if err != nil {
		return err
	}

	for i, s := range reg.Sites {
		if s.Name == site.Name {
			reg.Sites[i] = site
			return SaveSites(reg)
		}
	}

	reg.Sites = append(reg.Sites, site)
	return SaveSites(reg)
}

// RemoveSite removes a site by name from the registry.
func RemoveSite(name string) error {
	reg, err := LoadSites()
	if err != nil {
		return err
	}

	filtered := reg.Sites[:0]
	for _, s := range reg.Sites {
		if s.Name != name {
			filtered = append(filtered, s)
		}
	}
	reg.Sites = filtered
	return SaveSites(reg)
}

// FindSite returns the site with the given name, or an error if not found.
func FindSite(name string) (*Site, error) {
	reg, err := LoadSites()
	if err != nil {
		return nil, err
	}

	for _, s := range reg.Sites {
		if s.Name == name {
			s := s
			return &s, nil
		}
	}
	return nil, fmt.Errorf("site %q not found", name)
}

// FindSiteByDomain returns the site with the given domain, or an error if not found.
func FindSiteByDomain(domain string) (*Site, error) {
	reg, err := LoadSites()
	if err != nil {
		return nil, err
	}

	for _, s := range reg.Sites {
		if s.Domain == domain {
			s := s
			return &s, nil
		}
	}
	return nil, fmt.Errorf("site with domain %q not found", domain)
}
