/*
Copyright © 2023 dengxiaoshan
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"transfer/internal/app"
	"transfer/internal/app/config"
)

var cfgFile string

var cfg config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "transfer",
	Short: "transfer file to minio API",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			fmt.Println("configure file is not exist.", err)
			os.Exit(1)
		}
		app.Run(ctx)

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.transfer.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".transfer" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".transfer")
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.WatchConfig()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
		// Unmarshal 将配置文件转成对象
		if cfgErr := viper.Unmarshal(&cfg); cfgErr != nil {
			fmt.Println("parse config file error.", cfgErr)
			os.Exit(1)
		}
	}
}
