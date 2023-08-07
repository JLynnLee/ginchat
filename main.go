package main

import (
	"encoding/json"
	"fmt"
	"ginchat/router"
	"ginchat/utils"
)

func test(m map[string]int) {
	m["jjj"] = 3
}

func deepCopy(originalMap map[string]int) map[string]int {
	jsonData, _ := json.Marshal(originalMap)
	var newMap map[string]int
	err := json.Unmarshal(jsonData, &newMap)
	if err != nil {
		return nil
	}
	return newMap
}

func main() {
	utils.InitConfig()
	utils.InitMySql()
	utils.InitRedis()

	r := router.Router()
	r.Run(":8081")
}

// 冒泡排序
func bubblingSort(arr []int) []int {
	length := len(arr)
	if 0 == length {
		return arr
	}

	for num := 0; num < length; num++ {
		for k, _ := range arr {
			if k+1 >= length {
				break
			}
			if k+num >= length {
				break
			}
			if arr[k] > arr[k+1] {
				arr[k], arr[k+1] = arr[k+1], arr[k]
			}
		}
	}

	return arr
}

// 阶乘
func factorial(num uint64) uint64 {
	if num <= 1 {
		return num
	}
	diff := num - 1
	fmt.Println("----", num, "=====", diff)
	return num * factorial(diff)
}

// 选择排序
func changeSort(arr []int) []int {
	length := len(arr)
	for i := 0; i < length; i++ {
		minIdx := i
		for j := i + 1; j < length; j++ {
			if arr[j] < arr[minIdx] {
				minIdx = j
			}
		}
		arr[minIdx], arr[i] = arr[i], arr[minIdx]
	}

	return arr
}

// 插入排序
func insertSort(arr []int) []int {
	for i := range arr {
		preIdx := i - 1
		current := arr[i]
		for preIdx >= 0 && current < arr[preIdx] {
			arr[preIdx+1] = arr[preIdx]
			preIdx--
			fmt.Println("===========", preIdx)
		}
		arr[preIdx+1] = current
	}

	return arr
}
