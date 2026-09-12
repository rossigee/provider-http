// Package apis contains Kubernetes API for the Http provider.
package apis

import (
	disposablerequestv1beta1 "github.com/rossigee/provider-http/apis/disposablerequest/v1beta1"
	requestv1beta1 "github.com/rossigee/provider-http/apis/request/v1beta1"
	httpv1beta1 "github.com/rossigee/provider-http/apis/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes,
		httpv1beta1.SchemeBuilder.AddToScheme,
		disposablerequestv1beta1.SchemeBuilder.AddToScheme,
		requestv1beta1.SchemeBuilder.AddToScheme,
	)
}

// AddToSchemes may be used to add all resources defined in the project to a Scheme
var AddToSchemes runtime.SchemeBuilder

// AddToScheme adds all Resources to the Scheme
func AddToScheme(s *runtime.Scheme) error {
	return AddToSchemes.AddToScheme(s)
}
