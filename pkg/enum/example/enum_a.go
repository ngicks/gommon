// Code generated by github.com/ngicks/gommon/pkg/enum/cmd/generate_enum/generate_enum.go. DO NOT EDIT.
package enum

import "os"

type EnumA[T any] struct {
	data any
}

func EnumAInt[T any](val int) EnumA[T] {
	return EnumA[T]{
		data: val,
	}
}

func EnumAString[T any](val string) EnumA[T] {
	return EnumA[T]{
		data: val,
	}
}

func EnumAOsFile[T any](val *os.File) EnumA[T] {
	return EnumA[T]{
		data: val,
	}
}

type EnumAMatcher[T any] struct {
	Int    func(int) T
	String func(string) T
	OsFile func(*os.File) T
	Any    func() T
}

func (e EnumA[T]) Match(m EnumAMatcher[T]) T {
	switch x := e.data.(type) {
	case int:
		if m.Int != nil {
			return m.Int(x)
		}
	case string:
		if m.String != nil {
			return m.String(x)
		}
	case *os.File:
		if m.OsFile != nil {
			return m.OsFile(x)
		}
	}

	if m.Any != nil {
		return m.Any()
	}

	var ret T
	return ret
}