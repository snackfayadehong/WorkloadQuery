package clientDb

import (
	"WorkloadQuery/conf"
	logger2 "WorkloadQuery/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	DB *gorm.DB
)

func Init() error {
	var DbPwd string
	DbPwd, err := conf.DecryptionPwd()
	if err != nil {
		return err
	}
	log := logger.New(newMyWriter(), logger.Config{LogLevel: logger.Info})
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?&loc=%s&encrypt=disable", userName, password, ipAddr, port, dbName, loc)
	// dsn := "sqlserver://sa:密码@127.0.0.1:1433?database=dbStatus&encrypt=disable"
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s&encrypt=disable", conf.Configs.DBClient.Username, DbPwd,
		conf.Configs.DBClient.IP, conf.Configs.DBClient.DbName)
	DB, _ = gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger:      log,
		PrepareStmt: true, // 执行任何 SQL 时都创建并缓存预编译语句，可以提高后续的调用速度
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "TB_",
			SingularTable: true,
			NoLowerCase:   true,
		},
	})
	// 根据配置文件设置选择程序环境
	switch conf.Configs.Server.RunModel {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "release":
		gin.SetMode(gin.ReleaseMode)
		gin.DisableConsoleColor() // 禁用控制台颜色，将日志写入文件时不需要控制台颜色
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
	return nil
}

type MyWriter struct {
	mlog *zap.SugaredLogger
}

func (m *MyWriter) Printf(format string, v ...interface{}) {
	if len(v) <= 3 { // 执行sql语句len(v)>3
		return
	}
	logstr := fmt.Sprintf("\r事件：SQL执行\r调用：%s\r时间：%.3fms\r行数：%v\rSQL：%s", v...)
	m.mlog.Infof("%s\r%s\r", logstr, logger2.LoggerEndStr)
}
func newMyWriter() *MyWriter {
	log := zap.L().Sugar()
	return &MyWriter{mlog: log}
}
