package components

import (
	"context"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	"github.com/lburgazzoli/odh-cli/pkg/resources"
	"github.com/lburgazzoli/odh-cli/pkg/util/client"
	discoverypkg "github.com/lburgazzoli/odh-cli/pkg/util/discovery"
)

// ListComponents lists all component resources from the components.platform.opendatahub.io group.
// It discovers all resource types in the group and aggregates them into a single list.
func ListComponents(
	ctx context.Context,
	client *client.Client,
) (*unstructured.UnstructuredList, error) {
	// Discover all resources in the components.platform.opendatahub.io group
	componentResources, err := discoverypkg.GetGroupVersionResources(
		client.Discovery,
		resources.Components,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to discover component resources: %w", err)
	}

	// Aggregate all component instances
	result := &unstructured.UnstructuredList{
		Items: []unstructured.Unstructured{},
	}

	// List each resource type and add to the result
	for _, resource := range componentResources {
		// Skip subresources (e.g., status, scale)
		// Subresources have a "/" in their name (e.g., "dashboards/status")
		if strings.Contains(resource.Name, "/") {
			continue
		}

		// Additional check: skip if this is explicitly marked as a subresource
		if resource.Kind == "" {
			continue
		}

		gvr := schema.GroupVersionResource{
			Group:    resources.Components.Group,
			Version:  resources.Components.Version,
			Resource: resource.Name,
		}

		list, err := client.Dynamic.Resource(gvr).List(ctx, metav1.ListOptions{})
		if err != nil {
			// Skip resources that can't be listed (e.g., permissions issues)
			continue
		}

		result.Items = append(result.Items, list.Items...)
	}

	return result, nil
}

// GetComponent retrieves a specific component by name.
// Components follow a singleton pattern - there should only be one instance per type.
func GetComponent(
	ctx context.Context,
	dynamicClient dynamic.Interface,
	componentName string,
	componentResource string,
) (*unstructured.Unstructured, error) {
	gvr := schema.GroupVersionResource{
		Group:    resources.Components.Group,
		Version:  resources.Components.Version,
		Resource: componentResource,
	}

	component, err := dynamicClient.Resource(gvr).Get(ctx, componentName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get component %s: %w", componentName, err)
	}

	return component, nil
}

// GetComponentByType retrieves a component by matching its type name (case-insensitive).
// It discovers available component types, finds a match, and returns the singleton instance.
// Examples: "kserve" matches "kserves", "dashboard" matches "dashboards"
func GetComponentByType(
	ctx context.Context,
	client *client.Client,
	typeName string,
) (*unstructured.Unstructured, error) {
	// Discover all component resource types
	componentResources, err := discoverypkg.GetGroupVersionResources(
		client.Discovery,
		resources.Components,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to discover component resources: %w", err)
	}

	// Find matching resource types (case-insensitive)
	var exactMatches []metav1.APIResource
	var partialMatches []metav1.APIResource

	lowerTypeName := strings.ToLower(typeName)

	for _, resource := range componentResources {
		// Skip subresources
		if strings.Contains(resource.Name, "/") || resource.Kind == "" {
			continue
		}

		lowerResourceName := strings.ToLower(resource.Name)

		// Check for exact match
		if lowerResourceName == lowerTypeName {
			exactMatches = append(exactMatches, resource)
		} else if strings.Contains(lowerResourceName, lowerTypeName) || strings.Contains(lowerTypeName, lowerResourceName) {
			// Check for partial match (either way)
			partialMatches = append(partialMatches, resource)
		}
	}

	// Prefer exact matches
	var matchedResource metav1.APIResource
	if len(exactMatches) == 1 {
		matchedResource = exactMatches[0]
	} else if len(exactMatches) > 1 {
		return nil, fmt.Errorf("ambiguous component type %q: multiple exact matches found", typeName)
	} else if len(partialMatches) == 1 {
		matchedResource = partialMatches[0]
	} else if len(partialMatches) > 1 {
		// List the matching types
		var matchNames []string
		for _, m := range partialMatches {
			matchNames = append(matchNames, m.Name)
		}
		return nil, fmt.Errorf("ambiguous component type %q: matches %v", typeName, matchNames)
	} else {
		return nil, fmt.Errorf("no component type matching %q found", typeName)
	}

	// List instances of the matched resource type
	gvr := schema.GroupVersionResource{
		Group:    resources.Components.Group,
		Version:  resources.Components.Version,
		Resource: matchedResource.Name,
	}

	list, err := client.Dynamic.Resource(gvr).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list %s: %w", matchedResource.Name, err)
	}

	// Return the first instance (singleton pattern)
	if len(list.Items) == 0 {
		return nil, fmt.Errorf("no instances of %s found", matchedResource.Name)
	}

	return &list.Items[0], nil
}
