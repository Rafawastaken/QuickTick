package cli

import (
	"fmt"

	"github.com/rafawastaken/quicktick/internal/config"
	"github.com/rafawastaken/quicktick/internal/storage"
	"github.com/spf13/cobra"
)

func NewConfigCmd() *cobra.Command {
	var url, key string

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			changed := false
			if url != "" {
				cfg.SupabaseURL = url
				changed = true
			}
			if key != "" {
				cfg.SupabaseKey = key
				changed = true
			}

			if changed {
				if err := cfg.Save(); err != nil {
					return err
				}
				fmt.Println("Configuration saved.")
			} else {
				// Don't print secrets anymore. Check for logged in user.
				session, err := storage.LoadSession()
				if err != nil {
					return fmt.Errorf("failed to load session: %w", err)
				}

				if session != nil && session.Email != "" {
					fmt.Printf("Logged in as: %s\n", session.Email)
				} else {
					fmt.Println("Not logged in.")
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&url, "url", "", "Set Supabase URL")
	cmd.Flags().StringVar(&key, "key", "", "Set Supabase Anon Key")

	return cmd
}
