package cmd

import (
	"fmt"
	"io"
	"strings"

	"curriculum/internal/repository"

	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List skills available in the central repository",
		Long:  "Lists all skills and their versions in ~/.curriculum/repository/.",
		RunE:  runList,
	}
}

func runList(cmd *cobra.Command, args []string) error {
	infos, err := repository.List()
	if err != nil {
		return err
	}

	output(cmd, infos, func(w io.Writer) {
		if len(infos) == 0 {
			fmt.Fprintln(w, "No skills in repository.")
			return
		}
		rows := make([][]string, len(infos))
		for i, info := range infos {
			// Mark the latest version with a *.
			marked := make([]string, len(info.Versions))
			for j, v := range info.Versions {
				if v == info.Latest {
					marked[j] = v + " *"
				} else {
					marked[j] = v
				}
			}
			rows[i] = []string{info.Name, strings.Join(marked, ", ")}
		}
		printTable(w, []string{"SKILL", "VERSIONS"}, rows)
	})
	return nil
}
