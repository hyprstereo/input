package utils

import (
	"strings"

	godao "github.com/hyprstereo/go-dao"

	cmap "github.com/orcaman/concurrent-map"
	"github.com/tidwall/gjson"
)

type Map = godao.Map
type DMap struct {
	cmap.ConcurrentMap
}

func (d *DMap) String(key string) (res string) {
	if v, ok := d.Get(key); ok {
		res = v.(string)
	}
	return
}

func (d *DMap) JsonGet(key string, p string) (res gjson.Result) {
	if v, ok := d.Get(key); ok {
		val := v.(string)
		res = gjson.Get(val, p)
	}
	return
}

func (d *DMap) Map(key string) (res Map) {
	if v, ok := d.Get(key); ok {
		res = Map(v.(map[string]any))
	}
	return
}

func (d *DMap) KeysByPrefix(prefix string) (res []string) {
	res = make([]string, 0)
	keys := d.Keys()
	for _, k := range keys {
		if strings.HasPrefix(k, prefix) {
			res = append(res, k)
		}
	}
	return
}

func (d *DMap) GetByPrefix(prefix string) (res []any) {
	res = make([]any, 0)
	keys := d.Keys()
	for _, k := range keys {
		if strings.HasPrefix(k, prefix) {
			v, _ := d.Get(k)
			res = append(res, v)
		}
	}
	return
}

func (d *DMap) Match(key string, limit ...int) (res []any) {
	res = make([]any, 0)
	keys := d.Keys()
	for _, k := range keys {
		if Match(k, key) {
			v, _ := d.Get(k)
			res = append(res, v)
			if len(limit) > 0 && len(res) == limit[0] {
				break
			}
		}
	}
	return
}

func NewDynamicMap() *DMap {
	return &DMap{
		ConcurrentMap: cmap.New(),
	}
}
