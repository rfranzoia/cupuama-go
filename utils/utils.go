package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"
)

// YMFormat Date format for Year/Month equivalent to YYYYMM
const YMFormat = "200601"

// YMDFormat Date format for Year/Month/Day equivalent to YYYYMM
const YMDFormat = "20060102"

// DateTimeLong Date format for Year/Month/Day Hour:Minute equivalent to YYYY-MM-DD HH:MM
const DateTimeLong = "2006-01-02 15:04"

// GetMD5Hash Returns the MD5 hash code for a string
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// MessageJSON Default return values for services
type MessageJSON struct {
	Message string      `json:"message,ommitempty"`
	Value   interface{} `json:"value,ommitempty"`
}

// AuthLog log de autorizacao por token
// func AuthLog(c echo.Context, funcName string) {
// 	user := c.Get("user")
// 	token := user.(*jwt.Token)
// 	claims := token.Claims.(jwt.MapClaims)

// 	log.Println("Func:", funcName, "Id:", claims["jti"], "Issuer:", claims["iss"])
// }

// JSONTime Formato de data/hora para JSON
type JSONTime struct {
	time.Time
}

// MarshalJSON Implementacao de Marshal para JSONTime
func (t *JSONTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", t.Time.Format("2006-01-02 15:04"))
	return []byte(stamp), nil
}

// UnmarshalJSON Implementacao de Unmarshal para JSONTime
func (t *JSONTime) UnmarshalJSON(buf []byte) error {
	tt, err := time.Parse("2006-01-02 15:04", strings.Trim(string(buf), `"`))
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}

// Substring returns the substring of a string value
func Substring(value string, start, end uint) string {
	runes := []rune(value)
	result := string(runes[start:end])
	return result
}

// CreateSQLCache creates a map of queries to be used by the repository
func CreateSQLCache(queriesLocation ...string) (map[string]string, error) {

	if len(queriesLocation) == 0 {
		queriesLocation = append(queriesLocation, "./queries/*.sql")
	}

	myCache := map[string]string{}
	var queries []string
	var err error

	for _, queryPath := range queriesLocation {
		queries, err = filepath.Glob(queryPath)
		if err != nil {
			log.Fatal(fmt.Sprintf("cannot read queries from path %s", queryPath), err)
		}
	}

	if len(queries) == 0 {
		log.Fatal("no queries were found")
	}

	for _, query := range queries {
		name := filepath.Base(query)
		sql, err := ioutil.ReadFile(query)
		if err != nil {
			log.Println(err)
			return myCache, err
		}
		myCache[name] = string(sql)
	}

	return myCache, nil
}
