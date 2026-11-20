package initialize

import (
	"fmt"
	"os"
	"path/filepath"
	"gin-spa-app-framework-back/common"
	"gin-spa-app-framework-back/configs"
	"gin-spa-app-framework-back/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DbInterface interface {
	connect() *gorm.DB
}

type SqlLiteDB struct{}

// func (m *MysqlDB) connect() *gorm.DB {
// 	//处理配置
// 	config := configs.SystemConfigs.Database
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
// 		config.User,
// 		config.Password,
// 		config.Host,
// 		config.Port,
// 		config.Database)

// 	//连接
// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
// 		NamingStrategy: schema.NamingStrategy{
// 			TablePrefix:   config.Prefix, // 表名前缀，`User` 的表名应该是 `t_users`
// 			SingularTable: true,          // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
// 		},
// 	})
// 	if err != nil {
// 		panic("Connect failed, check your database configuration. Error: " + err.Error())
// 	}

// 	db.AutoMigrate(&models.User{}, &models.SystemSettings{})
// 	return db
// }

func CheckDBConnection(dbPath string) error {
	fmt.Printf("=== 检查数据库连接 ===\n")
	fmt.Printf("数据库路径: %s\n", dbPath)

	// 1. 检查文件是否存在
	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		return fmt.Errorf("获取绝对路径失败: %w", err)
	}
	fmt.Printf("绝对路径: %s\n", absPath)

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Printf("⚠️  数据库文件不存在，GORM 将自动创建\n")
	} else {
		fmt.Printf("✅ 数据库文件存在\n")

		// 获取文件信息
		fileInfo, err := os.Stat(absPath)
		if err == nil {
			fmt.Printf("   文件大小: %d bytes\n", fileInfo.Size())
			fmt.Printf("   修改时间: %s\n", fileInfo.ModTime().Format("2006-01-02 15:04:05"))
		}
	}

	// 2. 检查目录权限
	dbDir := filepath.Dir(absPath)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		fmt.Printf("⚠️  目录不存在，尝试创建...\n")
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return fmt.Errorf("创建目录失败: %w", err)
		}
		fmt.Printf("✅ 目录创建成功\n")
	}

	// 3. 测试文件读写权限
	testFile := filepath.Join(dbDir, ".test_write")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return fmt.Errorf("目录无写入权限: %w", err)
	}
	os.Remove(testFile)
	fmt.Printf("✅ 目录有读写权限\n")

	return nil
}

func (s *SqlLiteDB) connect() *gorm.DB {

	//连接
	config := configs.SystemConfigs.Database
	// fmt.Println(config.DbPath)
	if err := CheckDBConnection(config.DbPath); err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open(config.DbPath), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   config.Prefix, // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true,          // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
	})
	if err != nil {
		panic("Connect failed, check your database configuration. Error: " + err.Error())
	}

	db.AutoMigrate(&models.User{}, &models.SystemSettings{})
	// 非常重要：获取底层 sql.DB 并设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic("获取底层数据库连接失败: " + err.Error())
	}

	// 设置连接池参数（SQLite 推荐设置）
	sqlDB.SetMaxOpenConns(1)    // SQLite 不支持并发写入，最大连接数设为1
	sqlDB.SetMaxIdleConns(1)    // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(0) // 连接最大生命周期，0表示不限制

	// 迁移
	if err := db.AutoMigrate(&models.User{}, &models.SystemSettings{}); err != nil {
		panic("迁移失败: " + err.Error())
	}

	return db
}

func connectIng() *gorm.DB {
	var db DbInterface = new(SqlLiteDB)
	return db.connect()

}

func initDB() {
	common.COM_DB = connectIng()
	// fmt.Println("common.COM_DB", common.COM_DB)
}
