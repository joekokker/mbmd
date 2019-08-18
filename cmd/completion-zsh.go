package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// completionZshCmd represents the completion command
var completionZshCmd = &cobra.Command{
	Use:   "completion-zsh",
	Short: "Generates zsh completion scripts",
	Long: fmt.Sprintf(`To load completion run

. <(%s completion-zsh)

To configure your zsh shell to load completions for each session add to your zshrc

# ~/.zshrc
. <(%s completion-zsh)
`, rootCmd.Name(), rootCmd.Name()),
	Run: func(cmd *cobra.Command, args []string) {
		_ = rootCmd.GenZshCompletion(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(completionZshCmd)
}
