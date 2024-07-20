package clientDb

import (
	"WorkloadQuery/conf"
	logger2 "WorkloadQuery/logger"
	"fmt"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 新思路实现以下接口 通过.m.mlog.Infof(format, logstr, logger2.LoggerEndStr)打印日志 详见logger.gormLogger
//	LogMode(LogLevel) Interface
//	Info(context.Context, string, ...interface{})
//	Warn(context.Context, string, ...interface{})
//	Error(context.Context, string, ...interface{})
//	Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)

var (
	DB *gorm.DB
)

func Init() error {
	var DbPwd string
	DbPwd, err := conf.DecryptionPwd()
	if err != nil {
		return err
	}
	//log := logger.New(newMyWriter(), logger.Config{LogLevel: logger.Info})
	log := logger2.NewGormCustomLogger()
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
	//// 根据配置文件设置选择程序环境
	//switch conf.Configs.Server.RunModel {
	//case "debug":
	//	gin.SetMode(gin.DebugMode)
	//case "release":
	//	gin.SetMode(gin.ReleaseMode)
	//	gin.DisableConsoleColor() // 禁用控制台颜色，将日志写入文件时不需要控制台颜色
	//case "test":
	//	gin.SetMode(gin.TestMode)
	//default:
	//	gin.SetMode(gin.DebugMode)
	//}
	return nil
}

//type MyWriter struct {
//	logger.Interface
//	mlog *zap.SugaredLogger
//}
//
//func (m *MyWriter) Printf(format string, v ...interface{}) {
//	//if len(v) <= 3 { // 执行sql语句len(v)>3
//	//	return
//	//}
//	fmt.Println(format)
//	//for i, v1 := range v {
//	//	fmt.Printf("i:%d,v1:%v\n", i, v1)
//	//}
//	logstr := fmt.Sprintf("\r\n事件：SQL执行\r\n调用：%s\r\n时间：%.3fms\r\n行数：%v\r\nSQL：%s", v...)
//	m.mlog.Infof("%s\r\n%s\r\n", logstr, logger2.LoggerEndStr)
//}
//func newMyWriter() *MyWriter {
//	log := zap.L().Sugar()
//	return &MyWriter{mlog: log}
//}
