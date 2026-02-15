package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/nedpals/supabase-go"
	"github.com/rafawastaken/quicktick/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func NewSignUpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signup",
		Short: "Sign up for a new account",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			if cfg.SupabaseURL == "" || cfg.SupabaseKey == "" {
				return fmt.Errorf("Supabase credentials not set. Run 'qt config' first.")
			}

			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Email: ")
			email, _ := reader.ReadString('\n')
			email = strings.TrimSpace(email)

			fmt.Print("Password: ")
			bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return err
			}
			password := string(bytePassword)
			fmt.Println()

			client := supabase.CreateClient(cfg.SupabaseURL, cfg.SupabaseKey)
			ctx := context.Background()

			user, err := client.Auth.SignUp(ctx, supabase.UserCredentials{
				Email:    email,
				Password: password,
			})
			if err != nil {
				return fmt.Errorf("signup failed: %w", err)
			}

			if err != nil {
				return fmt.Errorf("signup failed: %w", err)
			}

			// nedpals/supabase-go SignUp returns *User, not session.
			// We can't save session here.
			// Ideally we would auto-login, but email confirmation might be needed.

			fmt.Printf("Sign up successful for %s!\n", user.Email)
			fmt.Println("Please check your email to confirm your account.")
			fmt.Println("After confirmation, you can run 'qt login'.")

			return nil
		},
	}

	return cmd
}
