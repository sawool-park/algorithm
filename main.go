package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"unsafe"
)

func main() {
	keys := []int{
		10, 20, 60, 70, 80, 100, 110, 120, 160, 170, 200, 210, 220, 260, 270, 280,
		300, 310, 320, 360, 370, 380, 400, 410, 420, 460, 470, 480, 500, 510, 520, 560, 570, 580,
		600, 610, 620, 660, 670, 680, 700, 710, 720, 760, 770, 780, 800, 810, 820, 860, 870, 880}

	slice1 := keys[:]
	slice2 := make([]int, len(keys))
	copy(slice2, keys)

	fmt.Println(len(slice1), cap(slice1))
	fmt.Println(unsafe.Pointer(&keys[0]), unsafe.Pointer(&slice1[0]), unsafe.Pointer(&slice2[0]))

	t := NewBTree()
	for _, v := range keys {
		t.Insert(v)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("정수를 입력하세요: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		number, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("정수로 변환할 수 없습니다.")
			break
		}
		index, node := t.Query(number)
		if index >= 0 {
			fmt.Println("Query Found.. ", index, node.x, node)
		} else {
			fmt.Println("Query Not Found.. ", index, node.x, node)
		}
		index = sort.Search(len(keys), func(i int) bool {
			return keys[i] >= number
		})
		if index < len(keys) && keys[index] == number {
			fmt.Println("found: ", index, keys[index])
		} else {
			fmt.Println("not found: ", keys[index])
		}
	}
	fmt.Println("end")
}
