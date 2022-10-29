package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ricdon41/go-client-test/awsutils"
	"github.com/ricdon41/go-client-test/internals"
	"github.com/ricdon41/go-client-test/kubeutils"

	api "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {

	config := internals.GetJsonConfig("test_config1")

	clientset := kubeutils.KubeClient()

	// filter nodes using label selector
	nodegroups := config.OnDemandNodeGroups
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
				fmt.Printf("List of filtered on demand instance contain one of the labels in %v:\n", config.OnDemandNodeGroups)
				print(node.Name)
				for _, condition := range node.Status.Conditions {
					if condition.Type == "Ready" && condition.Status == "True" {
						// append to the slice
						nodesToDrain = append(nodesToDrain, node)
					}
				}
			}
		}

		if len(nodesToDrain) == 0 {
			fmt.Println("No matching nodes found")
		}

		taintValue := api.Taint{Key: "key1", Value: "value1", Effect: "NoExecute"}
		taintexists := false
		var nodeStatus api.ConditionStatus

		for i, nodeToDrain := range nodesToDrain {

			// get a node from nodesToDrain using the name
			apinode, _ := clientset.CoreV1().Nodes().Get(context.TODO(), nodeToDrain.Name, v1.GetOptions{})
			print(apinode.Spec)

			for _, condition := range apinode.Status.Conditions {
				if condition.Type == "Ready" {
					nodeStatus = condition.Status
				}
			}
			if nodeStatus != "True" {
				continue
			}

			score := awsutils.SpotInstanceScore(config.SpotNodeTypes, config.AWSRegion)
			fmt.Printf("The Spot placement score is : %v\n", score)

			if score < int64(config.Ec2SpotScoreRequired) {
				fmt.Printf("Cancelling current run since spot placement score (%v) is lesser than expected %v\n", score, config.Ec2SpotScoreRequired)
				break
			}

			// if score is greater than required, drain one node at a time
			if nodeStatus == "True" {
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

			// Add a wait time between node drain
			durations, _ := time.ParseDuration(config.DrainPauseDuration)
			fmt.Printf("Pause for %v before draining more nodes", durations)
			time.Sleep(durations)

			if i == len(nodesToDrain)-1 {
				fmt.Println("All nodes were processed")
			}
		}
		// Add a wait time between script re-run
		durations, _ := time.ParseDuration(config.AppRerunInterval)
		fmt.Printf("Sleep until next rerun in %v ", durations)
		time.Sleep(durations)
	}

}

func print(value interface{}) {
	fmt.Println(value)
}

// pending:

//Send statsd to datadog
// Set up monitoring
// setup ci
// in cluster setup
// rbac setup for the cluster.
