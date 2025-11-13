package resources

import "k8s.io/apimachinery/pkg/runtime/schema"

// Components contains the group and version for ODH/RHOAI components.
// Individual component types (dashboards, kserves, etc.) are discovered dynamically.
var Components = schema.GroupVersion{
	Group:   "components.platform.opendatahub.io",
	Version: "v1alpha1",
}
