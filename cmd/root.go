package cmd

import (
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mountup",
	Short: "mountup command line client",
	Long: `mountup is a code syncing tool used with remote machines.
Write code remotely from the comfort of your own IDEs with local performance.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hi! Welcome to the faster way to develop on remote machines!\nPlease take a look at --help for more info.")
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
	//cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mountup.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	homeDir, _ := os.UserHomeDir()
	mountDir := filepath.Join(homeDir, "mountup")
	keyDir := filepath.Join(homeDir, ".mountup/keys")
	err := os.MkdirAll(keyDir, 0755)

	if err != nil {
		fmt.Printf("failed to init ~/.mountup/ folder %s\n", err)
	}

	err = os.MkdirAll(mountDir, 0755)
	if err != nil {
		fmt.Printf("failed to init ~/mountup/ folder %s\n", err)
	}
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

		// Search config in home directory with name ".mountup" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".mountup")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
