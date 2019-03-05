package http

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/securecookie"
	"github.com/hongkailiu/test-go/pkg/lib/util"
	log "github.com/sirupsen/logrus"
)

type config struct {
	sessionKey     []byte
	ghClientID     string
	ghClientSecret string
	ggClientID     string
	ggClientSecret string
	ggRedirectURL  string
}

func loadConfig() *config {
	config := &config{}
	config.sessionKey = securecookie.GenerateRandomKey(32)
	newSessionKey, err := getSessionKey()
	if err != nil {
		log.Warnf("error found when getSessionKey(): %s", err.Error())
	} else {
		config.sessionKey = newSessionKey
	}

	config.ghClientID = util.Getenv("gh_client_id", "gh_client_id")
	config.ghClientSecret = util.Getenv("gh_client_secret", "gh_client_secret")
	config.ggClientID = util.Getenv("gg_client_id", "gg_client_id")
	config.ggClientSecret = util.Getenv("gg_client_secret", "gg_client_secret")
	config.ggRedirectURL = util.Getenv("gg_redirect_url", "gg_redirect_url")

	if config.ggRedirectURL == "" {
		config.ggRedirectURL = "http://127.0.0.1:8080/google_oauth_cb"
	}
	return config
}

func getSessionKey() ([]byte, error) {
	key := util.Getenv("session_key", "session_key")
	trimKey := strings.TrimSuffix(strings.TrimPrefix(key, "["), "]")
	bytes := strings.Split(trimKey, " ")
	if len(bytes) != 32 {
		return nil, fmt.Errorf("key length is %d (only 32 allowed)", len(bytes))
	}

	var result []byte
	for _, b := range bytes {
		i, err := strconv.ParseUint(b, 10, 8)
		if err != nil {
			return nil, fmt.Errorf("get error when strconv.ParseUint(b, 10, 8) for b=%s", b)
		}
		result = append(result, byte(i))
	}
	log.WithFields(log.Fields{"result": result}).Warnf("got secret from env. var. session_key")
	return result, nil
}

type dbConfig struct {
	// "host=myhost port=myport user=gorm dbname=gorm password=mypassword"
	host     string
	port     string
	user     string
	dbname   string
	password string
}

func loadDBConfig() *dbConfig {
	config := &dbConfig{}
	config.host = os.Getenv("db_host")
	if config.host == "" {
		config.host = "localhost"
	}
	config.port = os.Getenv("db_port")
	if config.port == "" {
		config.port = "5432"
	}
	config.user = os.Getenv("db_user")
	if config.user == "" {
		config.user = "redhat"
	}
	config.dbname = os.Getenv("db_name")
	if config.dbname == "" {
		config.dbname = "ttt"
	}
	config.password = os.Getenv("db_password")
	if config.password == "" {
		config.password = "redhat"
	}
	return config
}

func (c dbConfig) getDBString() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", c.host, c.port, c.user, c.dbname, c.password)
}
