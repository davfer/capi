package cmd

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/davfer/capi/config"
	"github.com/davfer/capi/internal/datum"
	"github.com/davfer/capi/internal/finder"
	"github.com/davfer/capi/internal/repository"
	"github.com/davfer/capi/internal/requester"
	"github.com/davfer/capi/pkg/model"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

const (
	CfgFileDefaultPath      = ".capi.yaml"
	VerboseDefaultValue     = false
	DryRunDefaultValue      = false
	EnvironmentDefaultValue = "default"
)

var cfgFile string
var verbose bool
var dryRun bool
var env string
var confirmDialog bool
var diogenesMode bool

var CapiConfig config.CapiConfiguration

var rootCmd = &cobra.Command{
	Use:   "capi [resource]",
	Short: "CApi is a tool to perform API calls rapidly.",
	Long:  `Quick-to-execute CLI API tool that avoids the need to use a GUI and don't need to remember/copy&paste the API calls and play heavily with curl.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		localRepo := repository.NewLocalRepository(CapiConfig.Collections)
		capiFinder := finder.NewSimpleFinder(localRepo)
		capiData := datum.NewSystem()
		capiBuilder := requester.NewHttpBuilderResolver(capiData)
		capiRequester := requester.NewSimpleRequester(capiBuilder)

		var c *model.Collection
		c, err := capiFinder.FindPotentialCollection(args[0])
		if err != nil {
			logrus.WithError(err).WithField("collection", args[0]).Error("error while finding collection")
			os.Exit(1)
		}
		if c == nil {
			if capiFinder.IsANewPotentialCollection(args[0]) {
				// TODO v2 check if available in remote repository
				// TODO create a temporal collection with the base URL and a name
			} else {
				fmt.Printf("Collection '%s' not found, use `capi new %s` to create a new one or use -n flag", args[0], args[0])
				os.Exit(1)
			}
		}

		logrus.WithField("collection", c).Info("found collection")
		r, err := capiFinder.FindPotentialResource(c, args)
		if err != nil {
			logrus.WithError(err).WithField("collection", args[0]).Error("error while finding resource")
			os.Exit(1)
		}
		if r == nil {
			if capiFinder.IsANewPotentialResource(c, args) {
				// TODO create a resource in the collection
			} else {
				fmt.Printf("Resource '%s' not found, use `capi new %s` to create a new one or use -n flag", args[0], args[0])
				os.Exit(1)
			}
		}

		logrus.WithField("collection", c).WithField("resource", r).Info("found resource")
		req := capiRequester.New(c, r)
		if !req.IsComplete() {
			// TODO ask for missing params
		}

		if confirmDialog {
			// TODO ask for confirmation
			fmt.Println("TODO ask for confirmation")
		}

		if diogenesMode {
			// TODO store collection and resource in config
			fmt.Println("TODO store collection and resource in config")
		}

		if !dryRun {
			execute, err := req.Execute()
			if err != nil {
				logrus.WithError(err).Error("error while executing request")
				os.Exit(1)
			}

			fmt.Println(execute)
		}

		os.Exit(0)
	},
}

func init() {
	/**
	-d : dry run
	-y : no-interaction
	-n : create if no exists
	-v : debug, verbose
	-e : environment
	*/
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", CfgFileDefaultPath, "config file (default is $HOME/.capi.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", VerboseDefaultValue, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", DryRunDefaultValue, "dry run: outputs the intended calls but does not execute them")
	rootCmd.PersistentFlags().StringVarP(&env, "env", "e", EnvironmentDefaultValue, "environment: environment running the collections")
	viper.SetDefault("author", "davfer <iam@davfer.com>")
	viper.SetDefault("license", "apache")
}

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if cfgFile != CfgFileDefaultPath {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(home)
		viper.SetConfigName(".capi")
	}

	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// TODO: Create empty config file
		} else {
			fmt.Printf("Error reading config file, %s", err)
			os.Exit(1)
		}
	}

	err = viper.Unmarshal(&CapiConfig)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
		os.Exit(1)
	}

	spew.Dump("Config Read", CapiConfig)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
