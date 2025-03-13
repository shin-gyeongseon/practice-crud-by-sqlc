package tutorial

import (
	"fmt"
	"testing"
)

var string1 string = string("abcd")
var string2 string = string("abcd")

func TestErrors(t *testing.T) {
	fmt.Println(&string1)
	fmt.Println(&string2)
	if string1 == string2 {
		fmt.Println("같은 타입입니다.")
	}
}
