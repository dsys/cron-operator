package crontab

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CronTab struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`
}

type CronTabList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata,omitempty"`
	Items         []*CronTab `json:"items"`
}
