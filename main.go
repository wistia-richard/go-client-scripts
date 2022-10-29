package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ricdon41/go-client-test/awsutils"
	"github.com/ricdon41/go-client-test/internals"
	"github.com/ricdon41/go-client-test/kubeutils"

	api "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {

	internals.GetJsonConfig("test_config1")

	clientset := kubeutils.KubeClient()

	nodeList, _ := clientset.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{})
	nodes := nodeList.Items

	f, _ := os.Create("test_nodeData")
	w := bufio.NewWriter(f)

	// get the node names and ready status
	for _, node := range nodes {
		print(node.Name)
		for _, condition := range node.Status.Conditions {
			if condition.Type == "Ready" {
				print(condition.Status)
			}
		}

		// example of spec
		jsonblob, err := json.MarshalIndent(node, "", "   ")
		n4, err := w.WriteString(string(jsonblob))
		_, _ = w.WriteString("\n")
		print(err)
		fmt.Printf("wrote %d bytes\n", n4)
		w.Flush()
	}

	score := awsutils.SpotInstanceScore()
	print(score)

	// filter nodes using label selector
	nodegroups := []string{"generic-1d-od-1"}
	nodesToDrain := []api.Node{}
	for {

		for _, nodegroup := range nodegroups {
			labelFilter := v1.LabelSelector{
				MatchLabels: map[string]string{"alpha.eksctl.io/nodegroup-name": nodegroup},
			}

			labelString := v1.FormatLabelSelector(&labelFilter)
			nodeListFiltered, _ := clientset.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{LabelSelector: labelString})
			nodesFiltered := nodeListFiltered.Items

			// get the node names and ready status
			for _, node := range nodesFiltered {
				print(node.Name)
				for _, condition := range node.Status.Conditions {
					if condition.Type == "Ready" {
						print(condition.Status)
					}
				}

				// append to the slice
				nodesToDrain = append(nodesToDrain, node)
			}
		}

		if len(nodesToDrain) == 0 {
			log.Fatal("no nodes found")
		}

		// print filtered nodes data to a different file
		f, _ = os.Create("test_filteredNodesData")
		w = bufio.NewWriter(f)
		jsonblob, err := json.MarshalIndent(nodesToDrain, "", "   ")

		n4, err := w.WriteString(string(jsonblob))
		_, _ = w.WriteString("\n")
		print(err)
		fmt.Printf("wrote %d bytes\n", n4)
		w.Flush()

		taintValue := api.Taint{Key: "key1", Value: "value1", Effect: "NoExecute"}
		taintexists := false

		for _, nodeToDrain := range nodesToDrain {

			// get a node from nodesToDrain using the name
			apinode, _ := clientset.CoreV1().Nodes().Get(context.TODO(), nodeToDrain.Name, v1.GetOptions{})
			print(apinode.Spec)

			// if score is greater than 7, drain one node at a time
			if score > 7 {
				for _, taint := range apinode.Spec.Taints {
					if taint.Key == taintValue.Key && taint.Value == taintValue.Value && taint.Effect == taintValue.Effect {
						fmt.Print("Taint already exists")
						taintexists = true
						break
					}
				}
				if !taintexists {
					fmt.Println("Taint doesn't exist")
					apinode.Spec.Taints = append(apinode.Spec.Taints, taintValue)

					// do  a dry run

					apinodeUpdated, err := clientset.CoreV1().Nodes().Update(context.TODO(), apinode, v1.UpdateOptions{DryRun: []string{"All"}, FieldManager: "", FieldValidation: "strict"})
					print(apinodeUpdated.Spec)
					if err != nil {
						print(err)
					}

				}
			}

			durations, _ := time.ParseDuration("10s")

			// Add a wait time between node drain
			time.Sleep(durations)
		}

		fmt.Println("all nodes were drained")
		durations, _ := time.ParseDuration("1h")

		// Add a wait time between node drain
		time.Sleep(durations)
	}
}

func print(value interface{}) {
	fmt.Println(value)
}
