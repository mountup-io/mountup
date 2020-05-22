package cmd

import (
	"errors"
	"fmt"
	"github.com/mountup-io/mountup/api"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a clientname")
		}
		stringSlice := strings.FieldsFunc(args[0], Split)
		if len(stringSlice) != 1 && len(stringSlice) != 3 {
			return errors.New("invalid syntax, see `mountup help sync` for more details")
		}
		if len(args) > 3 {
			return errors.New("too many arguments")
		}
		if args[0] != "push" && args[0] != "pull" {
			return errors.New("first arg must be either `push` or `pull`")
		}
		return nil
	},
	Use:   "sync <push/pull> <servername>:<directory_on_remote>",
	Short: "Syncs files in your ~/mountup/servername folder with your remote server",
	Long: `sync <push/pull> <servername>
syncs ~/mountup/servername directory with your mountup instances or own servers

sync <push/pull> username@remote_host:directory_on_remote <ssh_key_path>
syncs with your own servers

sync push will upload files from ~/mountup/<servername> to the remote server before local syncing

sync pull downloads files from the remote server to ~/mountup/<servername>
`,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse input
		// Tokenize
		stringSlice := strings.FieldsFunc(args[1], Split)

		var servername string
		var username string
		var host string
		var destDir string
		var pkeyDir string

		if len(stringSlice) == 1 || len(stringSlice) == 2 {
			username = "ubuntu"
			servername = stringSlice[0]

			if len(stringSlice) == 1 {
				destDir = "~"
			} else {
				destDir = stringSlice[1]
			}

			db := api.NewDB()

			var err error
			host, err = db.GetHostnameForClientname(servername)
			if err != nil {
				//Either table or client with this name doesn't exist
				fmt.Printf("virtual machine named %s not found.\nTo create a client run:\n\tmountup create %s\n", servername, servername)
				return
			}

			homeDir, _ := os.UserHomeDir()
			pkeyDir = filepath.Join(homeDir, ".mountup/keys", servername+".pem")

		} else if len(stringSlice) == 3 {
			username = stringSlice[0]
			host = stringSlice[1]
			destDir = stringSlice[2]
			pkeyDir = args[1]
			servername = host
		}

		shellCmd := &exec.Cmd{
			Path: "/Users/danielwang/go/src/github.com/mountup-io/mountup/cmd/sync.sh",
			Args: []string{
				"/Users/danielwang/go/src/github.com/mountup-io/mountup/cmd/sync.sh",
				username + "@" + host,
				destDir,
				pkeyDir,
				servername,
				args[0],
			},
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}

		err := shellCmd.Start()
		if err != nil {
			return
		}

		err = shellCmd.Wait()
		if err != nil {
			return
		}

		fmt.Println("Hi this exited")

		return
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}

func Split(r rune) bool {
	return r == ':' || r == '@'
}
