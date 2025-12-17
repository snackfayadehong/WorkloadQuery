package conf

import (
	"WorkloadQuery/encry"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"io"
	"os"
)

var Configs *Config

type Config struct {
	Server struct {
		IP              string `json:"Ip"`
		Port            string `json:"Port"`
		RunModel        string `json:"RunModel"`
		HealthCheckPort string `json:"HealthCheckPort"`
	} `json:"Server"`
	DBClient struct {
		IP       string `json:"ip"`
		Username string `json:"username"`
		Password string `json:"password"`
		DbName   string `json:"db_name"`
		IsEc     int    `json:"isEc"`
	} `json:"DBClient"`
	IPWhite struct {
		IPWhiteList []string `json:"IPWhiteList"`
	} `json:"IPWhite"`
	CustomTaskTime struct {
		Run       int `json:"Run"`
		StartTime int `json:"StartTime"`
		EndTime   int `json:"EndTime"`
	} `json:"CustomTaskTime"`
}

func InitSetting(rootPath string) error {
	file, err := os.OpenFile(rootPath+"\\config.json", os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	v, _ := io.ReadAll(file)
	err = json.Unmarshal(v, &Configs)
	if Configs.DBClient.IsEc == 0 {
		file.Close()
		if err = writeEncryptionPwd(rootPath); err != nil {
			return err
		}
	}
	return nil
}

// 读取配置文件密码加密后重新写入配置文件
func writeEncryptionPwd(rootPath string) error {
	// 生成公钥密钥文件
	encry.GenerateRSAKey(2048)
	// 打开public.pem公钥文件
	file, err := os.Open(rootPath + "\\encry\\public.pem")
	defer file.Close()
	// 读取公钥
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	// pem解码
	block, _ := pem.Decode(buf)
	// x509
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	// 类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	configFile, err := os.OpenFile(rootPath+"\\config.json", os.O_RDWR, 0666)
	defer configFile.Close()
	if err != nil {
		return err
	}
	v, _ := io.ReadAll(configFile)
	if err = json.Unmarshal(v, &Configs); err != nil {
		return err
	}
	// 密码加密
	// 对明文进行加密
	encPwd, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(Configs.DBClient.Password))
	if err != nil {
		return err
	}
	// 转base64
	base64Pwd := base64.StdEncoding.EncodeToString(encPwd)
	Configs.DBClient.Password = base64Pwd
	// 标记加密
	Configs.DBClient.IsEc = 1
	newConfig, err := json.Marshal(Configs)
	if _, err = configFile.WriteAt(newConfig, 0); err != nil {
		return err
	}
	return err
}

// DecryptionPwd 解密
func DecryptionPwd(rootPath string) (pwd string, err error) {
	// 打开私钥文件
	file, err := os.Open(rootPath + "\\encry\\private.pem")
	if err != nil {
		return
	}
	defer file.Close()
	// 读取私钥文件内容
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	// pem解码
	block, _ := pem.Decode(buf)
	// x509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return
	}
	pwdByte, err := base64.StdEncoding.DecodeString(Configs.DBClient.Password)
	if err != nil {
		return
	}
	decryptedBytes, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, pwdByte)
	if err != nil {
		return
	}
	return string(decryptedBytes), err
}
