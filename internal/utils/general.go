package utils

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/hyprstereo/go-dao/encoding/json"

	"github.com/savsgio/gotils/uuid"
	"github.com/spf13/afero"
)

func UUID() string {
	return uuid.V4()
}

var fs = afero.NewOsFs()

func GetEnv(key string, defaultValue ...string) (value string) {
	val := os.Getenv(key)
	if val == "" && len(defaultValue) > 0 {
		os.Setenv(key, defaultValue[0])
		value = defaultValue[0]
	} else {
		value = val
	}
	return
}

func SetEnv(key, value string) error {
	return os.Setenv(key, value)
}

func Fetch(addr string, query map[string]interface{}, vs ...interface{}) (r interface{}) {

	//addr := API + resourceId
	//addr = fmt.Sprintf("%s?%s", addr, "populate=*")

	uri, _ := url.Parse(addr)
	cli := fiber.AcquireClient()
	a := cli.Get(uri.String())

	if t, ok := query["token"]; ok {
		a.Request().Header.Add("Authorization", "Bearer "+t.(string))
		//mt.Println(uri.String(), t)
	}

	_, res, _ := a.Get([]byte{}, uri.String())
	if len(vs) > 0 {
		json.Decode(res, vs[0])
	}

	r = res

	fiber.ReleaseAgent(a)
	fiber.ReleaseClient(cli)
	return
}

func FetchURL(uri string, query map[string]interface{}) (r []byte, err error) {

	addr := uri
	if len(query) > 0 {
		addr += `?`
		for n, v := range query {
			addr += n + `="` + fmt.Sprint(v) + `" `
		}
	}
	//addr = fmt.Sprintf("%s?%s", addr, "populate=*")

	//uri, _ := url.Parse(addr)
	cli := fiber.AcquireClient()
	a := cli.Get(addr)

	// if t, ok := query["token"]; ok {
	// 	a.Request().Header.Add("Authorization", "Bearer "+t.(string))
	// 	//mt.Println(uri.String(), t)
	// }

	_, r, err = a.Get([]byte{}, addr)

	fiber.ReleaseAgent(a)
	fiber.ReleaseClient(cli)
	return
}

func GetPath(p string) (v interface{}) {
	return
}

func GetArray(data, p string, keys []string) (v []fiber.Map) {
	// res := gjson.Get(data, p).Array()
	// for _, v := range res {
	// 	val := v.Map()
	// 	v = append(v, fiber.Map{})
	// }
	return
}

func HaveAny(src string, list []string) (b bool) {
	for _, l := range list {
		if strings.TrimSpace(src) == strings.TrimSpace(l) {
			b = true
			return
		}
	}
	return
}

func SearchGlob(src string) ([]string, error) {
	return afero.Glob(fs, src)
}

func SearchGlobMap(root, src string) (m map[string]interface{}) {
	m = make(map[string]interface{})
	res, _ := SearchGlob(src)
	for _, v := range res {
		if isDir, _ := afero.IsDir(fs, v); isDir {
			m[v] = SearchGlobMap(v, src)
		}
	}
	return
}

func Invoke(v interface{}, args ...interface{}) (res interface{}) {
	ref := reflect.ValueOf(v)

	out := ref.Call(ItoValue(args))
	res = ValueToI(out)
	return
}

func ValueToI(vals []reflect.Value) []interface{} {
	r := make([]interface{}, len(vals))
	for x, v := range vals {
		r[x] = reflect.ValueOf(v)
	}
	return r
}

func ItoValue(args []interface{}) []reflect.Value {
	vals := make([]reflect.Value, len(args))
	for x, v := range args {
		vals[x] = reflect.ValueOf(v)
	}
	return vals
}

func StripQuotes(v string) string {
	v = strings.TrimPrefix(strings.TrimSpace(v), "'")
	v = strings.TrimSuffix(v, "'")
	v = strings.TrimPrefix(v, "`")
	v = strings.TrimSuffix(v, "`")
	v = strings.TrimPrefix(v, `"`)
	v = strings.TrimSuffix(v, `"`)
	return v
}
