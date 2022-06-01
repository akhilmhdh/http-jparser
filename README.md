[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/akhilmhdh/http-jparser">
    <h3 align="center">Go HTTP JParser</h3>  
  </a>

  <p align="center">
    A simple JSON parser and encoding for your Go REST servers.
    <br />
    <a href="https://github.com/akhilmhdh/http-jparser/issues">Report Bug</a>
    ·
    <a href="https://github.com/akhilmhdh/http-jparser/issues">Request Feature</a>
  </p>
</div>

`http-jparser` is a simple Go library used for parsing and validating JSON request body and JSON responding. The parser is compatible with any framework that exposes the `HTTP` interface like [go-chi](https://github.com/go-chi/chi). Under the hood, the library uses [go-json](https://github.com/goccy/go-json) for JSON encoding/decoding and [go-validator](https://github.com/go-playground/validator) for validation.

## Install

```bash
go get -u github.com/akhilmhdh/http-jparser
```

## Features

- Simple API
- Request body validation
- Compatibility with `net/http`
- Uses go-json for fast encoding/decoding JSON

## Examples

The following examples uses `go-chi`

```go

package main

import (
	"net/http"

	jparser "github.com/akhilmhdh/http-jparser"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Response struct {
	HTTPStatusCode int         `json:"-"`               // http response status code
	Success        bool        `json:"success"`         // flag to indicated failed to success
	Message        string      `json:"message"`         // user-level status message
	AppCode        int64       `json:"code,omitempty"`  // application-specific error code
	ErrorText      interface{} `json:"error,omitempty"` // application-level error message, for debugging
	Data           interface{} `json:"data,omitempty"`
}

func (e *Response) GetStatusCode() int {
	return e.HTTPStatusCode
}

type ReqBody struct {
	Key   string `json:"key" validate:"lowercase"`
	Value string `json:"value"`
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		body := &Response{
			HTTPStatusCode: 200,
			Success:        true,
			Message:        "Success Request Message",
		}

		jparser.Send(w, r, body)
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		data := &ReqBody{}

		// GET JSON request body
		if err := jparser.Get(r, data); err != nil {
			// To check if error is validation error or not
			// if vErr, ok := err.(*jparser.ValidationErrors); ok {}
			// ErrInvalidRequest can be a wrapper on Response to return error messages
			jparser.Send(w, r, ErrInvalidRequest())
			return
		}

		// success response
		jparser.Send(w, r, &Response{
			HTTPStatusCode: 200,
			Message:        "Data returned successfully",
			Success:        true,
			Data:           data,
		})
	})

	http.ListenAndServe(":3000", r)
}

```

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

[contributors-shield]: https://img.shields.io/github/contributors/akhilmhdh/http-jparser.svg?style=for-the-badge
[contributors-url]: https://github.com/akhilmhdh/http-jparser/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/akhilmhdh/http-jparser.svg?style=for-the-badge
[forks-url]: https://github.com/akhilmhdh/http-jparser/network/members
[stars-shield]: https://img.shields.io/github/stars/akhilmhdh/http-jparser.svg?style=for-the-badge
[stars-url]: https://github.com/akhilmhdh/http-jparser/stargazers
[issues-shield]: https://img.shields.io/github/issues/akhilmhdh/http-jparser.svg?style=for-the-badge
[issues-url]: https://github.com/akhilmhdh/http-jparser/issues
[license-shield]: https://img.shields.io/github/license/akhilmhdh/http-jparser.svg?style=for-the-badge
[license-url]: https://github.com/akhilmhdh/http-jparser/blob/master/LICENSE
