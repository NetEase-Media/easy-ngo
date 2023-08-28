package xgo

import (
	"fmt"
	"runtime"

	"github.com/NetEase-Media/easy-ngo/utils/xstring"
	"github.com/pkg/errors"
)

func SerialUntilError(fns ...func() error) func() error {
	return func() error {
		for _, fn := range fns {
			if err := try(fn, nil); err != nil {
				return err
			}
		}
		return nil
	}
}

func try(fn func() error, cleaner func()) (ret error) {
	if cleaner != nil {
		defer cleaner()
	}
	defer func() {
		if err := recover(); err != nil {
			_, _, line, _ := runtime.Caller(2)
			if _, ok := err.(error); ok {
				ret = err.(error)
			} else {
				ret = fmt.Errorf("%+v", err)
			}
			ret = errors.Wrap(ret, fmt.Sprintf("%s:%d", xstring.FunctionName(fn), line))
		}
	}()
	return fn()
}
