package cli

import (
	"fmt"

	"github.com/rafawastaken/quicktick/internal/storage"
	"github.com/spf13/cobra"
)

func NewLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Log out of the current session",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := storage.ClearSession(); err != nil {
				return fmt.Errorf("failed to clear session: %w", err)
			}
			fmt.Println("Logged out successfully.")
			return nil
		},
	}
}
