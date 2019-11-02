package appy

import "os"

func newDbCreateCommand(config *Config, dbManager *DbManager, logger *Logger) *Cmd {
	cmd := &Cmd{
		Use:   "db:create",
		Short: "Creates all databases for the current environment",
		Run: func(cmd *Cmd, args []string) {
			CheckConfig(config, logger)
			CheckDbManager(config, dbManager, logger)

			if len(dbManager.dbs) < 1 {
				logger.Infof("No database is defined in pkg/config/.env.%s", config.AppyEnv)
				os.Exit(0)
			}

			logger.Infof("Creating databases from pkg/config/.env.%s...", config.AppyEnv)

			err := dbManager.ConnectAll(false)
			if err != nil {
				logger.Fatal(err)
			}
			defer dbManager.CloseAll()

			var errs []error
			for _, db := range dbManager.dbs {
				if db.config.Replica {
					continue
				}

				tmpErrs := db.Create()
				errs = append(errs, tmpErrs...)
			}

			if len(errs) > 0 {
				for _, err := range errs {
					logger.Infof(err.Error())
				}

				os.Exit(-1)
			}
		},
	}

	return cmd
}