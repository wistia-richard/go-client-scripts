package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func spotInstanceScore() {
	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	fmt.Println(session)

	svc := ec2.New(session)
	// if err != nil {
	// 	log.Fatal("error creating session")
	// }
	dryrun := false
	m5 := "m5.xlarge"
	region := "us-east-1"
	var capacity int64 = 1

	options := ec2.GetSpotPlacementScoresInput{
		DryRun:         &dryrun,
		InstanceTypes:  []*string{&m5},
		RegionNames:    []*string{&region},
		TargetCapacity: &capacity,
	}

	score, err := svc.GetSpotPlacementScores(&options)
	if err != nil {
		log.Fatal(err)
	}
	print(score)

}
