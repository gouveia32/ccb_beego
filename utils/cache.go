// Copyright 2018 ccb_beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"time"

	"github.com/beego/beego/v2/client/cache"
	"github.com/beego/beego/v2/core/logs"

	beego "github.com/beego/beego/v2/server/web"
)

var cc cache.Cache

func InitCache() {
	host, _ := beego.AppConfig.String("cache::redis_host")
	//passWord := beego.AppConfig.String("cache::redis_password")
	var err error
	defer func() {
		if r := recover(); r != nil {
			cc = nil
		}
	}()
	//cc, err = cache.NewCache("redis", `{"conn":"`+host+`","password":"`+passWord+`"}`)
	cc, err = cache.NewCache("redis", `{"conn":"`+host+`"}`)
	if err != nil {
		logs.Error("Connect to the redis host " + host + " failed")
		logs.Error(err)
	}
}

// SetCache
func SetCache(key string, value interface{}, timeout int) error {
	data, err := Encode(value)
	if err != nil {
		return err
	}
	if cc == nil {
		return errors.New("cc is nil")
	}

	defer func() {
		if r := recover(); r != nil {
			logs.Error(r)
			cc = nil
		}
	}()
	timeouts := time.Duration(timeout) * time.Second
	err = cc.Put(context.Background(), key, data, timeouts)
	if err != nil {
		logs.Error(err)
		logs.Error("SetCache falhou com a chave: " + key)
		return err
	} else {
		return nil
	}
}

func GetCache(key string, to interface{}) error {
	if cc == nil {
		return errors.New("cc is nil")
	}

	defer func() {
		if r := recover(); r != nil {
			logs.Error(r)
			cc = nil
		}
	}()

	data, _ := cc.Get(context.Background(), key)
	if data == nil {
		return errors.New("Cache não existe")
	}

	err := Decode(data.([]byte), to)
	if err != nil {
		logs.Error(err)
		logs.Error("GetCache falhou com a chave: " + key)
	}

	return err
}

// DelCache
func DelCache(key string) error {
	if cc == nil {
		return errors.New("cc is nil")
	}
	defer func() {
		if r := recover(); r != nil {
			//fmt.Println("get cache error caught: %v\n", r)
			cc = nil
		}
	}()
	err := cc.Delete(context.Background(), key)
	if err != nil {
		return errors.New("Cache Falha ao excluir")
	} else {
		return nil
	}
}

// Encode
// 用gob进行数据编码
func Encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode
// 用gob进行数据解码
func Decode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}
