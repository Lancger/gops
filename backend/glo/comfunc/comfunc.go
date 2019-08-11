package comfunc

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"gops/backend/glo/encrypt"
	"os"
	"reflect"
	"strings"
	"time"
)

// DefaultPassword const Params
const DefaultPassword string = "123456"

// TokenInvaild Const Params
const TokenInvaild string = "Token invaild"

// FormatTs 格式化时间戳为字符串
func FormatTs(ts int64) (res string) {
	res = time.Unix(ts, 0).Format("2006-01-02 15:04:05")
	return
}

// FormatShortTs 格式化时间戳为字符串
func FormatShortTs(ts int64) (res string) {
	res = time.Unix(ts, 0).Format("20060102150405")
	return
}

// FlorPageInt 分页整除判断
func FlorPageInt(pageSize int, total int) (res int) {
	res = 0
	if pageSize <= 0 {
		return
	}
	mix := total % pageSize
	if mix != 0 {
		res = total/pageSize + 1
	} else {
		res = total / pageSize
	}
	return
}

// GetDefaultPage 校验页数，小于0返回默认1
func GetDefaultPage(page int) (res int) {
	if page <= 0 {
		res = 1
		return
	}
	res = page
	return
}

// GetDefaultPageSize 校验每页行数，小于0返回默认20
func GetDefaultPageSize(pageSize int) (res int) {
	if pageSize <= 0 {
		res = 20
		return
	}
	res = pageSize
	return
}

// EncryptToken Init Token
func EncryptToken(srcStr string, timeTs int64, key string) (token string) {
	encStr := fmt.Sprintf("%s_%d", srcStr, timeTs)
	token, _ = encrypt.AesEncryptString([]byte(encStr), []byte(key))
	return
}

// DecryptToken Token
func DecryptToken(token string, key string) (srcStr string) {
	encStr, _ := encrypt.AesDecryptString(token, []byte(key))
	srcStr = strings.Split(string(encStr), "_")[0]
	return
}

// GetDayByTimeStampRange 获取时间戳之间的日期
func GetDayByTimeStampRange(startTs, endTs int64) []string {
	var (
		retDay []string
		tmpTs  int64
		i      int64
	)
	if startTs > endTs {
		tmpTs = startTs
		startTs = endTs
		endTs = tmpTs
	}
	i = 0
	for {
		tmpTs = i*86400 + startTs
		retDay = append(retDay, time.Unix(tmpTs, 0).Format("2006-01-02"))
		if tmpTs > endTs {
			break
		}
		i++
	}
	return retDay
}

// GetTodayFirstTs 获取当天凌晨时间戳
func GetTodayFirstTs() int64 {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	timeNumber := t.Unix()
	return timeNumber
}

// FormatSubTimeStr 格式化时间间隔
func FormatSubTimeStr(duration int64) (res string) {
	var (
		day    int64
		hour   int64
		minute int64
	)
	day = duration / 86400
	hour = (duration - day*86400) / 3600
	minute = (duration - day*86400 - hour*3600) / 60
	return fmt.Sprintf("%d天%d小时%d分", day, hour, minute)
}

// StrArrayIndexOf 检查字符串数组中是否存在
func StrArrayIndexOf(arr []string, val string) (index int, status bool) {
	index = 0
	status = false
	for i, r := range arr {
		if r == val {
			index = i
			status = true
			break
		}
	}
	return
}

// Md5String 加密
func Md5String(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// PathExists 判断目录是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// UniqArr 字符串数组去重
func UniqArr(arr []string) (ret []string) {
	mapDict := map[string]bool{}
	for _, a := range arr {
		mapDict[a] = true
	}
	for k := range mapDict {
		if k == `` {
			continue
		} else {
			ret = append(ret, k)
		}
	}
	return
}

/*
@func           将map转为Struct
@param
    mmap        需要转换的map[string]interface
    structure   转换后的结构体指针
@return
    error       错误信息

    暂不支持递归转换
*/
func MapToStruct(mmap map[string]interface{}, structure interface{}) (err error) {
	defer func() {
		if errs := recover(); errs != nil {
			err = errors.New("调用出错")
		}
	}()
	ptp := reflect.TypeOf(structure)
	pv := reflect.ValueOf(structure)
	switch ptp.Kind() {
	case reflect.Ptr:
		if ptp.Elem().Kind() == reflect.Struct {
			fmt.Println("sss")
			break
		} else {
			return errors.New("需要*struct类型，却传入*" + ptp.Elem().Kind().String() + "类型")
		}
	default:
		return errors.New("需要*struct类型，却传入" + ptp.Kind().String() + "类型")
	}
	tp := ptp.Elem()
	v := pv.Elem()
	num := tp.NumField()
	for i := 0; i < num; i++ {
		name := tp.Field(i).Name
		tag := tp.Field(i).Tag.Get("map")
		if len(tag) != 0 {
			name = tag
		}
		value, ok := mmap[name]
		if !ok {
			continue
		}
		//能够设置值，且类型相同
		if v.Field(i).CanSet() {
			if v.Field(i).Type() == reflect.TypeOf(value) {
				v.Field(i).Set(reflect.ValueOf(value))
			} else {
				continue
			}
		} else {
			continue
		}
	}
	return nil
}
