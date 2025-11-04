package config

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

// SupportedCustomResources represents the set of custom resources that VPA should treat as supported targets.
// The map key is a schema.GroupVersionKind marshalled via the resourceKey helper.
type SupportedCustomResources map[string]struct{}

// resourceKey builds a unique key for a custom resource identified by apiVersion and kind.
func resourceKey(apiVersion, kind string) string {
	return fmt.Sprintf("%s|%s", apiVersion, kind)
}

// Contains verifies whether the given apiVersion/kind pair is present in the set.
func (s SupportedCustomResources) Contains(apiVersion, kind string) bool {
	if len(s) == 0 {
		return false
	}
	_, ok := s[resourceKey(apiVersion, kind)]
	return ok
}

// ParseSupportedCustomResources converts a comma-separated flag value into a SupportedCustomResources set.
// Each entry must use the format "<apiVersion>:<kind>", for example "k8s.sentio.xyz/v1:DriverJob".
func ParseSupportedCustomResources(flagValue string) (SupportedCustomResources, error) {
	result := SupportedCustomResources{}
	if flagValue == "" {
		return result, nil
	}

	tokens := strings.Split(flagValue, ",")
	for _, token := range tokens {
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}

		parts := strings.Split(token, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid custom resource entry %q, expected <apiVersion>:<kind>", token)
		}

		apiVersion := strings.TrimSpace(parts[0])
		kind := strings.TrimSpace(parts[1])
		if apiVersion == "" || kind == "" {
			return nil, fmt.Errorf("invalid custom resource entry %q, apiVersion and kind must be non-empty", token)
		}

		if _, err := schema.ParseGroupVersion(apiVersion); err != nil {
			return nil, fmt.Errorf("invalid apiVersion %q for custom resource entry %q: %w", apiVersion, token, err)
		}

		result[resourceKey(apiVersion, kind)] = struct{}{}
	}

	return result, nil
}
