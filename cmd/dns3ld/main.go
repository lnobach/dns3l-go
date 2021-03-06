package main

import (
	"os"

	"github.com/dta4/dns3l-go/context"
	"github.com/dta4/dns3l-go/service"
	"github.com/dta4/dns3l-go/state"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {

	log.SetLevel(log.DebugLevel)

	err := Execute()
	if err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "dns3ld",
	Short: "dns3l backend daemon",
	Long:  `Foo bar, fill me, version ` + context.Version,
	Run: func(cmd *cobra.Command, args []string) {

		confPath, err := cmd.PersistentFlags().GetString("config")
		if err != nil {
			panic(err)
		}
		conf := service.Config{}
		err = conf.FromFile(confPath)
		if err != nil {
			panic(err)
		}
		err = conf.Initialize()
		if err != nil {
			panic(err)
		}
		socket, err := cmd.PersistentFlags().GetString("socket")
		if err != nil {
			panic(err)
		}
		svc := service.Service{Config: &conf, Socket: socket}
		err = svc.Run()
		if err != nil {
			panic(err)
		}
	},
}

var dbCreateCmd = &cobra.Command{
	Use:   "dbcreate",
	Short: "dns3l backend daemon",
	Long:  `Foo bar, fill me, version ` + context.Version,
	Run: func(cmd *cobra.Command, args []string) {

		confPath, err := cmd.Parent().PersistentFlags().GetString("config")
		if err != nil {
			panic(err)
		}
		conf := service.Config{}
		err = conf.FromFile(confPath)
		if err != nil {
			panic(err)
		}

		state.CreateSQLDB(conf.DB)
	},
}

func Execute() error {
	rootCmd.PersistentFlags().StringP("config", "c", "config.yaml",
		`YAML-formatted configuration for dns3ld.`)
	rootCmd.PersistentFlags().StringP("socket", "s", ":80",
		`L4 socket on which the service should listen.`)
	rootCmd.AddCommand(dbCreateCmd)
	return rootCmd.Execute()
}
