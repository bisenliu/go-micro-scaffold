package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"

	commonDI "common/di"
	"services/internal/infrastructure/persistence/ent"
	"services/internal/infrastructure/persistence/ent/gen"
	"services/internal/infrastructure/persistence/ent/gen/migrate"
)

func main() {
	// 创建CLI应用
	app := fx.New(
		fx.Options(
			commonDI.ConfigModule,
			commonDI.LoggerModule,
			commonDI.DatabasesModule,

			ent.Module,
		),

		// CLI入口点
		fx.Invoke(runCLI),
	)

	if err := app.Start(context.Background()); err != nil {
		fmt.Printf("Failed to start CLI application: %v\n", err)
		os.Exit(1)
	}

	// CLI应用不需要长期运行，执行完命令后退出
	if err := app.Stop(context.Background()); err != nil {
		fmt.Printf("Failed to stop CLI application: %v\n", err)
		os.Exit(1)
	}
}

// runCLI 运行CLI命令
func runCLI(logger *zap.Logger, client *gen.Client) error {
	// 创建根命令
	rootCmd := &cobra.Command{
		Use:   "services-cli",
		Short: "",
		Long:  "命令行工具",
	}

	// 添加迁移命令
	rootCmd.AddCommand(&cobra.Command{
		Use:   "migrate",
		Short: "执行数据库迁移",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Info("开始执行数据库迁移...")
			ctx := context.Background()
			if err := client.Schema.Create(
				ctx,
				migrate.WithDropIndex(true),
				migrate.WithDropColumn(true),
			); err != nil {
				logger.Error("数据库迁移失败", zap.Error(err))
				return err
			}
			logger.Info("数据库迁移完成")
			return nil
		},
	})

	// 执行命令
	if err := rootCmd.Execute(); err != nil {
		logger.Error("CLI command execution failed", zap.Error(err))
		return err
	}

	return nil
}
