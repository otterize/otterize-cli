package kafkamapper

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

const (
	ServiceNameAnnotation = "intents.otterize.com/service-name"
)

type ServiceIDResolver struct {
	clientset     *kubernetes.Clientset
	dynamicclient *dynamic.DynamicClient
}

func NewServiceIDResolver() (*ServiceIDResolver, error) {
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	dynamicclient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &ServiceIDResolver{clientset: clientset, dynamicclient: dynamicclient}, nil
}

func ResolvePodToServiceIdentityUsingAnnotationOnly(pod *corev1.Pod) (string, bool) {
	annotation, ok := pod.Annotations[ServiceNameAnnotation]
	return annotation, ok
}

type ServiceIdentity struct {
	Name string
	// OwnerObject used to resolve the service name. May be nil if service name was resolved using annotation.
	OwnerObject client.Object
}

// ResolvePodToServiceIdentity resolves a pod object to its otterize service ID, referenced in intents objects.
// It calls GetOwnerObject to recursively iterates over the pod's owner reference hierarchy until reaching a root owner reference.
// In case the pod is annotated with an "intents.otterize.com/service-name" annotation, that annotation's value will override
// any owner reference name as the service name.
func (r *ServiceIDResolver) ResolvePodToServiceIdentity(ctx context.Context, pod *corev1.Pod) (ServiceIdentity, error) {
	annotatedServiceName, ok := ResolvePodToServiceIdentityUsingAnnotationOnly(pod)
	if ok {
		return ServiceIdentity{Name: annotatedServiceName}, nil
	}
	ownerObj, err := r.GetOwnerObject(ctx, pod)
	if err != nil {
		return ServiceIdentity{}, err
	}

	return ServiceIdentity{Name: ownerObj.GetName(), OwnerObject: ownerObj}, nil
}

// GetOwnerObject recursively iterates over the pod's owner reference hierarchy until reaching a root owner reference
// and returns it.
func (r *ServiceIDResolver) GetOwnerObject(ctx context.Context, pod *corev1.Pod) (client.Object, error) {
	log := logrus.WithFields(logrus.Fields{"pod": pod.Name, "namespace": pod.Namespace})
	var obj client.Object
	obj = pod
	for len(obj.GetOwnerReferences()) > 0 {
		owner := obj.GetOwnerReferences()[0]

		gv, err := schema.ParseGroupVersion(owner.APIVersion)
		if err != nil {
			return nil, err
		}

		resourceID := gv.WithResource(strings.ToLower(owner.Kind) + "s") // bah
		ownerObj, err := r.dynamicclient.Resource(resourceID).Namespace(obj.GetNamespace()).Get(ctx, owner.Name, v1.GetOptions{})
		if err != nil && errors.IsForbidden(err) {
			// We don't have permissions for further resolving of the owner object,
			// and so we treat it as the identity.
			log.WithFields(logrus.Fields{"owner": owner.Name, "resourceID": resourceID}).Warning(
				"permission error resolving owner, will use owner object as service identifier",
			)
			ownerObj.SetName(owner.Name)
			return ownerObj, nil
		} else if err != nil {
			log.WithFields(logrus.Fields{"owner": owner.Name, "ownerKind": resourceID}).Error(
				"failed querying owner reference",
			)
			return nil, fmt.Errorf("error querying owner reference: %w", err)
		}

		// recurse parent owner reference
		obj = ownerObj
	}

	log.WithFields(logrus.Fields{"owner": obj.GetName(), "ownerKind": obj.GetObjectKind().GroupVersionKind()}).Debug("pod resolved to owner name")
	return obj, nil
}
