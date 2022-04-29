package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

type Code struct {
	Num  int
	Code string
}

func TestCache(t *testing.T) {
	var num1 float64 = 123456.45
	data, err := json.Marshal(num1)
	if err != nil {
		fmt.Printf("序列化错误 err=%v\n", err)
	}
	fmt.Printf("float64 序列化后=%v\n", string(data))
	// int test
	fmt.Println("----------int test")
	ctx := context.Background()
	c := NewCache(50, "test", "int")
	for i := 0; i < 10; i++ {
		c.Set(ctx, fmt.Sprint(i), i)
	}
	var a int
	err = c.Get(ctx, "5", &a)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(a)
	c.Set(ctx, "5", 10000000)
	var ttt []int
	err = c.List(ctx, &ttt)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, kk := range ttt {
		fmt.Println(kk)
	}
	fmt.Println("keys len:", c.Len(ctx))

	fmt.Println("----------string test")
	// int test
	c = NewCache(5, "test", "string")
	for i := 0; i < 10; i++ {
		c.Set(ctx, fmt.Sprint(i), fmt.Sprint(i))
	}
	var as string
	err = c.Get(ctx, "5", &as)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(as)
	c.Set(ctx, "5", fmt.Sprint(10000000))
	var ttts []string
	err = c.List(ctx, &ttts)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, kk := range ttts {
		fmt.Println(kk)
	}
	fmt.Println("keys len:", c.Len(ctx))

	//
	fmt.Println("----------struct test")
	c = NewCache(5, "test", "struct")
	for i := 0; i < 10; i++ {
		c.Set(ctx, fmt.Sprint(i), &Code{Num: i})
	}
	var aa Code
	err = c.Get(ctx, "5", &aa)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(aa.Num)
	fmt.Println(c.Exist(ctx, "5"))
	c.Delete(ctx, "6")
	c.Set(ctx, "5", &Code{Num: 10000000})
	var tttaa []Code
	err = c.List(ctx, &tttaa)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, kk := range tttaa {
		fmt.Println(kk.Num)
	}
	fmt.Println("keys len:", c.Len(ctx))
	fmt.Println(c.Exist(ctx, "5"))
	fmt.Println(c.Exist(ctx, "6"))
}
