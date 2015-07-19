package strictjson

import (
	"reflect"
	"testing"
)

func TestDecodeBasicTypes(t *testing.T) {
	type User struct {
		FirstName string
		LastName  *string
		Age       int64   `json:"age,omitempty"`
		FavNumber float64 `json:",omitempty"`
	}

	expected := User{
		FirstName: "Bob",
		LastName:  nil,
		Age:       32,
		FavNumber: 0,
	}
	input := []byte(`
	{
		"FirstName": "Bob",
		"age": 32
	}
	`)

	var got User
	if err := Unmarshal(input, &got); err != nil {
		t.Fatalf("cannot unmarshal: %s", err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("got: %#v", got)
	}
}
