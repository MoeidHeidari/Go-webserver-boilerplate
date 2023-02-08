package kubes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func (u KubeRequest) GetEvents(namespace string) (watch.Interface, error) {
	opts := metav1.ListOptions{
		TypeMeta: metav1.TypeMeta{
			Kind: "POD",
		},
	}
	err := u.Clientset.CoreV1().Events(namespace).DeleteCollection(context.TODO(), metav1.DeleteOptions{}, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	events, err := u.Clientset.CoreV1().Events(namespace).Watch(context.TODO(), opts)
	if err != nil || events == nil {
		return nil, err
	}
	return events, nil
}
