package task

import (
	"fmt"
	"strconv"
)

func TestTask1() {
	//单数
	var nums = []int{1, 2, 1, 2, 5, 5, 4, 3, 3}
	SingleNumber(nums)
	//回文数字
	var num = 14322341
	IsPalindrome(num)
	//有效括号
	s := "{}[]()[{[]}]{}[]"
	IsValid(s)
	//最长公共前缀
	strs := []string{"efrlower", "flow", "flight"}
	LongestCommonPrefix(strs)
	//加1
	ss := []int{9, 9, 9, 9, 9, 8}
	AddOne(ss)
	//去重
	duplicates := []int{0, 1, 1, 2, 2, 3, 3, 3, 4, 5, 6, 6, 7, 7, 8}
	RemoveDuplicates(duplicates)
	//两数之和
	two := []int{2, 7, 11, 15}
	TwoSum(two, 18)
}

func SingleNumber(nums []int) {
	countMap := make(map[int]int)
	for _, num := range nums {
		count, exists := countMap[num]
		if exists {
			countMap[num] = count + 1
		} else {
			countMap[num] = 1
		}
	}
	for num, count := range countMap {
		if count == 1 {
			fmt.Printf("只出现一次的数字: %d\n\n", num)
		}
	}
}

func IsPalindrome(num int) {
	numStr := strconv.Itoa(num)
	length := len(numStr)
	for i := 0; i < length/2; i++ {
		if numStr[i] != numStr[length-1-i] {
			fmt.Println("不是回文数: ", num)
			return
		}
	}
	fmt.Println("是回文数: ", num)
}

func IsValid(s string) {
	length := len(s)
	if length%2 != 0 {
		fmt.Println("不是有效的括号")
		return
	}
	pair := map[string]string{"[": "]", "(": ")", "{": "}"}
	stack := make([]string, 0)
	for i := 0; i < length; i++ {
		char := string(s[i])
		// 左括号入栈
		if _, exists := pair[char]; exists {
			stack = append(stack, char)
		} else { // 右括号
			if len(stack) == 0 {
				fmt.Println("不是有效的括号")
				return
			}
			top := pair[stack[len(stack)-1]]
			if char != top {
				fmt.Println("不是有效的括号")
				return
			}
			stack = stack[:len(stack)-1]
		}
	}
	// 栈必须为空，否则还有未匹配的左括号
	if len(stack) == 0 {
		fmt.Println("是有效的括号")
	} else {
		fmt.Println("不是有效的括号")
	}
}

func Merge(intervals [][]int) {

}

func LongestCommonPrefix(strs []string) {
	if len(strs) == 0 {
		fmt.Println("字符串为空, 请重新输入！")
		return
	}
	prefix := strs[0]
	for i := 1; i < len(strs); i++ {
		for !isPrefix(strs[i], prefix) {
			if len(prefix) == 0 {
				fmt.Println("没有最长公共前缀")
				return
			} else {
				prefix = prefix[:len(prefix)-1]
			}
		}
	}
	if len(prefix) == 0 {
		fmt.Println("没有最长公共前缀")
		return
	}
	fmt.Println("最长公共前缀是：", prefix)
}

func isPrefix(s string, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	} else {
		return s[:len(prefix)] == prefix
	}

}

func AddOne(strs []int) {
	for i := len(strs) - 1; i >= 0; i-- {
		if strs[i] == 9 {
			strs[i] = 0
		} else {
			strs[i]++
			fmt.Println("加1后的值为:", strs)
			return
		}
	}
	strs = append([]int{1}, strs...)
	fmt.Println("加1后的值为:", strs)
}

func RemoveDuplicates(nums []int) {
	if len(nums) == 0 {
		return
	}
	newNums := make([]int, 0)
	i := 0 // 慢指针，记录唯一元素的位置
	for j := 1; j < len(nums); j++ {
		if nums[j] != nums[i] {
			i++
			nums[i] = nums[j] // 把新元素放到 i 的位置
		}
	}
	newNums = nums[:i+1]
	fmt.Println("去重后的数组为:", newNums)
}

func TwoSum(nums []int, target int) {
	result := make(map[int]int) //key为数组元素，value为target-nums[i]，差值
	for _, num := range nums {
		result[num] = target - num
	}
	for _, num := range nums {
		value := result[num]
		if _, exists := result[value]; exists {
			fmt.Println("两数之和的两个数是:", num, result[num])
			return
		}
	}
}
