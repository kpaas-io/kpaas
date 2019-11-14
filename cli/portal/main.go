package main

import (
	goflag "flag"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	command := getCommand()
	decorateFlags(command)

	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

func decorateFlags(command *cobra.Command) {

	command.Flags().Int16("port", 8080, "web service listening port")
}

func getCommand() *cobra.Command {

	return &cobra.Command{
		Use:   "portal",
		Short: "portal is a tool for deploying Kubernetes",
		Long:  `portal is a tool for deploying Kubernetes clusters using web interface`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
}
