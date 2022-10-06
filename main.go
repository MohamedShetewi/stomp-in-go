package main

type Object struct {
	name string
}

type Arr struct {
	arr []*Object
}

func main() {
	obj1 := &Object{
		name: "Ahmed1",
	}
	obj2 := &Object{
		name: "Ahmed2",
	}
	obj3 := &Object{
		name: "Ahmed3",
	}
	obj4 := &Object{
		name: "Ahmed4",
	}
	arr := make([]*Object, 4)
	arr[0] = obj1
	arr[1] = obj2
	arr[2] = obj3
	arr[3] = obj4

}
