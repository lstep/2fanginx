package server

import (
	"2fanginx/database"

	"github.com/spf13/cobra"
)

func CreateDB(cmd *cobra.Command, args []string) {
	database.CreateDB()
}
