package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/git"
	"github.com/mpppk/hlb/hlblib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "hlb",
	Short: "multi git hosting service manager",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var config etc.Config
		err := viper.Unmarshal(&config)
		etc.PanicIfErrorExist(err)
		remote, err := git.GetDefaultRemote(".")
		etc.PanicIfErrorExist(err)
		serviceConfig, ok := config.FindServiceConfig(remote.ServiceHost)
		if !ok {
			fmt.Println(remote.ServiceHost, "is unknown host. Please add the service configuration to config file(~/.hlb.yaml)")
			os.Exit(1)
		}
		if serviceConfig.Token == "" {
			if !hlblib.CanCreateToken(serviceConfig.Type) {
				fmt.Println("The token of", serviceConfig.Host, "can not create via hlb.")
				fmt.Println("Please add token to config file(~/.hlb.yaml) manually.")
				os.Exit(1)
			}
			serviceUrl := serviceConfig.Protocol + "://" + serviceConfig.Host
			fmt.Println(serviceUrl)
			addServiceCmd.Run(cmd, []string{serviceConfig.Type, serviceUrl})
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hlb.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".hlb") // name of config file (without extension)
	dir, err := homedir.Dir()
	if err != nil {
		etc.PanicIfErrorExist(err)
	}

	viper.AddConfigPath(dir) // adding home directory as first search path
	viper.AutomaticEnv()     // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		initCmd.Run(nil, nil)
		err := viper.ReadInConfig()
		etc.PanicIfErrorExist(err)
	}
}
