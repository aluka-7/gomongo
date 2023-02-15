package gomongo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aluka-7/configuration"
	"github.com/aluka-7/utils"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"strings"
	"time"
)

type Config struct {
	Uri             string `json:"uri"`             // 连接uri
	Database        string `json:"database"`        // 数据库
	MaxPoolSize     uint64 `json:"maxPoolSize"`     // 连接池最大连接数
	MinPoolSize     uint64 `json:"minPoolSize"`     // 连接池最小连接数
	MaxConnecting   uint64 `json:"maxConnecting"`   // 连接池可以同时建立的最大连接数
	MaxConnIdleTime uint64 `json:"maxConnIdleTime"` // 连接最大空闲时间
	TimeOut         int64  `json:"timeOut"`
}

type goMongo struct {
	systemId   string
	cfg        configuration.Configuration
	privileges map[string][]string
}

type GoMongo interface {
	Config(dsID string) *Config
	Connection(dsID string) *mongo.Client
}

/**
 * 获取mongo引擎的唯一实例。
 * @return
 */
func Engine(cfg configuration.Configuration, systemId string) GoMongo {
	fmt.Println("Loading Mongo Engine")
	return &goMongo{cfg: cfg, systemId: systemId, privileges: make(map[string][]string, 0)}
}

func (d *goMongo) Config(dsID string) *Config {
	ds, dsID, err := d.getConfiguration(dsID, d.systemId)
	if len(ds.Uri) == 0 || err != nil {
		panic(fmt.Sprintf("数据源[%s]配置未指定或者读取时发生错误:%+v", dsID, err))
	}
	return ds
}
func (d *goMongo) Connection(dsID string) *mongo.Client {
	c := d.Config(dsID)
	opt := options.Client().ApplyURI(c.Uri)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.TimeOut))
	defer cancel()
	opt.SetMaxPoolSize(c.MaxPoolSize)
	opt.SetMinPoolSize(c.MinPoolSize)
	opt.SetMaxConnecting(c.MaxConnecting)
	opt.SetMaxConnIdleTime(time.Duration(c.MaxConnIdleTime))
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		panic(fmt.Sprintf("初始化mongo引擎出错%+v", err))
	}
	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		panic(fmt.Sprintf("mongo连接ping错误%+v", err))
	}
	return client.Database(c.Database).Client()
}

/**
 * 获取指定标示的数据源的配置信息，返回的配置Config对象
 * 需要特别说明：如果给定的数据源标示为Null，则表明是要获取当前业务系统的默认数据源配置信息。
 * </p>
 *
 * @param dsID 数据源的标示，如果为Null则表明是默认数据源
 * @return
 */
func (d *goMongo) getConfiguration(dsID, csID string) (*Config, string, error) {
	config := &Config{}
	// 如果是获取默认的数据源，则使用当前系统的标示，否则鉴权
	if len(dsID) == 0 || dsID == csID {
		dsID = csID
	} else {
		plist := d.systemPrivileges(csID) // 数据库的访问权限鉴权
		if len(plist) == 0 || utils.ContainsString(plist, dsID) == -1 {
			return config, "", fmt.Errorf("系统[%s]无数据源[%s]的访问权限", csID, dsID)
		}
	}
	err := d.readFromConfiguration(dsID, config)
	return config, dsID, err
}

/*
*
加载数据库的访问权限鉴权
*/
func (d *goMongo) systemPrivileges(csID string) []string {
	d.cfg.Get("base", "mongo", "", []string{"privileges"}, d)
	plist := d.privileges[csID]
	fmt.Printf("系统[%s]的数据源权限:%s", csID, strings.Join(plist, ","))
	return plist
}
func (d *goMongo) Changed(data map[string]string) {
	for _, v := range data {
		var vl map[string][]string
		if err := json.Unmarshal([]byte(v), &vl); err == nil {
			for k, _v := range vl {
				d.privileges[k] = _v
			}
		}
	}
}
func (d *goMongo) readFromConfiguration(dsID string, config *Config) error {
	ex := d.readCommonProperties(config)
	if ex != nil {
		return ex
	}
	fmt.Printf("从配置中心读取数据源配置:/base/mongo/%s\n", dsID)
	ex = d.cfg.Clazz("base", "mongo", "", dsID, config)
	if ex != nil {
		log.Error().Err(ex).Msgf("数据源[%s]的配置获取失败", dsID)
	}
	return ex
}

func (d *goMongo) readCommonProperties(config *Config) error {
	fmt.Println("从配置中心的读取通用数据源配置:/base/mongo/common")
	vl, err := d.cfg.String("base", "mongo", "", "common")
	if err != nil {
		log.Error().Err(err).Msg("配置中心的通用数据源配置获取失败:%v")
	} else {
		if err = json.Unmarshal([]byte(vl), config); err != nil {
			log.Error().Err(err).Msg("解析数据源的通用配置失败")
		}
	}
	return err
}
