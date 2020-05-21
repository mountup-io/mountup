package cmd

import (
	"fmt"
	"github.com/mountup-io/mountup/api"
	"github.com/mountup-io/mountup/util"
	"os"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <servername>",
	Short: "Create mountup server that is named clientname",
	Long: `Create mountup server that is named clientname
create <clientname>
creates a new mountup provisioned virtual machine as your remote server
`,
	//create username@remote_host:destination_directory <optional_ssh_key_path>
	//creates a client that connects to an existing remote server
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

		homeDir, _ := os.UserHomeDir()
		serverDir := filepath.Join(homeDir, "mountup", args[0])
		err = os.MkdirAll(serverDir, 0755)

		s.Stop()
		fmt.Printf("%s is now ready to roll!\n", args[0])
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
