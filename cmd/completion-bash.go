package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// completionBashCmd represents the completion command
var completionBashCmd = &cobra.Command{
	Use:   "completion-bash",
	Short: "Generates bash completion scripts",
	Long: fmt.Sprintf(`To load completion run

. <(%s completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(%s completion-bash)
`, rootCmd.Name(), rootCmd.Name()),
	Run: func(cmd *cobra.Command, args []string) {
		_ = rootCmd.GenBashCompletion(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(completionBashCmd)
}
