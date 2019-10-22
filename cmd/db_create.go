package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/appist/appy/core"
)

// NewDbCreateCommand creates all databases from app/config/.env.<APPY_ENV>.
func NewDbCreateCommand(config core.AppConfig, dbMap map[string]*core.AppDb) *AppCmd {
	cmd := &AppCmd{
		Use:   "db:create",
		Short: "Creates all databases from app/config/.env.<APPY_ENV>.",
		Run: func(cmd *AppCmd, args []string) {
			fmt.Printf("Creating databases from app/config/.env.%s...\n", config.AppyEnv)

			if len(dbMap) < 1 {
				fmt.Printf("No database is defined in app/config/.env.%s.\n", config.AppyEnv)
				os.Exit(-1)
			}

			var msgs, errs []string
			for _, db := range dbMap {
				if db.Config.Replica {
					continue
				}

				tmpMsgs, tmpErrs := dbCreate(db)
				msgs = append(msgs, tmpMsgs...)
				errs = append(errs, tmpErrs...)
			}

			if len(errs) > 0 {
				fmt.Println(strings.Join(errs, "\n"))
				os.Exit(-1)
			}

			for _, msg := range msgs {
				fmt.Println(msg)
			}
		},
	}

	return cmd
}