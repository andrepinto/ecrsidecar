package main

import (
	"encoding/base64"
	"flag"
	"os"
	"strings"
	"time"
	"io/ioutil"
	"log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"gopkg.in/yaml.v2"

)

type EcrData struct {
	Interval           string `yaml:"interval"`
	AwsAccessKeyId     string `yaml:"accessKeyId"`
	AwsSecretAccessKey string `yaml:"secretAccessKey"`
	Region             string `yaml:"region"`
	TokenFile          string `yaml:"tokenFile"`
	RegistryID 	   string `yaml:"registryID"`
}



func main() {

	var file = flag.String("config", "config.yaml", "Ecr configuration file")

	flag.Parse()

	data, err := ioutil.ReadFile(*file)

	if err != nil {
		log.Fatal(err)
	}

	var ecrData EcrData
	err = yaml.Unmarshal([]byte(data), &ecrData)
	if err != nil {
		log.Fatalf("error: %v", err)
	}


	os.Setenv("AWS_ACCESS_KEY_ID", ecrData.AwsAccessKeyId)
	os.Setenv("AWS_SECRET_ACCESS_KEY", ecrData.AwsSecretAccessKey)

	sess := session.New(
		aws.NewConfig().WithMaxRetries(10).WithRegion(ecrData.Region),
	)

	interval, err := time.ParseDuration(ecrData.Interval)

	ecrCredential(ecrData, sess)

	c := time.Tick(interval)

	for range c {
		log.Println("refreshing...")
		go ecrCredential(ecrData, sess)
	}



}

func ecrCredential(ecrData EcrData, sess *session.Session) {
	instance := ecr.New(sess, &aws.Config{
		Region: aws.String(ecrData.Region),
	})

	input := &ecr.GetAuthorizationTokenInput{
		RegistryIds: []*string{aws.String(ecrData.RegistryID)},
	}

	authToken, err := instance.GetAuthorizationToken(input)

	if err != nil {
		log.Printf("Failed to get credential for %s in region %s (%s)", ecrData.RegistryID, ecrData.Region, err.Error())
		return
	}

	for _, data := range authToken.AuthorizationData {
		output, err := base64.StdEncoding.DecodeString(*data.AuthorizationToken)

		if err != nil {
			log.Printf("Failed to decode credential for %s in region %s (%s)", ecrData.RegistryID, ecrData.Region, err.Error())
			return
		}

		var password string
		split := strings.Split(string(output), ":")

		if len(split) == 2 {
			password = strings.TrimSpace(split[1])
		} else {
			log.Print("Failed to parse password.")
			return
		}
		err = ioutil.WriteFile(ecrData.TokenFile, []byte(password), 0644)

		if err != nil {
			log.Printf("Failed to write password to file %s (%s)", ecrData.TokenFile, err.Error())
		}
	}
}


