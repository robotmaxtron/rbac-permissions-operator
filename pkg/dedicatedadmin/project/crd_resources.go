// Copyright 2018 RedHat
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package project

import (
	"google.golang.org/genproto/googleapis/ads/googleads/v0/resources"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ClusterRoles = map[string]rbacv1.ClusterRole{
	"dedicated-admins-cluster-crds": {
		apiVersion: "rbac.authorization.k8s.io/v1",
		kind: "ClusterRole",
		metadata: metav1.ObjectMeta{
			Name: "dedicated-admins-cluster-crds",
		},
		Rules: []rbacv1.ClusterRole.Rules{
			{
				APIGroup: "rbac.authorization.k8s.io",
				attributeRestrictions: "null",
				resources: []rbacv1.ClusterRole.Rules.resources{
					""
				},
				verbs: []rbacv1.ClusterRole.Rules.verbs{
					""
				},
			},
		},
	},
		"dedicated-admins-project-crds": {
			apiVersion: "rbac.authorization.k8s.io/v1",
			kind: "ClusterRole",
			metadata: metav1.ObjectMeta{
				Name: "dedicated-admins-project-crds",
				Namespace: ""
			},
			Rules: []rbacv1.ClusterRole.Rules{
				{
					APIGroup: "rbac.authorization.k8s.io",
					attributeRestrictions: "null",
					resources: []rbacv1.ClusterRole.Rules.resources{
						""
					},
					verbs: []rbacv1.ClusterRole.Rules.verbs{
						""
					},
				},
			},
		},
}