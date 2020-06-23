/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"
	"github.com/spf13/cobra"
)

// instancesSendKeyCmd represents the instancesSendKey command
var instancesSendKeyCmd = &cobra.Command{
	Use:   "send-key",
	Short: "Send SSH public key to specified instance",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetString("instance-id")
		if err != nil {
			panic(err)
		}

		session := ec2instanceconnect.New(session.New(&aws.Config{
			Region: aws.String("us-east-1"),
		}))

		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		pkeydir := fmt.Sprintf("%s/.ssh/id_rsa.pub", home)
		pkey, err := ioutil.ReadFile(pkeydir)
		if err != nil {
			panic(err)
		}
		input := &ec2instanceconnect.SendSSHPublicKeyInput{
			AvailabilityZone: aws.String("us-east-1a"),
			InstanceId:       aws.String(id),
			InstanceOSUser:   aws.String("ubuntu"),
			SSHPublicKey:     aws.String(string(pkey)),
		}

		result, err := session.SendSSHPublicKey(input)
		if err != nil {
			panic(err)
		}
		fmt.Println(result)
	},
}

func init() {
	instancesCmd.AddCommand(instancesSendKeyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// instancesSendKeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	instancesSendKeyCmd.Flags().StringP("instance-id", "i", "", "AWS instance ID")
}
