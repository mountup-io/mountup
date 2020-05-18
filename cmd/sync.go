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
package cmd

import (
	"errors"
	"fmt"
	"github.com/mountup-io/mountup/api"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a clientname")
		}
		if len(args) > 1 {
			return errors.New("too many arguments, only one clientname is required")
		}
		return nil
	},
	Use:   "sync <clientname>",
	Short: "Syncs files in your ~/mountup folder with your remote client",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Printf("syncing on client: %s\n", args[0])
		//
		//s := spinner.New(spinner.CharSets[26], 100*time.Millisecond)
		//s.Start()
		//time.Sleep(1 * time.Second)
		//s.Stop()
		//
		//fmt.Printf("%s sync engaged.\n", args[0])
		//
		//return

		db := api.NewDB()
		host, err := db.GetHostnameForClientname(args[0])
		if err != nil {
			//Either table or client with this name doesn't exist
			fmt.Printf("virtual machine named %s not found.\nTo create a client run:\n\tmountup create %s\n", args[0], args[0])
			return
		}

		homeDir, _ := os.UserHomeDir()
		pkeyDir := filepath.Join(homeDir, ".mountup/keys", args[0]+".pem")

		shellCmd := &exec.Cmd{
			Path:   "/Users/danielwang/go/src/github.com/mountup-io/mountup/cmd/sync.sh",
			Args:   []string{"/Users/danielwang/go/src/github.com/mountup-io/mountup/cmd/sync.sh", "ubuntu@" + host, pkeyDir},
			Stdout: os.Stdout,
			Stderr: os.Stdout,
		}

		err = shellCmd.Start()
		if err != nil {
			return
		}

		err = shellCmd.Wait()
		if err != nil {
			return
		}

		return
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
