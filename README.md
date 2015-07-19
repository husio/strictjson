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

Although package provides strict type checking, it does not validate
unmarshaled content. Any data validation must be done manually.


## TODO

* [ ] Tests
* [ ] Allow root value to be any type, not only `struct`
* [ ] Documentation & examples
* [ ] Basic validation through tags
* [ ] Benchmarks
