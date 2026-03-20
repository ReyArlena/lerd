package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// pausedServicesFile returns the path to paused-services.yaml.
func pausedServicesFile() string {
	return filepath.Join(DataDir(), "paused-services.yaml")
}

// loadPausedServices reads the set of manually-stopped service names.
func loadPausedServices() (map[string]bool, error) {
	data, err := os.ReadFile(pausedServicesFile())
	if os.IsNotExist(err) {
		return map[string]bool{}, nil
	}
	if err != nil {
		return nil, err
	}
	var names []string
	if err := yaml.Unmarshal(data, &names); err != nil {
		return nil, err
	}
	m := make(map[string]bool, len(names))
	for _, n := range names {
		m[n] = true
	}
	return m, nil
}

func savePausedServices(paused map[string]bool) error {
	var names []string
	for n := range paused {
		names = append(names, n)
	}
	if err := os.MkdirAll(DataDir(), 0755); err != nil {
		return err
	}
	data, err := yaml.Marshal(names)
	if err != nil {
		return err
	}
	return os.WriteFile(pausedServicesFile(), data, 0644)
}

// ServiceIsPaused returns true if the service was manually stopped by the user.
func ServiceIsPaused(name string) bool {
	paused, err := loadPausedServices()
	if err != nil {
		return false
	}
	return paused[name]
}

// SetServicePaused marks or clears the manual-stop flag for the named service.
func SetServicePaused(name string, paused bool) error {
	m, err := loadPausedServices()
	if err != nil {
		m = map[string]bool{}
	}
	if paused {
		m[name] = true
	} else {
		delete(m, name)
	}
	return savePausedServices(m)
}
