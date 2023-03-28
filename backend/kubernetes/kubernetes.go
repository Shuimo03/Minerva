package kubernetes

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

type EventSrouce struct {
	rateLimitingQeueue workqueue.RateLimitingInterface
}

func (es *EventSrouce) GetEvent() ([]*corev1.Event, error) {
	queue := es.rateLimitingQeueue
	events := make([]*corev1.Event, 0)
	for {
		item, quit := queue.Get()
		if quit {
			break
		}
		event, ok := item.(*corev1.Event)
		if !ok {
			klog.Errorf("item is not an Event")
			queue.Forget(item)
			continue
		}
		events = append(events, event)
	}
	return events, nil
}

func NewKubernetesEventSource() *EventSrouce {
	source := &EventSrouce{
		rateLimitingQeueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "events"),
	}
	eventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			event, ok := obj.(*corev1.Event)
			if !ok {
				klog.Errorf("object is not an Event")
				return
			}
			klog.Infof("kubernetes '%s' added", event.Name)
			queue.Add(event)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			event, ok := newObj.(*corev1.Event)
			if !ok {
				klog.Errorf("object is not an Event")
				return
			}
			klog.Infof("kubernetes '%s' updated", event.Name)
			queue.Add(event)
		},
		DeleteFunc: func(obj interface{}) {
			event, ok := obj.(*corev1.Event)
			if !ok {
				klog.Errorf("object is not an Event")
				return
			}
			klog.Infof("kubernetes '%s' deleted", event.Name)
			queue.Add(event)
		},
	}
	return source
}
