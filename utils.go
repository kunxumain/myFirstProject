package common

import (
	"math"
	"strconv"
	"strings"
)

//用于uuid加密

func StringToArray(input string) []int {
	output := []int{}
	for _, v := range input {
		output = append(output, int(v))
	}
	for i, j := 0, len(output)-1; i < j; i, j = i+1, j-1 {
		output[i], output[j] = output[j], output[i]
	}
	return output
}

func GetInput(input string) <-chan int {
	out := make(chan int)
	go func() {
		for _, b := range StringToArray(input) {
			out <- b
		}
		close(out)
	}()
	return out
}

func SQ(in <-chan int) <-chan int {
	out := make(chan int)
	var base, i float64 = 2, 0
	go func() {
		for n := range in {
			out <- (n - 48) * int(math.Pow(base, i)) //48是“0”对应的ASCII码，依次递增
			i++
		}
		close(out)
	}()
	return out
}

// 转换成int
func ToInt(input string) int {
	c := GetInput(input)
	out := SQ(c)
	sum := 0
	for o := range out {
		sum += o
	}
	return sum
}

// int转二进制
func ConverToBinary(n int) string {
	res := ""
	for ; n > 0; n /= 2 {
		lsb := n % 2
		res = strconv.Itoa(lsb) + res
	}
	return res
}

// 格式化页面传入的IDs
func SplitToInt32List(str string, sep string) (int32List []int32) {
	tempStr := strings.Split(str, sep)
	if len(tempStr) > 0 {
		for _, item := range tempStr {
			if item == "" {
				continue
			}
			val, err := strconv.ParseInt(item, 10, 32)
			if err != nil {
				continue
			}
			int32List = append(int32List, int32(val))
		}
	}
	return int32List
}
