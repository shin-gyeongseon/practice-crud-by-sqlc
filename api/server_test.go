package api

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSer(t *testing.T) {
	err1 := errors.New("error1")
	err2 := err1
	fmt.Printf("err1 address: %v, err2 address: %v \n\n", &err1, &err2)

	resultBoolean := err1 == err2
	fmt.Println(resultBoolean)
	require.True(t, resultBoolean)
}

type Temporary struct {
	s string
}

func (t *Temporary) Error() string {
	return "temporaryr error function"
}

func TestPlainGo(t *testing.T) {
	tmp := Temporary{
		"ddffdsaf",
	}
	fmt.Println(tmp)
}

type customError struct {
	Err string
}

func (c customError) Error() string {
	return c.Err
}

func TestEqualError(t *testing.T) {
	err1 := customError{
		Err: "this is custom error",
	}

	err2 := customError{
		Err: "this is custom error",
	}

	err3 := &err1
	fmt.Println(*err3)

	compare_result := err1 == err2
	if compare_result {
		fmt.Println("같습니다.")
		fmt.Printf("err1: %v , err2: %v\n", &err1, &err2)
		// fmt.Printf("err1: %v , err2: %v\n", *err1, *err2)

		return
	}

	fmt.Println("다릅니다.")
}

type customError2 struct {
	Err          string
	OtherMessage string
}

func (c *customError2) ChangeOtherMessageField() {
	c.OtherMessage = "changing !! in other print pointer receiver !! "
}

func (c *customError2) OtherPrint() string {
	return c.OtherMessage
}

type Terrible int64

func (t *Terrible) RoyalPrint() {
	*t = 999
}

func TestEqualError2(t *testing.T) {
	err1 := customError2{
		Err:          "err11",
		OtherMessage: "initialize OtherMessage Field",
	}
	err2 := err1
	err3 := err1

	fmt.Printf("%p\n", &err1)
	fmt.Println(&err2)
	fmt.Printf("%p\n", &err3)

	err1.ChangeOtherMessageField()

	fmt.Println("err1 ::: ", err1.OtherMessage)
	fmt.Println("err2 ::: ", err2.OtherMessage)
	fmt.Println("err3 ::: ", err3.OtherMessage)
}