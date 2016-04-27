package main

import "fmt"

func Test() {
}

func main() {
	var res [][]string
	for _, v := range res {
		v = make([]string, 3)
		v = append(v, "fff")
	}
	r := []string{"t","g","sd"}
	res = append(res, r)
	//res = append(res, []string{"f", "ok"})
	fmt.Println(res)

	k := []string{}
	k = append(k, "d", "g")
	fmt.Println(k)
}