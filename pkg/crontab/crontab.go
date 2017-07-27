package crontab

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

const (
	CronTabPlural = "crontabs"
	CronTabKind   = "CronTab"
)

type CronTabGetter interface {
	CronTabs(namespace string) CronTabInterface
}

type CronTabInterface interface {
	Create(*CronTab) (*CronTab, error)
	Get(name string) (*CronTab, error)
	Update(*CronTab) (*CronTab, error)
	Delete(name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (runtime.Object, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type cronTabs struct {
	restClient rest.Interface
	client     *dynamic.ResourceClient
	ns         string
}

func newCronTabs(r rest.Interface, c *dynamic.Client, namespace string) *cronTabs {
	return &cronTabs{
		r,
		c.Resource(
			&metav1.APIResource{
				Kind:       CronTabKind,
				Name:       CronTabPlural,
				Namespaced: true,
			},
			namespace,
		),
		namespace,
	}
}

func (p *cronTabs) Create(o *CronTab) (*CronTab, error) {
	up, err := UnstructuredFromCronTab(o)
	if err != nil {
		return nil, err
	}

	up, err = p.client.Create(up)
	if err != nil {
		return nil, err
	}

	return CronTabFromUnstructured(up)
}

func (p *cronTabs) Get(name string) (*CronTab, error) {
	obj, err := p.client.Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return CronTabFromUnstructured(obj)
}

func (p *cronTabs) Update(o *CronTab) (*CronTab, error) {
	up, err := UnstructuredFromCronTab(o)
	if err != nil {
		return nil, err
	}

	up, err = p.client.Update(up)
	if err != nil {
		return nil, err
	}

	return CronTabFromUnstructured(up)
}

func (p *cronTabs) Delete(name string, options *metav1.DeleteOptions) error {
	return p.client.Delete(name, options)
}

func (p *cronTabs) List(opts metav1.ListOptions) (runtime.Object, error) {
	req := p.restClient.Get().
		Namespace(p.ns).
		Resource(CronTabPlural).
		FieldsSelectorParam(nil)

	b, err := req.DoRaw()
	if err != nil {
		return nil, err
	}
	var crontabs CronTabList
	return &crontabs, json.Unmarshal(b, &crontabs)
}

func (p *cronTabs) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	r, err := p.restClient.Get().
		Prefix("watch").
		Namespace(p.ns).
		Resource(CronTabPlural).
		FieldsSelectorParam(nil).
		Stream()
	if err != nil {
		return nil, err
	}
	return watch.NewStreamWatcher(&reportDecoder{
		dec:   json.NewDecoder(r),
		close: r.Close,
	}), nil
}

// CronTabFromUnstructured unmarshals a CronTab object.
func CronTabFromUnstructured(r *unstructured.Unstructured) (*CronTab, error) {
	b, err := json.Marshal(r.Object)
	if err != nil {
		return nil, err
	}
	var p CronTab
	if err := json.Unmarshal(b, &p); err != nil {
		return nil, err
	}
	p.TypeMeta.Kind = CronTabKind
	p.TypeMeta.APIVersion = Group + "/" + Version
	return &p, nil
}

// UnstructuredFromCronTab marshals a CronTab object.
func UnstructuredFromCronTab(p *CronTab) (*unstructured.Unstructured, error) {
	p.TypeMeta.Kind = CronTabKind
	p.TypeMeta.APIVersion = Group + "/" + Version
	b, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	var r unstructured.Unstructured
	if err := json.Unmarshal(b, &r.Object); err != nil {
		return nil, err
	}
	return &r, nil
}

type reportDecoder struct {
	dec   *json.Decoder
	close func() error
}

func (d *reportDecoder) Close() {
	d.close()
}

func (d *reportDecoder) Decode() (action watch.EventType, object runtime.Object, err error) {
	var e struct {
		Type   watch.EventType
		Object CronTab
	}
	if err := d.dec.Decode(&e); err != nil {
		return watch.Error, nil, err
	}
	return e.Type, &e.Object, nil
}
