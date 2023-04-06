package kubernetes

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
	"sync"
)

const (
	resource          = "events"
	namespace         = metav1.NamespaceAll
	maxEventBatchSize = 100 //TODO 可能会有事件丢失
)

type EventSrouce struct {
	rateLimitingQeueue workqueue.RateLimitingInterface
	informer           cache.SharedIndexInformer
	mutex              sync.Mutex
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
		if len(events) == maxEventBatchSize {
			break
		}
	}
	return events, nil
}

func (es *EventSrouce) Run(stopCh <-chan struct{}) {
	go es.informer.Run(stopCh)
}

func (es *EventSrouce) addedToQueueEvent(obj interface{}) {
	if obj == nil {
		return
	}
	event, ok := obj.(*corev1.Event)
	if ok {
		event.SetManagedFields(nil)
		es.rateLimitingQeueue.Add(event)
	}
}

func NewKubernetesEventSource(clientset *kubernetes.Clientset) *EventSrouce {
	source := &EventSrouce{
		rateLimitingQeueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "events"),
	}
	lw := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), resource, namespace, fields.Everything())
	source.informer = cache.NewSharedIndexInformer(lw, &corev1.Event{}, 0, cache.Indexers{})
	source.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: source.addedToQueueEvent,
		UpdateFunc: func(oldObj, newObj interface{}) {
			source.addedToQueueEvent(newObj)
		}})
	return source
}
