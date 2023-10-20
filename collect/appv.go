package collect

import (
	"context"
	"fmt"
	"strings"

	"github.com/Albertwzp/cli-go/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/flowcontrol"
)

type Tag struct {
	ns        string
	app       string
	upImage   string
	downImage string
	upCount   int32
	downCount int32
}

func AppFilter(filter string) ([]Tag, error) {
	k, _ := config.GetK8sConfig()
	k.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(100, 100)
	clientSet := kubernetes.NewForConfigOrDie(k)
	dep := clientSet.AppsV1().Deployments("")
	opt := metav1.ListOptions{
		LabelSelector: filter,
	}
	depL, err := dep.List(context.TODO(), opt)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	var tag []Tag
	for _, c := range depL.Items {
		tag_app := Tag{}
		ns := c.GetNamespace()
		labels := c.GetLabels()
		tag_app.ns = ns
		tag_app.app = labels["app_name"]
		rsLabel := "app=" + c.Spec.Template.Labels["app"]

		//fmt.Printf("%v:%v\n", c.Status.Replicas, c.Status.ReadyReplicas)
		if c.Status.Replicas > 0 {
			if c.Status.ReadyReplicas == c.Status.Replicas {
				tag_app.upImage = strings.Split(c.Spec.Template.Spec.Containers[0].Image, ":")[1]
				//tag_app.downImage = ""
				tag_app.upCount = c.Status.ReadyReplicas
				tag_app.downCount = -1
			} else {
				rsL(clientSet, &tag_app, ns, rsLabel)
			}
		}
		tag = append(tag, tag_app)
	}
	return tag, nil
}

func rsL(kc *kubernetes.Clientset, podS *Tag, ns, label string) {
	rs := kc.AppsV1().ReplicaSets(ns)
	opt := metav1.ListOptions{
		LabelSelector: label,
	}
	rsL, err := rs.List(context.TODO(), opt)
	if err != nil {
		fmt.Println(err.Error())
	}
	var count int32
	for _, rs := range rsL.Items {
		if rs.Status.Replicas > 0 {
			count += rs.Status.Replicas
			if rs.Status.ReadyReplicas > 0 {
				podS.upImage = strings.Split(rs.Spec.Template.Spec.Containers[0].Image, ":")[1]
				podS.upCount = rs.Status.ReadyReplicas
			} else {
				podS.downImage = strings.Split(rs.Spec.Template.Spec.Containers[0].Image, ":")[1]
			}
		}
	}
	podS.downCount = count - podS.upCount
}
