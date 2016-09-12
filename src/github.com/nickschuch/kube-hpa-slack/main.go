package main

import (
	"fmt"
	"log"
	"time"

	slack "github.com/nickschuch/go-slack"
	"gopkg.in/alecthomas/kingpin.v2"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/autoscaling"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
)

var (
	cliKubernetes = kingpin.Flag("kubernetes", "Kubernetes API server").Default("http://localhost:8080").String()

	cliSlackEmoji   = kingpin.Flag("slack-emoji", "Slack - Emoji").Default(":slack:").String()
	cliSlackChannel = kingpin.Flag("slack-channel", "Slack - Channel").Default("general").String()
	cliSlackUrl     = kingpin.Flag("slack-url", "Slack - Url").Required().String()
)

func main() {
	kingpin.Parse()

	prevHPAs := make(map[string]autoscaling.HorizontalPodAutoscaler)

	throttle := time.Tick(time.Second / 60)

	for {
		<-throttle

		k8s, err := client.New(&restclient.Config{
			Host: *cliKubernetes,
		})
		if err != nil {
			log.Println(err)
			continue
		}

		ns, err := k8s.Namespaces().List(api.ListOptions{})
		if err != nil {
			log.Println(err)
			continue
		}

		for _, n := range ns.Items {
			hpas, err := k8s.Autoscaling().HorizontalPodAutoscalers(n.ObjectMeta.Name).List(api.ListOptions{})
			if err != nil {
				log.Println(err)
				continue
			}

			for _, h := range hpas.Items {
				id := n.ObjectMeta.Name + "/" + h.ObjectMeta.Name

				var msg string

				// Ensure we have a previously set value.
				if _, ok := prevHPAs[id]; !ok {
					prevHPAs[id] = h
					continue
				}

				// Check if the values have been changed.
				if h.Spec.MinReplicas != prevHPAs[id].Spec.MinReplicas {
					msg = msg + fmt.Sprintf("Minimum changed from *%d* to *%d*\n", prevHPAs[id].Spec.MinReplicas, h.Spec.MinReplicas)
				}

				if h.Spec.MaxReplicas != prevHPAs[id].Spec.MaxReplicas {
					msg = msg + fmt.Sprintf("Maximum changed from *%d* to *%d*\n", prevHPAs[id].Spec.MaxReplicas, h.Spec.MaxReplicas)
				}

				if h.Status.CurrentReplicas != prevHPAs[id].Status.CurrentReplicas {
					msg = msg + fmt.Sprintf("Current changed from *%d* to *%d*\n", prevHPAs[id].Status.CurrentReplicas, h.Status.CurrentReplicas)
				}

				if h.Status.DesiredReplicas != prevHPAs[id].Status.DesiredReplicas {
					msg = msg + fmt.Sprintf("Desired changed from *%d* to *%d*\n", prevHPAs[id].Status.DesiredReplicas, h.Status.DesiredReplicas)
				}

				// Save it for later.
				prevHPAs[id] = h

				if msg != "" {
					if err := slack.Send(id, *cliSlackEmoji, msg, *cliSlackChannel, *cliSlackUrl); err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}
