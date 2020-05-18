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
	"fmt"
	"github.com/mountup-io/mountup/api"
	"github.com/mountup-io/mountup/util"
	"time"

	"github.com/briandowns/spinner"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <clientname>",
	Short: "Create a virtual machine that is named clientname",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("setting up your ssh keys")
		fmt.Println("provisioning your new virtual machine...")

		s := spinner.New(spinner.CharSets[36], 100*time.Millisecond)
		s.Start()

		vm, pkey, err := api.MakeCreateVMRequest(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		// Save the private key
		err = util.SavePrivateKeyToFS(pkey)
		if err != nil {
			fmt.Printf("Error saving private key: %s\n", err)
			return
		}

		db := api.NewDB()

		err = db.InitTables()
		if err != nil {
			fmt.Printf("Error initializing mountup local db: %s\n", err)
			return
		}

		err = db.PutVM(vm)
		if err != nil {
			fmt.Printf("Error inserting vm into mountup local db: %s\n", err)
			return
		}

		s.Stop()
		fmt.Printf("%s is now ready to roll!\n", args[0])
	},
}

func init() {
	// DISABLED
	//rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
