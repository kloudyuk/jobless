package main

import (
	"context"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func init() {

	// Set log level
	lvlStr, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		lvlStr = "INFO"
	}
	lvl, err := log.ParseLevel(lvlStr)
	if err != nil {
		panic(err)
	}
	log.SetLevel(lvl)
}

func main() {

	ctx := context.Background()

	log.Info("Getting in-cluster client config")
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Getting current namespace")
	b, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		log.Fatal(err)
	}
	namespace := string(b)
	log.Infof("Namespace: %s", namespace)

	log.Info("Creating client")
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Creating job client for namespace: %s", namespace)
	jobs := clientset.BatchV1().Jobs(namespace)

	log.Info("Setting up watch for successful jobs")
	w, err := jobs.Watch(ctx, metav1.ListOptions{
		FieldSelector: "status.successful=1",
	})
	if err != nil {
		log.Fatal(err)
	}

	propagationPolicy := metav1.DeletePropagationBackground

	watchc := w.ResultChan()

	// The watch will timeout so run it in a loop to ensure it's restarted
	for {
		// read events from the channel
		e, ok := <-watchc
		// break the loop if the channel has been closed
		if !ok {
			log.Info("Watch channel closed")
			break
		}
		// skip events that are not "ADDED" events
		// the fieldSelector on the watch means we'll get
		// "ADDED" events when a job completes
		if e.Type != watch.Added {
			continue
		}
		// Cast the object the job type
		job, ok := e.Object.(*batchv1.Job)
		if !ok {
			continue
		}
		// Delete the job
		log.Infof("Deleting job: %s", job.Name)
		if err := jobs.Delete(ctx, job.Name, metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		}); err != nil {
			log.Fatal(err)
		}
	}

}
