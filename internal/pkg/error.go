package pkg

import "fmt"

type CustomError struct {
	Code    int    `json:"code"`
	Service string `json:"service"`
	Cause   error  `json:"detail"`
}

func (err *CustomError) Error() interface{} {
	return map[string]interface{}{
		"message": fmt.Sprintf("Error at %s : %s", err.Service, err.Cause.Error()),
	}
}
