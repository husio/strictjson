package strictjson

import (
	"encoding/json"
	"reflect"
	"testing"
)

type User struct {
	FirstName string
	LastName  *string
	Age       int64   `json:"age,omitempty"`
	FavNumber float64 `json:",omitempty"`
	Address   *Address
}

type Address struct {
	Street     string
	HouseNo    *uint64
	PostalCode *string
	City       *City
}

type City struct {
	Name string
}

type Box struct {
	Name  string  `json:"name"`
	Items []Item  `json:"items"`
	Extra []*Item `json:"extra"`
}

type Item struct {
	Name   string `json:"name"`
	Amount *int64 `json:"amount"`
}

func TestSuccess(t *testing.T) {
	var testCases = []struct {
		Desc     string
		JSON     string
		Expected interface{}
	}{
		{
			Desc: "All values provided",
			JSON: `
				{
					"FirstName": "Bob",
					"age": 32,
					"LastName": "Doe",
					"FavNumber": 92,
					"Address": {
						"Street": "Wide Street",
						"HouseNo": 51,
						"PostalCode": "11EAS",
						"City": {
							"Name": "Berlin"
						}
					}
				}
			`,
			Expected: User{
				FirstName: "Bob",
				LastName:  stringPtr("Doe"),
				Age:       32,
				FavNumber: 92,
				Address: &Address{
					Street:     "Wide Street",
					HouseNo:    uint64Ptr(51),
					PostalCode: stringPtr("11EAS"),
					City: &City{
						Name: "Berlin",
					},
				},
			},
		},
		{
			Desc: "Only required values are provided",
			JSON: `
			{
				"FirstName": "bob",
				"age": 55,
				"FavNumber": 0
			}
			`,
			Expected: User{
				FirstName: "bob",
				Age:       55,
				FavNumber: 0,
			},
		},
		{
			Desc: "Optional fields provided with zero value",
			JSON: `
				{
					"FirstName": "Bob",
					"LastName": "",
					"age": 99,
					"FavNumber": 0,
					"Address": {
						"Street": "Any Street",
						"HouseNo": 0,
						"PostalCode": "",
						"City": {
							"Name": "New York"
						}
					}
				}
			`,
			Expected: User{
				FirstName: "Bob",
				LastName:  stringPtr(""),
				Age:       99,
				FavNumber: 0,
				Address: &Address{
					Street:     "Any Street",
					HouseNo:    uint64Ptr(0),
					PostalCode: stringPtr(""),
					City: &City{
						Name: "New York",
					},
				},
			},
		},
		{
			Desc: "Array with all elements",
			JSON: `
			{
				"name": "whatever",
				"items": [
					{"name": "first", "amount": 4},
					{"name": "second"}
				],
				"extra": [
					{"name": "first"}
				]
			}
			`,
			Expected: Box{
				Name: "whatever",
				Items: []Item{
					{Name: "first", Amount: int64Ptr(4)},
					{Name: "second"},
				},
				Extra: []*Item{
					{Name: "first"},
				},
			},
		},
		{
			Desc: "No value for array of pointers",
			JSON: `
			{
				"name": "whatever",
				"items": [
					{"name": "first", "amount": 4}
				]
			}
			`,
			Expected: Box{
				Name: "whatever",
				Items: []Item{
					{Name: "first", Amount: int64Ptr(4)},
				},
			},
		},
	}

	for _, tc := range testCases {
		result := reflect.New(reflect.TypeOf(tc.Expected))
		if err := Unmarshal([]byte(tc.JSON), result.Interface()); err != nil {
			t.Error("          failed:", tc.Desc)
			t.Error("cannot unmarshal:", err.String())
			continue
		}
		if !reflect.DeepEqual(tc.Expected, result.Elem().Interface()) {
			t.Error("  failed:", tc.Desc)
			t.Error("expected:", asJSON(tc.Expected))
			t.Error("  result:", asJSON(result.Elem().Interface()))
		}
	}
}

func asJSON(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func stringPtr(s string) *string {
	return &s
}

func int64Ptr(i int64) *int64 {
	return &i
}

func uint64Ptr(ui uint64) *uint64 {
	return &ui
}
