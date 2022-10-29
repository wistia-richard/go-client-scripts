package awsutils

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func SpotInstanceScore() int64 {
	val := true
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{CredentialsChainVerboseErrors: &val},
	}))

	id := StsGetCallerId(sess)
	fmt.Println(id)

	ec2svc := ec2.New(sess)

	// filters for instance types and region
	dryrun := false
	m5 := "r5.8xlarge"
	m51 := "r5.16xlarge"
	m52 := "r5.4xlarge"
	m53 := "c5.12xlarge"
	region := "us-east-1"
	var capacity int64 = 1

	options := ec2.GetSpotPlacementScoresInput{
		DryRun:         &dryrun,
		InstanceTypes:  []*string{&m5, &m51, &m52, &m53},
		RegionNames:    []*string{&region},
		TargetCapacity: &capacity,
	}

	fmt.Println(options)

	score, err := ec2svc.GetSpotPlacementScores(&options)
	if err != nil {
		log.Fatal(err)
	}

	return (*score.SpotPlacementScores[0].Score)
}
