package internals

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type appconfig struct {
	ClusterName            string   `json:"ClusterName"`
	SpotNodeTypes          []string `json:"SpotNodeTypes"`
	OnDemandNodeGroups     []string `json:"OnDemandNodeGroups"`
	AppRerunInterval       string   `json:"RerunInterval"`
	MaxConcurrentNodeDrain int      `json:"MaxConcurrentNodeDrain"`
	DrainPauseDuration     string   `json:"DrainPauseDuration"`
	Ec2SpotScoreRequired   int      `json:"Ec2SpotScoreRequired"`
	DryRun                 bool     `json:"DryRun"`
}

func GenerateJsonTemplate(name string) {
	f := createFile(name)
	fmt.Printf("The file \"%s\" was created\n", f.Name())

	config := appconfig{}
	jsonbyte, err := json.MarshalIndent(config, "", "   ")
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println(string(jsonbyte))
	w := bufio.NewWriter(f)
	n4, err := w.WriteString(string(jsonbyte))
	checkerr(err)
	fmt.Printf("wrote %d bytes\n", n4)

	w.Flush()
}

func GetJsonConfig(name string) {
	data, err := os.ReadFile(name)
	checkerr(err)
	config := appconfig{}
	err = json.Unmarshal(data, &config)
	checkerr(err)

	fmt.Println(config)

}

func createFile(name string) *os.File {
	f, err := os.Create(name)
	if err != nil {
		log.Fatalf(err.Error())
	} else {
		fmt.Printf("The file \"%s\" was created\n", f.Name())
	}
	return f
}

func checkerr(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}
