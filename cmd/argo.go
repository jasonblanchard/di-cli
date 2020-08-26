package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

var argoCmd = &cobra.Command{
	Use:   "argo",
	Short: "Port forward to Argo and show login info",
	RunE: func(cmd *cobra.Command, args []string) error {
		podName, err := exec.Command("kubectl", "get", "pods", "-n=argocd", "-l=app.kubernetes.io/name=argocd-server", "-o=name").Output()
		if err != nil {
			return err
		}
		password := strings.TrimSpace(string(podName))

		fmt.Println("")
		fmt.Println("Username: admin")
		fmt.Println(fmt.Sprintf("Password: %s", password))

		portForwardCmd := exec.Command("kubectl", "port-forward", "svc/argocd-server", "-n=argocd", "8080:443")

		err = portForwardCmd.Start()
		if err != nil {
			return err
		}
		defer portForwardCmd.Process.Signal(syscall.SIGTERM)

		fmt.Println("")
		fmt.Println("Argo running at https://localhost:8080/")
		fmt.Println("")

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		for range c {
			fmt.Println("")
			fmt.Println("Received SIGINT, cleaning up...")
			fmt.Println("")
			return nil
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(argoCmd)
}
