package main

import (
	"log"
	"os"

	sample_grpc_server "sample-grpc-server"
	"sample-grpc-server/cmd/migration"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	SilenceUsage:  false,
	SilenceErrors: false,
	Run: func(cmd *cobra.Command, args []string) {
		sample_grpc_server.Start()
	},
}

func init() {
	rootCmd.AddCommand(migration.Cmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Println("execute command:", err)
		os.Exit(1)
	}
}
