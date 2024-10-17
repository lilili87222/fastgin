package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
)

// 结构体转为json
func Struct2Json(obj interface{}) string {
	str, err := json.Marshal(obj)
	if err != nil {
		panic(fmt.Sprintf("[Struct2Json]转换异常: %v", err))
	}
	return string(str)
}

// json转为结构体
func Json2Struct(str string, obj interface{}) {
	// 将json转为结构体
	err := json.Unmarshal([]byte(str), obj)
	if err != nil {
		panic(fmt.Sprintf("[Json2Struct]转换异常: %v", err))
	}
}

// json interface转为结构体
func JsonI2Struct(str interface{}, obj interface{}) {
	JsonStr := str.(string)
	Json2Struct(JsonStr, obj)
}

// StructsToMap 函数，直接调用 StructToMap
func StructsToMap(objs []any, include bool, fields ...string) ([]map[string]any, error) {
	var results []map[string]any

	for _, obj := range objs {
		result, err := StructToMap(obj, include, fields...)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}
func StructToMap(obj any, include bool, fields ...string) (map[string]any, error) {
	result := make(map[string]any)
	v := reflect.ValueOf(obj)

	// 处理结构体指针
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, fmt.Errorf("nil pointer provided")
		}
		v = v.Elem() // 获取指针指向的值
	}

	// 确保传入的是结构体类型
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %s", v.Kind())
	}

	// 遍历所有字段，包括嵌套结构体字段
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		fieldName := field.Name
		fieldValue := v.Field(i)

		// 处理嵌入结构体的字段
		if field.Anonymous {
			// 如果是嵌入字段，递归调用 StructToMap
			embeddedMap, err := StructToMap(fieldValue.Interface(), include, fields...)
			if err != nil {
				return nil, err
			}
			for k, v := range embeddedMap {
				result[k] = v
			}
			continue
		}
		contains := slices.Contains(fields, fieldName)

		// 检查字段是否在包含或排除列表中
		if include && contains {
			result[fieldName] = fieldValue.Interface()
		} else if !include && !contains {
			result[fieldName] = fieldValue.Interface()
		}
	}

	return result, nil
}
