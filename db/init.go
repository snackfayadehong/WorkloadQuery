package clientDb

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"io"
	"os"
)

var (
	DB         *gorm.DB
	privateKey *rsa.PrivateKey // 秘钥
	configs    *config
)

type config struct {
	DBClient struct {
		IP       string `json:"ip"`
		Username string `json:"username"`
		Password string `json:"password"`
		DbName   string `json:"db_name"`
	} `json:"DBClient"`
}

func InitDb() error {
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?&loc=%s&encrypt=disable", userName, password, ipAddr, port, dbName, loc)
	// dsn := "sqlserver://sa:密码@127.0.0.1:1433?database=dbStatus&encrypt=disable"
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s&encrypt=disable", configs.DBClient.Username, configs.DBClient.Password, configs.DBClient.IP, configs.DBClient.DbName)
	DB, _ = gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "TB_",
			SingularTable: true,
			NoLowerCase:   true,
		},
	})
	return nil
}
func writeEncryptionPwd() {
	file, err := os.OpenFile("../config.json", os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	v, _ := io.ReadAll(file)
	err = json.Unmarshal(v, &configs)
	if err != nil {
		panic(err)
	}
	// 密码加密
	encPwd, err := encryptionPwd(configs.DBClient.Password)
	if err != nil {
		panic(err)
	}
	configs.DBClient.Password = encPwd
	newConfig, err := json.Marshal(configs)
	_, err = file.WriteAt(newConfig, 0)
	if err != nil {
		panic(nil)
	}
}

func encryptionPwd(pwd string) (encryptionPwd string, err error) {
	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return
	}
	// 生成公钥
	publicKey := privateKey.PublicKey
	// 根据公钥加密
	encryptionBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &publicKey, []byte(pwd), nil)
	if err != nil {
		return
	}
	encryptionPwd = base64.StdEncoding.EncodeToString(encryptionBytes)
	return encryptionPwd, nil
}

// func Decryption() (pwd string) {
// 	decryptedBytes, err := privateKey.Decrypt(nil)
// }
