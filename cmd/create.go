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
	Short: "creates a new mountup provisioned virtual machine as your remote server",
	Long:  `creates a new mountup provisioned virtual machine as your remote server named <servername>`,
	Run: func(cmd *cobra.Command, args []string) {

		s := spinner.New(spinner.CharSets[36], 100*time.Millisecond)
		s.Start()

		fmt.Println("provisioning your new virtual machine...")

		vm, pkey, err := api.MakeCreateVMRequest(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("setting up ssh keys")

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
}
