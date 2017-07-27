package crontab

import (
	"fmt"

	extensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	Group   = "kubeheads.github.io"
	Version = "v1"
)

var ReportResource = extensions.CustomResourceDefinition{
	ObjectMeta: metav1.ObjectMeta{
		Name: fmt.Sprintf("%s.%s", CronTabPlural, Group),
	},
	Spec: extensions.CustomResourceDefinitionSpec{
		Group:   Group,
		Version: Version,
		Names: extensions.CustomResourceDefinitionNames{
			Plural: CronTabPlural,
			Kind:   CronTabKind,
		},
		Scope: extensions.ClusterScoped,
	},
}
