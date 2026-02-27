package database

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"SuperBizAgent/internal/model"
)

var DB *gorm.DB

// Init 初始化数据库连接
func Init() error {
	ctx := gctx.New()

	// 从配置读取数据库连接信息
	host, _ := g.Cfg().Get(ctx, "database.host", "localhost")
	port, _ := g.Cfg().Get(ctx, "database.port", "3306")
	user, _ := g.Cfg().Get(ctx, "database.user", "root")
	pass, _ := g.Cfg().Get(ctx, "database.pass", "")
	name, _ := g.Cfg().Get(ctx, "database.name", "fo_sentinel")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local",
		user.String(), pass.String(), host.String(), port.String(), name.String())

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// 自动迁移数据库表
	if err := AutoMigrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate() error {
	// 修复表字符集
	DB.Exec("ALTER TABLE reports CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
	DB.Exec("ALTER TABLE report_templates CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")

	return DB.AutoMigrate(
		&model.Subscription{},
		&model.SecurityEvent{},
		&model.Report{},
		&model.ReportTemplate{},
		&model.FetchLog{},
	)
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
