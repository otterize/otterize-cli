package kafkamapper

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

type PodMapper struct {
	clientset *kubernetes.Clientset

	serviceIDResolver *ServiceIDResolver

	podsByIPCache map[string][]corev1.Pod
	podToService  map[types.NamespacedName]ServiceIdentity
}

func NewPodMapper() (*PodMapper, error) {
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	serviceIDResolver, err := NewServiceIDResolver()
	if err != nil {
		return nil, err
	}

	return &PodMapper{
		clientset:         clientset,
		serviceIDResolver: serviceIDResolver,
		podsByIPCache:     map[string][]corev1.Pod{},
		podToService:      map[types.NamespacedName]ServiceIdentity{},
	}, nil
}

func (m *PodMapper) InitIndexes(ctx context.Context) error {
	if err := m.loadPodsCache(ctx); err != nil {
		return err
	}
	return nil
}

func (m *PodMapper) GetPodByIP(ip string) (*corev1.Pod, error) {
	pods, ok := m.podsByIPCache[ip]
	if !ok {
		return nil, fmt.Errorf("no pod found for ip %s", ip)
	} else if len(pods) > 1 {
		return nil, fmt.Errorf("multiple pods found for ip %s", ip)
	} else {
		return &pods[0], nil
	}
}

func (m *PodMapper) GetServiceIDByPod(ctx context.Context, pod *corev1.Pod) (ServiceIdentity, error) {
	podName := types.NamespacedName{Namespace: pod.Namespace, Name: pod.Name}
	serviceID, ok := m.podToService[podName]
	if ok {
		return serviceID, nil
	}

	serviceID, err := m.serviceIDResolver.ResolvePodToServiceIdentity(ctx, pod)
	if err != nil {
		return ServiceIdentity{}, err
	}

	m.podToService[podName] = serviceID
	return serviceID, nil
}

func (m *PodMapper) listNamespaces(ctx context.Context) (*corev1.NamespaceList, error) {
	return m.clientset.CoreV1().Namespaces().List(ctx, v1.ListOptions{})
}

func (m *PodMapper) listPods(ctx context.Context, ns string) (*corev1.PodList, error) {
	return m.clientset.CoreV1().Pods(ns).List(ctx, v1.ListOptions{})
}

func (m *PodMapper) loadPodsCache(ctx context.Context) error {
	namespaces, err := m.listNamespaces(ctx)
	if err != nil {
		return err
	}
	for _, ns := range namespaces.Items {
		pods, err := m.listPods(ctx, ns.Name)
		if err != nil {
			return err
		}
		for _, pod := range pods.Items {
			if err != nil {
				logrus.WithError(err).Warn("Skipping resolution to service")
			}
			for _, ip := range pod.Status.PodIPs {
				m.podsByIPCache[ip.IP] = append(m.podsByIPCache[ip.IP], pod)
			}
		}
	}

	return nil
}
