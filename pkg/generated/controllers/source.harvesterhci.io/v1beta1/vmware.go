/*
Copyright 2022 Rancher Labs, Inc.

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

// Code generated by main. DO NOT EDIT.

package v1beta1

import (
	"context"
	"time"

	v1beta1 "github.com/harvester/vm-import-controller/pkg/apis/source.harvesterhci.io/v1beta1"
	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/kv"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type VmwareHandler func(string, *v1beta1.Vmware) (*v1beta1.Vmware, error)

type VmwareController interface {
	generic.ControllerMeta
	VmwareClient

	OnChange(ctx context.Context, name string, sync VmwareHandler)
	OnRemove(ctx context.Context, name string, sync VmwareHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() VmwareCache
}

type VmwareClient interface {
	Create(*v1beta1.Vmware) (*v1beta1.Vmware, error)
	Update(*v1beta1.Vmware) (*v1beta1.Vmware, error)
	UpdateStatus(*v1beta1.Vmware) (*v1beta1.Vmware, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1beta1.Vmware, error)
	List(namespace string, opts metav1.ListOptions) (*v1beta1.VmwareList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.Vmware, err error)
}

type VmwareCache interface {
	Get(namespace, name string) (*v1beta1.Vmware, error)
	List(namespace string, selector labels.Selector) ([]*v1beta1.Vmware, error)

	AddIndexer(indexName string, indexer VmwareIndexer)
	GetByIndex(indexName, key string) ([]*v1beta1.Vmware, error)
}

type VmwareIndexer func(obj *v1beta1.Vmware) ([]string, error)

type vmwareController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewVmwareController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) VmwareController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &vmwareController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromVmwareHandlerToHandler(sync VmwareHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1beta1.Vmware
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1beta1.Vmware))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *vmwareController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1beta1.Vmware))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateVmwareDeepCopyOnChange(client VmwareClient, obj *v1beta1.Vmware, handler func(obj *v1beta1.Vmware) (*v1beta1.Vmware, error)) (*v1beta1.Vmware, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *vmwareController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *vmwareController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *vmwareController) OnChange(ctx context.Context, name string, sync VmwareHandler) {
	c.AddGenericHandler(ctx, name, FromVmwareHandlerToHandler(sync))
}

func (c *vmwareController) OnRemove(ctx context.Context, name string, sync VmwareHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromVmwareHandlerToHandler(sync)))
}

func (c *vmwareController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *vmwareController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *vmwareController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *vmwareController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *vmwareController) Cache() VmwareCache {
	return &vmwareCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *vmwareController) Create(obj *v1beta1.Vmware) (*v1beta1.Vmware, error) {
	result := &v1beta1.Vmware{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *vmwareController) Update(obj *v1beta1.Vmware) (*v1beta1.Vmware, error) {
	result := &v1beta1.Vmware{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *vmwareController) UpdateStatus(obj *v1beta1.Vmware) (*v1beta1.Vmware, error) {
	result := &v1beta1.Vmware{}
	return result, c.client.UpdateStatus(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *vmwareController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *vmwareController) Get(namespace, name string, options metav1.GetOptions) (*v1beta1.Vmware, error) {
	result := &v1beta1.Vmware{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *vmwareController) List(namespace string, opts metav1.ListOptions) (*v1beta1.VmwareList, error) {
	result := &v1beta1.VmwareList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *vmwareController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *vmwareController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v1beta1.Vmware, error) {
	result := &v1beta1.Vmware{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type vmwareCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *vmwareCache) Get(namespace, name string) (*v1beta1.Vmware, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1beta1.Vmware), nil
}

func (c *vmwareCache) List(namespace string, selector labels.Selector) (ret []*v1beta1.Vmware, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.Vmware))
	})

	return ret, err
}

func (c *vmwareCache) AddIndexer(indexName string, indexer VmwareIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1beta1.Vmware))
		},
	}))
}

func (c *vmwareCache) GetByIndex(indexName, key string) (result []*v1beta1.Vmware, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1beta1.Vmware, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1beta1.Vmware))
	}
	return result, nil
}

type VmwareStatusHandler func(obj *v1beta1.Vmware, status v1beta1.VmwareClusterStatus) (v1beta1.VmwareClusterStatus, error)

type VmwareGeneratingHandler func(obj *v1beta1.Vmware, status v1beta1.VmwareClusterStatus) ([]runtime.Object, v1beta1.VmwareClusterStatus, error)

func RegisterVmwareStatusHandler(ctx context.Context, controller VmwareController, condition condition.Cond, name string, handler VmwareStatusHandler) {
	statusHandler := &vmwareStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromVmwareHandlerToHandler(statusHandler.sync))
}

func RegisterVmwareGeneratingHandler(ctx context.Context, controller VmwareController, apply apply.Apply,
	condition condition.Cond, name string, handler VmwareGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &vmwareGeneratingHandler{
		VmwareGeneratingHandler: handler,
		apply:                   apply,
		name:                    name,
		gvk:                     controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterVmwareStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type vmwareStatusHandler struct {
	client    VmwareClient
	condition condition.Cond
	handler   VmwareStatusHandler
}

func (a *vmwareStatusHandler) sync(key string, obj *v1beta1.Vmware) (*v1beta1.Vmware, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status.DeepCopy()
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(&newStatus, "", nil)
		} else {
			a.condition.SetError(&newStatus, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, &newStatus) {
		if a.condition != "" {
			// Since status has changed, update the lastUpdatedTime
			a.condition.LastUpdated(&newStatus, time.Now().UTC().Format(time.RFC3339))
		}

		var newErr error
		obj.Status = newStatus
		newObj, newErr := a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
		if newErr == nil {
			obj = newObj
		}
	}
	return obj, err
}

type vmwareGeneratingHandler struct {
	VmwareGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *vmwareGeneratingHandler) Remove(key string, obj *v1beta1.Vmware) (*v1beta1.Vmware, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1beta1.Vmware{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *vmwareGeneratingHandler) Handle(obj *v1beta1.Vmware, status v1beta1.VmwareClusterStatus) (v1beta1.VmwareClusterStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.VmwareGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
