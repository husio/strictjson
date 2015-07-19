# Strict JSON

Strict JSON package provides wrapper on top of standard library
`json.Unmarshal` with strict rules for data deserialization.


## Motivation

From [Unmarshal documentation](http://golang.org/pkg/encoding/json/#Unmarshal):

>  If a JSON value is not appropriate for a given target type, or if a JSON
>  number overflows the target type, Unmarshal skips that field and completes
>  the unmarshalling as best it can. If no more serious errors are
>  encountered, Unmarshal returns an UnmarshalTypeError describing the
>  earliest such error.

Strict JSON package always return error if provided value does not meet
structure expectations.


## Validation

All non-pointer fields are required. All pointer fields  and fields with
`json` tag `omitempty` are optional (value is not required).

Empty containers (slice, map) are considered empty value. Value for container
of pointers is optional (not required) while at least one element for
container of non pointer values must be provided.

To unmarshal JSON data, type compatible non zero value must be provided for
all required fields.


```go
type User struct {
    FirstName    string                        // required
    LastName     string    `json:",omitempty"` // not required
    Nickname     *string                       // not required
    Hobby        []*string                     // not required
    LuckyNumbers []int64                       // required
    FavColors    []string  `json:",omitempty"` // not required
}
```

Although package provides strict type and data existence checking, it does not
validate unmarshaled content. Any data validation must be done manually.


## TODO

* [ ] Tests
* [ ] Allow root value to be any type, not only `struct`
* [ ] Documentation & examples
* [ ] Basic validation through tags (`min`, `max`)
* [ ] Benchmarks & comparison with plain `encoding/json`
