package kafkamapper

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ServiceMapper struct {
	clientset *kubernetes.Clientset

	podsByIPCache map[string][]corev1.Pod
}

func NewServiceMapper(clientset *kubernetes.Clientset) *ServiceMapper {
	return &ServiceMapper{
		clientset:     clientset,
		podsByIPCache: map[string][]corev1.Pod{},
	}
}

func (m *ServiceMapper) InitIndexes(ctx context.Context) error {
	if err := m.loadIPToPodMap(ctx); err != nil {
		return err
	}
	return nil
}

func (m *ServiceMapper) GetPodByIP(ip string) (*corev1.Pod, error) {
	pods, ok := m.podsByIPCache[ip]
	if !ok {
		return nil, fmt.Errorf("no pod found for ip %s", ip)
	} else if len(pods) > 1 {
		return nil, fmt.Errorf("multiple pods found for ip %s", ip)
	} else {
		return &pods[0], nil
	}
}

func (m *ServiceMapper) listNamespaces(ctx context.Context) (*corev1.NamespaceList, error) {
	return m.clientset.CoreV1().Namespaces().List(ctx, v1.ListOptions{})
}

func (m *ServiceMapper) listPods(ctx context.Context, ns string) (*corev1.PodList, error) {
	return m.clientset.CoreV1().Pods(ns).List(ctx, v1.ListOptions{})
}

func (m *ServiceMapper) loadIPToPodMap(ctx context.Context) error {
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
			for _, ip := range pod.Status.PodIPs {
				m.podsByIPCache[ip.IP] = append(m.podsByIPCache[ip.IP], pod)
			}
		}
	}

	return nil
}
