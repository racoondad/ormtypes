/*
 * @Author       : lptecodad lptecodad@sina.com
 * @Date         : 2023-01-05 10:20:06
 * @LastEditors  : lptecodad lptecodad@sina.com
 * Copyright (c) 2023 by lptecodad lptecodad@sina.com, All Rights Reserved.
 */
package ormtypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/duke-git/lancet/v2/convertor"
)

type Array[T ~string | ~int64 | ~int32 | ~int | ~int8 | ~uint64 | ~uint32 | ~uint | ~uint8] []T

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (arr *Array[T]) Scan(value interface{}) error {
	var bytes []byte
	if value == nil {
		bytes = nil
	} else {
		bytes = []byte(fmt.Sprintf("%v", value))
	}
	if len(bytes) > 0 {
		return json.Unmarshal(bytes, arr)
	}
	*arr = make([]T, 0)
	return nil
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (arr Array[T]) Value() (driver.Value, error) {
	if arr == nil {
		return "[]", nil
	}
	return convertor.ToString(arr), nil
}

func (arr Array[T]) Len() int {
	return len(arr)
}

func (arr Array[T]) Less(i, j int) bool {
	return arr[i] < arr[j]
}

func (arr Array[T]) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func (arr Array[T]) Contains(in T) bool {
	for _, v := range arr {
		if v == in {
			return true
		}
	}
	return false
}
