package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
)

func main() {
	test := &Student{
		Name:   "Xiluo",
		Male:   true,
		Scores: []int32{98, 85, 88},
	}
	data, err := proto.Marshal(test)
	fmt.Println(data)
	if err != nil {
		log.Fatal("marshaling error:", err)
	}
	newTest := &Student{}
	err = proto.Unmarshal(data, newTest)
	fmt.Printf("%+v \n", newTest)
	if err != nil {
		log.Fatal("unmarshaling error:", err)
	}

	if test.GetName() != newTest.GetName() {
		log.Fatal("data mismatch %q", test.GetName(), newTest.GetName())
	}

}
