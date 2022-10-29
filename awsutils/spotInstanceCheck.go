package awsutils

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func SpotInstanceScore(spotInstanceTypes []string, region string) int64 {
	val := true
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{CredentialsChainVerboseErrors: &val},
	}))

	StsGetCallerId(sess)
	// id := StsGetCallerId(sess)
	// fmt.Println(id)

	ec2svc := ec2.New(sess)

	// filters for instance types and region
	dryrun := false
	var capacity int64 = 1

	instanceTypes := arrToPointerArr(spotInstanceTypes)

	options := ec2.GetSpotPlacementScoresInput{
		DryRun:         &dryrun,
		InstanceTypes:  instanceTypes,
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

func arrToPointerArr(values []string) []*string {
	var pointerArr []*string
	for _, val := range values {
		pointerArr = append(pointerArr, &val)
	}

	return pointerArr

}
