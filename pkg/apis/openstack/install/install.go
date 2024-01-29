// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package install

import (
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	"github.com/gardener/machine-controller-manager-provider-openstack/pkg/apis/openstack"
	"github.com/gardener/machine-controller-manager-provider-openstack/pkg/apis/openstack/v1alpha1"
	"github.com/gardener/machine-controller-manager-provider-openstack/pkg/apis/openstack/v1alpha2"
)

var (
	schemeBuilder = runtime.NewSchemeBuilder(
		v1alpha1.AddToScheme,
		v1alpha2.AddToScheme,
		openstack.AddToScheme,
		setVersionPriority,
	)

	// AddToScheme adds all APIs to the scheme
	AddToScheme = schemeBuilder.AddToScheme
)

func setVersionPriority(scheme *runtime.Scheme) error {
	return scheme.SetVersionPriority(v1alpha2.SchemeGroupVersion, v1alpha1.SchemeGroupVersion)
}

// Install installs all APIs in the current scheme.
func Install(scheme *runtime.Scheme) *runtime.Scheme {
	utilruntime.Must(AddToScheme(scheme))
	return scheme
}
