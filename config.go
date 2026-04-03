package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	defKeyFile  = ""
	defKeyData  = ""
	defLogLevel = "info"
	defLogFile  = "/var/log/vmStarter.log"
	defVMs      = ""
	defDelaySec = 2.5
)

type Config struct {
	keyFile  string
	keyData  string
	logLevel string
	logFile  string
	vmList   []string
	delaySec float32
}
type YCToken struct {
	IAM string
}
type AuthorizedKey struct {
	KID              string `json:"id"`
	ServiceAccountID string `json:"service_account_id"`
	CreatedAt        string `json:"created_at"`
	KeyAlgorithm     string `json:"key_algorithm"`
	PublicKey        string `json:"public_key"`
	PrivateKey       string `json:"private_key"`
}

func getEnv(in string, def string) string {
	if value, ok := os.LookupEnv(in); ok {
		return strings.TrimSpace(value)
	}
	return def
}

func (c *Config) New() {
	c.keyFile = getEnv("YC_KEY_FILE", defKeyFile)
	c.keyData = getEnv("YC_KEY_DATA", defKeyData)
	c.logLevel = getEnv("YC_LOG_LEVEL", defLogLevel)
	c.logFile = getEnv("YC_LOG_FILE", defLogFile)
	c.vmList = commStrToStrArray(getEnv("YC_VMS", defVMs))
	c.delaySec = defDelaySec
	if c.keyFile == "" && c.keyData == "" {
		log.Fatal("YC key file or key data is empty, exiting...")
	}
}

func (c *Config) Print() string {
	return fmt.Sprintf("keyFile: %s; keyData: %s...; logLevel: %s; vmList: %s", c.keyFile, c.keyData[0:25], c.logLevel, c.vmList)
}

func (c *Config) setupLogs() {
	file, err := os.OpenFile(c.logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to log to file, using default stderr")
	}

	mw := io.MultiWriter(os.Stdout, file)

	log.SetOutput(mw)

	// Optional: Set formatting
	//log.SetFormatter(&log.JSONFormatter{})

	ll, err := log.ParseLevel(c.logLevel)
	if err != nil {
		ll = log.InfoLevel
	}
	log.SetLevel(ll)
}

func ReadData() []byte {
	if conf.keyFile != "" {
		data, err := os.ReadFile(conf.keyFile)
		if err != nil {
			panic(err)
		}
		return data
	} else {
		data, err := base64.StdEncoding.DecodeString(conf.keyData)
		if err != nil {
			panic(err)
		}
		return data
	}
	return nil
}

func ParseAuthData() AuthorizedKey {
	data := ReadData()
	var authData AuthorizedKey
	err := json.Unmarshal(data, &authData)
	if err != nil {
		panic(err)
	}
	log.Debugf("id: %s...", authData.KID[:10])
	log.Debugf("service_account_id: %s", authData.ServiceAccountID)
	return authData
}

func commStrToStrArray(str string) []string {
	if str == "" {
		return []string{}
	}
	return strings.Split(str, ",")
}
