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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/cobra"
)

// kubeconfigCmd represents the kubeconfig command
var kubeconfigCmd = &cobra.Command{
	Use:   "kubeconfig",
	Short: "Download kubeconfig and export it to KUBECONFIG env",
	Run: func(cmd *cobra.Command, args []string) {
		session := session.Must(session.NewSession(&aws.Config{
			Region: aws.String("us-east-1"),
		}))
		downloader := s3manager.NewDownloader(session)

		homedir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		kubedir := fmt.Sprintf("%s/.kube", homedir)
		filename := fmt.Sprintf("%s/kubeconfig-di-aws", kubedir)
		f, err := os.Create(filename)
		if err != nil {
			panic(err)
		}

		_, err = downloader.Download(f, &s3.GetObjectInput{
			Bucket: aws.String("di-kubeconfig"),
			Key:    aws.String("kubeconfig"),
		})

		if err != nil {
			panic(err)
		}

		fmt.Println(fmt.Sprintf("export KUBECONFIG=%s", filename))
		fmt.Println("# Run this command to configure your shell:")
		fmt.Println("# eval \"$(di kubeconfig)\"")
	},
}

func init() {
	rootCmd.AddCommand(kubeconfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kubeconfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kubeconfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
