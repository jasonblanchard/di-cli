package cmd

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

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

// instancesListCmd represents the instancesList command
var instancesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List worker node instances",
	Run: func(cmd *cobra.Command, args []string) {
		session, err := session.NewSession(&aws.Config{
			Region: aws.String("us-east-1"),
		})

		if err != nil {
			panic(err)
		}

		svc := ec2.New(session)
		input := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name: aws.String("tag:aws:autoscaling:groupName"),
					Values: []*string{
						aws.String(("di-cluster")),
					},
				},
			},
		}
		result, err := svc.DescribeInstances(input)

		if err != nil {
			panic(err)
		}

		writer := tabwriter.NewWriter(os.Stdout, 0, 9, 5, ' ', tabwriter.TabIndent)
		fmt.Fprintln(writer, "ID\tSTATE\tHOST\tIP")

		for _, reservation := range result.Reservations {
			for _, instance := range reservation.Instances {
				id := instance.InstanceId
				state := *instance.State.Name
				var publicHost string
				var publicIP string
				if len(instance.NetworkInterfaces) > 0 {
					publicHost = *instance.NetworkInterfaces[0].Association.PublicDnsName
					publicIP = *instance.NetworkInterfaces[0].Association.PublicIp
				}
				fmt.Fprintf(writer, "%v\t%v\t%v\t%v", *id, state, publicHost, publicIP)
				fmt.Fprint(writer, "\n")
			}
		}

		writer.Flush()
		fmt.Println()
	},
}

func init() {
	instancesCmd.AddCommand(instancesListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// instancesListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// instancesListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
