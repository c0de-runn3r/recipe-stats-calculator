package processor

import (
	"fmt"
	"reflect"
	"testing"
)

func assert(t *testing.T, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%+v != %+v", a, b)
	}
}

func TestDeliveryTimes(t *testing.T) {
	open, close, err := deliveryTime("Saturday 10AM - 4PM")
	fmt.Println("ERROR: ", err)
	assert(t, open, 10)
	assert(t, close, 4)
	open1, close1, err := deliveryTime("Wednesday 4AM - 7PM")
	fmt.Println("ERROR: ", err)
	assert(t, open1, 4)
	assert(t, close1, 7)
	open2, close2, err := deliveryTime("Wednesday 11AM - 10PM")
	fmt.Println("ERROR: ", err)
	assert(t, open2, 11)
	assert(t, close2, 10)
}
