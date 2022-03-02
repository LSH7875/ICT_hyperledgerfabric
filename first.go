package main

import (
	"fmt"
	"math"
)

type Shape interface {
	area() float64
	perimeter() float64
}

type Rect struct { //Rect 정의
	width, height float64
}
type Circle struct { //Circle 정의
	radius float64
}

//Rect 타입에 대한 Shape 인터페이스 구현
func (r Rect) area() float64 {
	return r.width * r.height
}
func (r Rect) perimeter() float64 {
	return 2 * (r.width + r.height)
}

//Circle 타입에 대한 Shape 인터페이스 구현
func (c Circle) area() float64 {
	return math.Pi * c.radius * c.radius
}
func (c Circle) perimeter() float64 {
	return 2 * math.Pi * c.radius
}

func main() {
	r := Rect{10., 20.}
	c := Circle{10}
	showArea(r, c)
}

func showArea(shapes ...Shape) {
	for _, s := range shapes { //순차적으로 shape를 받은 값을 돌리겠다. 첫번째에 인덱스 들어간다. 값도 들어가고, 그래서 인덱스만 쓸거면 index,_로 쓰고 인덱스 안 쓸거면 _,index쓰는 거임
		//메모리 낭비를 줄이기 위해서
		a := s.area() //인터페이스 메서드 호출
		fmt.Println(a)
	}
}
