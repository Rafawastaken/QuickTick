package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/rafawastaken/quicktick/internal/app"
	"github.com/rafawastaken/quicktick/internal/domain"
	"github.com/rafawastaken/quicktick/internal/storage"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	f := &Flags{
		Status: string(domain.StatusTodo),
	}

	cmd := &cobra.Command{
		Use:   "quicktick",
		Short: "QuickTick - fast terminal tasks",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !domain.IsValidStatus(domain.Status(f.Status)) {
				return fmt.Errorf("status inválido: %q (usa: todo|progress|completed|canceled)", f.Status)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			// Check Auth
			session, err := storage.LoadSession()
			if err != nil {
				return fmt.Errorf("failed to load session: %w", err)
			}
			if session == nil || session.UserID == "" {
				return fmt.Errorf("You are not logged in. Please run 'qt login' or 'qt signup'.")
			}

			a, err := app.New(ctx, session.UserID)
			if err != nil {
				return err
			}
			defer a.Close()

			action, err := resolveAction(*f)
			if err != nil {
				return err
			}
			if action == ActionNone {
				action = ActionShow
			}

			switch action {
			case ActionShow:
				filter := domain.Filter{}
				if f.Status != "" {
					// We reuse the --status flag, but typically --status is for ADD.
					// Let's check if --status was explicitly set?
					// Cobra sets default values. Status default is "todo".
					// If user runs `qt --show`, status is "todo".
					// We probably want to ignore default "todo" for SHOW unless we add specific filter flag.
					// Or we just abuse "Status" field.
					// Better: add specific filter flag?
					// Plan said: "CLI Flag: qt list --status=todo (or just reusing current --status flag logic)".

					// Issue: "qt --add" uses status default "todo".
					// "qt --show" should probably show ALL by default.
					// But f.Status has "todo" by default.

					// Let's assume if action is SHOW, we only filter if user provided it?
					// Cobra doesn't easily tell us "is default".
					// Let's use `cmd.Flags().Changed("status")`.

					if cmd.Flags().Changed("status") {
						filter.Status = domain.Status(f.Status)
					}
				}

				tasks, err := a.ListTasks(ctx, filter)
				if err != nil {
					return err
				}
				PrintTasks(tasks)
				return nil

			case ActionAdd:
				id, err := a.AddTask(ctx, f.AddText, domain.Status(f.Status))
				if err != nil {
					return err
				}
				fmt.Printf("Task added: [%d] %s\n", id, f.AddText)
				return nil

			case ActionComplete:
				// Fetch task first to show title
				t, err := a.Store.GetTask(ctx, int64(f.CompleteID))
				if err != nil {
					return err
				}

				if err := a.CompleteTask(ctx, int64(f.CompleteID)); err != nil {
					return err
				}
				fmt.Printf("Task completed: [%d] %s\n", f.CompleteID, t.Title)
				return nil

			case ActionOpen:
				if err := a.OpenTask(ctx, int64(f.OpenID)); err != nil {
					return err
				}
				return nil

			case ActionSync:
				if err := a.SyncTasks(ctx); err != nil {
					return err
				}
				fmt.Println("Sync completed")
				return nil

			case ActionEdit:
				if err := a.EditTask(ctx, int64(f.EditID), f.EditText); err != nil {
					return err
				}
				fmt.Printf("Task updated: [%d] %s\n", f.EditID, f.EditText)
				return nil

			case ActionDelete:
				// Fetch task first to show title
				t, err := a.Store.GetTask(ctx, int64(f.DeleteID))
				if err != nil {
					return err
				}

				if err := a.DeleteTask(ctx, int64(f.DeleteID)); err != nil {
					return err
				}
				fmt.Printf("Task deleted: [%d] %s\n", f.DeleteID, t.Title)
				return nil

			default:
				return errors.New("ação desconhecida")
			}
		},
	}

	// Flags (1 comando)
	cmd.Flags().StringVarP(&f.AddText, "add", "a", "", `Add a task. Ex: --add "comprar leite"`)
	cmd.Flags().BoolVarP(&f.Show, "show", "s", false, "Show task list")
	cmd.Flags().IntVarP(&f.CompleteID, "complete", "c", 0, "Complete task by id")
	cmd.Flags().IntVarP(&f.OpenID, "open", "o", 0, "Open task note (markdown) by id")
	cmd.Flags().BoolVar(&f.Sync, "sync", false, "Sync with cloud")
	cmd.Flags().StringVar(&f.Status, "status", string(domain.StatusTodo), "Status for --add: todo|progress|completed|canceled")

	// New Flags
	cmd.Flags().IntVar(&f.EditID, "edit", 0, "Edit task ID")
	cmd.Flags().StringVar(&f.EditText, "title", "", "New title for edit")
	cmd.Flags().IntVar(&f.DeleteID, "rm", 0, "Delete task by ID")

	cmd.AddCommand(NewConfigCmd())
	cmd.AddCommand(NewLoginCmd())
	cmd.AddCommand(NewSignUpCmd())
	cmd.AddCommand(NewLogoutCmd())

	return cmd
}

func resolveAction(f Flags) (Action, error) {
	var actions []Action

	if f.AddText != "" {
		actions = append(actions, ActionAdd)
	}
	if f.Show {
		actions = append(actions, ActionShow)
	}
	if f.CompleteID > 0 {
		actions = append(actions, ActionComplete)
	}
	if f.OpenID > 0 {
		actions = append(actions, ActionOpen)
	}
	if f.Sync {
		actions = append(actions, ActionSync)
	}
	if f.EditID > 0 {
		if f.EditText == "" {
			return ActionNone, errors.New("para editar, usa --edit <ID> --title \"Novo Titulo\"")
		}
		actions = append(actions, ActionEdit)
	}
	if f.DeleteID > 0 {
		actions = append(actions, ActionDelete)
	}

	if len(actions) > 1 {
		return ActionNone, fmt.Errorf("só podes usar 1 ação por comando (encontradas: %v)", actions)
	}
	if len(actions) == 0 {
		return ActionNone, nil
	}
	return actions[0], nil
}
