// pod-metrics-exporter
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	labelName  *string
	labelValue *string
	listenAddr *string
	kubeconfig *string
)

func init() {
	labelName = flag.String("label-name", "app", "Name of pod label to filter on")
	labelValue = flag.String("label-value", "demo-app", "Value corresponding to the pod label name `--label-name`")
	listenAddr = flag.String("metrics-listen-addr", ":8080", "Listen address of the prometheus metrics server")
	kubeconfig = flag.String("kubeconfig", "kubecfg", "Absolute path to kubeconfig")
}

func main() {
	flag.Parse()

	log.Println("Starting up...")

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	options := metav1.ListOptions{LabelSelector: *labelName + "=" + *labelValue}
	podList, err := client.CoreV1().Pods("").List(context.TODO(), options)
	if err != nil {
		log.Fatal(err.Error())
	}

	podCount := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "pod_count",
			Help:        "Total number of pods with given label",
			ConstLabels: prometheus.Labels{"label_name": *labelName, "label_value": *labelValue},
		},
		[]string{"phase"},
	)

	if err := prometheus.Register(podCount); err != nil {
		log.Fatal("podCount not registered:", err.Error())
	} else {
		log.Println("podCount registered.")
	}

	log.Println("Starting polling loop...")

	go func() {
		for {
			var phase string
			var count float64
			for _, p := range podList.Items {
				phase = string(p.Status.Phase)
				count++
				log.Printf("Pod: %v, Namespace: %v\n", p.Name, p.Namespace)
			}
			podCount.WithLabelValues(phase).Set(count)
			time.Sleep(10 * time.Second)
		}
	}()

	log.Println("Starting metrics handler...")
	log.Println("Ready.")

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
