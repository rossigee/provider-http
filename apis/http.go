/*
Copyright 2020 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package apis contains Kubernetes API for the Http provider.
package apis

import (
	httpv1alpha1 "github.com/rossigee/provider-http/apis/v1alpha1"
	disposablerequestv1alpha1 "github.com/rossigee/provider-http/apis/disposablerequest/v1alpha1"
	disposablerequestv1alpha2 "github.com/rossigee/provider-http/apis/disposablerequest/v1alpha2"
	disposablerequestv1beta1 "github.com/rossigee/provider-http/apis/disposablerequest/v1beta1"
	requestv1alpha1 "github.com/rossigee/provider-http/apis/request/v1alpha1"
	requestv1alpha2 "github.com/rossigee/provider-http/apis/request/v1alpha2"
	requestv1beta1 "github.com/rossigee/provider-http/apis/request/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes,
		httpv1alpha1.SchemeBuilder.AddToScheme,
		disposablerequestv1alpha1.SchemeBuilder.AddToScheme,
		disposablerequestv1alpha2.SchemeBuilder.AddToScheme,
		requestv1alpha1.SchemeBuilder.AddToScheme,
		requestv1alpha2.SchemeBuilder.AddToScheme,
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
