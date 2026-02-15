package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/nedpals/supabase-go"
	"github.com/rafawastaken/quicktick/internal/config"
	"github.com/rafawastaken/quicktick/internal/storage"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func NewLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to Supabase",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			if cfg.SupabaseURL == "" || cfg.SupabaseKey == "" {
				return fmt.Errorf("Supabase credentials not set. Run 'qt config --url <URL> --key <KEY>' first.")
			}

			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Email: ")
			email, _ := reader.ReadString('\n')
			email = strings.TrimSpace(email)

			fmt.Print("Password: ")
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return err
			}
			password := string(bytePassword)
			fmt.Println() // Newline after password input

			client := supabase.CreateClient(cfg.SupabaseURL, cfg.SupabaseKey)
			ctx := context.Background()

			user, err := client.Auth.SignIn(ctx, supabase.UserCredentials{
				Email:    email,
				Password: password,
			})
			if err != nil {
				return fmt.Errorf("login failed: %w", err)
			}

			// Clear old config token if exists (migration)
			cfg.Token = ""
			_ = cfg.Save()

			// Save new session
			session := &storage.Session{
				UserID:      user.User.ID,
				AccessToken: user.AccessToken,
				Email:       user.User.Email,
			}
			if err := storage.SaveSession(session); err != nil {
				return fmt.Errorf("failed to save session: %w", err)
			}

			fmt.Println("Login successful! Session saved.")
			return nil
		},
	}

	return cmd
}
