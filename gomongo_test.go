package gomongo

import (
	"github.com/aluka-7/configuration"
	"github.com/aluka-7/configuration/backends"
	"testing"
)

func Test_Connect(t *testing.T) {
	conf := configuration.MockEngine(t, backends.StoreConfig{Exp: map[string]string{
		"/system/base/mongo/1000":   "{\"uri\":\"mongodb://localhost:27017/\",\"timeOut\":2}",
		"/system/base/mongo/common": "{\"uri\":\"mongodb://localhost:27017/\",\"timeOut\":2}",
	}})

	Engine(conf, "1000").Connection("")
}
