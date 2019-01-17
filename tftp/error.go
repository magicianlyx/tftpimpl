package tftp

import (
	"runtime"
	"fmt"
)

func HandleError(err error) (error) {
	if err == nil {
		return nil
	}
	location := ""
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		function := runtime.FuncForPC(pc)
		location += fmt.Sprintf("FileName:%s, Line:%d, Func:%s", file, line, function.Name())
	}
	//fmt.Printf("[ERROR] | location{%s} | error msg{%s}", location, err.Error())
	fmt.Println(fmt.Sprintf("[ERROR] | location{%s} | error msg{%s}", location, err.Error()))
	return err
}
