package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/rafawastaken/quicktick/internal/domain"
)

func PrintTasks(tasks []domain.Task) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tSTATUS\tCREATED\tTITLE")

	for _, t := range tasks {
		status := t.Status
		// Simple ANSI colors (manual for now, can use library later if needed)
		// Reset: \033[0m
		// Green: \033[32m
		// Yellow: \033[33m
		// Red: \033[31m
		// Blue: \033[34m

		color := ""
		switch status {
		case domain.StatusCompleted:
			color = "\033[32m" // Green
		case domain.StatusTodo:
			color = "\033[33m" // Yellow
		case domain.StatusProgress:
			color = "\033[34m" // Blue
		case domain.StatusCanceled:
			color = "\033[31m" // Red
		}

		created := t.CreatedAt.Format("2006-01-02 15:04")

		if color != "" {
			fmt.Fprintf(w, "%d\t%s%s\033[0m\t%s\t%s\n", t.ID, color, status, created, t.Title)
		} else {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", t.ID, status, created, t.Title)
		}
	}
	_ = w.Flush()
}
