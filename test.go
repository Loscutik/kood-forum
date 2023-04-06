package main

import (
	"database/sql"
	"fmt"
	"forum/model"

	_ "github.com/mattn/go-sqlite3"
)

func test() {

	var s[]int
	fmt.Printf("%#v, %d\n",s,len(s))


	fmt.Println("---------")

	foftf(tf)
	foftf2(tf)
	for _, driver := range sql.Drivers() {
		fmt.Println(driver)
	}
}


func tf(i int) string {
	return fmt.Sprintf("i=%d\n", i)
}

type tfunc func(int) string

type tfunc2 func(int) string

func foftf2(f tfunc2) {
	fmt.Println("foftf2: ", f(20))
	fmt.Println("foftf2: ", f(200))
}

func foftf(f tfunc) {
	fmt.Println("foftf: ", f(0))
	fmt.Println("foftf: ", f(10))
}

func shortTest(){
	var m map[int]string
	fmt.Printf("m= %#v\n",m)
	it,ok:=m[0]
	fmt.Printf("it= %#v\n",it)
	fmt.Printf("ok= %#v\n",ok)

	post := &model.Post{}
	post.Message.Author=&model.User{}
	post.Message.Author.Name="r"
	fmt.Printf("%p,  %#v\n", post.Message.Author,post.Message.Author)
}
