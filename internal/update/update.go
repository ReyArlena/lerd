package update

import (
	"fmt"
	"net/http"
	"strings"
)

// ReleasesBaseURL is the base GitHub releases URL. Overridable in tests.
var ReleasesBaseURL = "https://github.com/geodro/lerd/releases"

// FetchLatestVersion returns the latest published release tag from GitHub.
func FetchLatestVersion() (string, error) {
	url := ReleasesBaseURL + "/latest"
	client := &http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(http.MethodGet, url, nil) //nolint:noctx
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "lerd-cli")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusFound && resp.StatusCode != http.StatusMovedPermanently {
		return "", fmt.Errorf("unexpected status from %s: HTTP %d", url, resp.StatusCode)
	}
	location := resp.Header.Get("Location")
	if location == "" {
		return "", fmt.Errorf("no Location header in redirect from %s", url)
	}
	parts := strings.Split(location, "/tag/")
	if len(parts) != 2 || parts[1] == "" {
		return "", fmt.Errorf("unexpected release URL format: %s", location)
	}
	return parts[1], nil
}

// StripV removes a leading "v" from a version string.
func StripV(v string) string {
	if len(v) > 0 && v[0] == 'v' {
		return v[1:]
	}
	return v
}
