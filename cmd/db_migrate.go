package cmd

import (
	"github.com/appist/appy/core"
)

// NewDbMigrateCommand migrates all the databases(or specific database with --database) for the current environment.
func NewDbMigrateCommand(config core.AppConfig, dbMap map[string]*core.AppDb) *AppCmd {
	var target string

	cmd := &AppCmd{
		Use:   "db:migrate",
		Short: "Migrates all the databases(or specific database with --database) for the current environment.",
		Run: func(cmd *AppCmd, args []string) {
			logger.Infof("Migrating databases from app/config/.env.%s...", config.AppyEnv)

			err := core.DbConnect(dbMap, logger, true)
			if err != nil {
				logger.Fatal(err)
			}
			defer core.DbClose(dbMap)

			if target != "" {
				return
			}

			for name, db := range dbMap {
				logger.Infof("Migrating '%s' database...", name)

				err := db.Migrate()
				if err != nil {
					logger.Fatal(err)
				}
			}
		},
	}

	cmd.Flags().StringVar(&target, "database", "", "The target database to migrate.")
	return cmd
}
