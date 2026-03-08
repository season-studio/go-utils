package arguments

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/season-studio/go-utils/misc"
)

type ArgumentEntry struct {
	Prefixs  []string
	Desc     string
	required bool
	v        any
}

func MakeSimpleEntry(pv any, prefixs ...string) *ArgumentEntry {
	return &ArgumentEntry{
		Prefixs: prefixs,
		v:       pv,
	}
}

type ArgumentRequiredError struct {
	Requireds []*ArgumentEntry
}

func (e *ArgumentRequiredError) Error() string {
	if len(e.Requireds) > 1 {
		return fmt.Sprintf("%d arguments are required", len(e.Requireds))
	} else {
		return "1 argument is required"
	}
}

var (
	locker      sync.Mutex
	argsEntries []*ArgumentEntry
	UnknownArgs []string
)

func RegisterEntry(v any, required bool, desc string, prefixs ...string) error {
	switch v.(type) {
	case *bool, *string,
		*int, *int8, *int16, *int32, *int64,
		*uint, *uint8, *uint16, *uint32, *uint64,
		*float32, *float64,
		func(), func(string):
	default:
		vType := reflect.TypeOf(v)
		if vType.Kind() != reflect.Ptr {
			return fmt.Errorf("the \"v\" parameter must be a pointer")
		}
		vType = vType.Elem()
		vVal := reflect.ValueOf(v).Elem()
		if vVal.Kind() != reflect.Struct {
			return fmt.Errorf("the typeof \"v\" parameter is unsupported")
		}
		fieldCount := vType.NumField()
		for i := 0; i < fieldCount; i++ {
			field := vVal.Field(i)

			// 跳过不可设置的字段
			if !field.CanSet() {
				continue
			}

			// 跳过不可寻址字段（如未导出字段）
			if !field.CanAddr() {
				continue
			}

			fieldType := vType.Field(i)
			tagsStr := fieldType.Tag.Get("arg_tags")
			if misc.IsStrEmptyAndWhitespace(tagsStr) {
				continue
			}
			tags := strings.Split(tagsStr, " ")
			argDesc := fieldType.Tag.Get("arg_desc")
			argRequired := fieldType.Tag.Get("arg_required") == "true"
			fieldPtr := field.Addr().Interface()
			RegisterEntry(fieldPtr, argRequired, argDesc, tags...)
		}
	}

	if len(prefixs) <= 0 {
		return fmt.Errorf("require at least one prefix")
	}

	if prefixs[0] == "" {
		return nil
	}

	locker.Lock()
	defer locker.Unlock()

	argsEntries = append(argsEntries, &ArgumentEntry{
		Prefixs:  prefixs,
		Desc:     desc,
		required: required,
		v:        v,
	})
	return nil
}

func ShowHelper(versionInfo string) {
	if len(versionInfo) > 0 {
		fmt.Println(versionInfo)
	}
	for _, entry := range argsEntries {
		if len(entry.Desc) == 0 {
			continue
		}
		var valType string
		switch entry.v.(type) {
		case *bool, func():
			valType = ""
		case func(string):
			valType = " <val>"
		default:
			valType = fmt.Sprintf(" <%s>", reflect.ValueOf(entry.v).Elem().Type().Name())
		}
		for idx, prefix := range entry.Prefixs {
			if idx > 0 {
				fmt.Print(", ")
			}
			fmt.Print(prefix)
			if len(valType) > 0 {
				fmt.Print(valType)
			}
		}
		fmt.Print("\n")
		if entry.required {
			fmt.Println("\tRequired")
		}
		if len(entry.Desc) > 0 {
			fmt.Print("\t")
			fmt.Println(strings.ReplaceAll(strings.TrimSpace(entry.Desc), "\n", "\n\t"))
		}
	}
}

func ParseBaseOnEntries(args []string, entries ...*ArgumentEntry) error {
	return ParseFromListBaseOnEntries(args, entries)
}

func ParseFromListBaseOnEntries(args []string, entries []*ArgumentEntry) (retErr error) {
	defer func() {
		if err := recover(); err != nil {
			if rerr, ok := err.(error); ok {
				retErr = rerr
			} else {
				retErr = fmt.Errorf("%v", err)
			}
		}
	}()

	parsedList := make([]bool, len(entries))

	lastIdx := len(args) - 1
	pass := true
	for idx, arg := range args {
		if pass {
			pass = false
			continue
		}
		eidx, entry := misc.Find(entries, func(e *ArgumentEntry) bool {
			return misc.IndexOf(e.Prefixs, arg) >= 0
		})
		if entry != nil {
			switch v := entry.v.(type) {
			case *bool:
				*v = true
			case func():
				v()
			case *string:
				if idx < lastIdx {
					*v = args[idx+1]
					pass = true
				}
			case func(string):
				if idx < lastIdx {
					v(args[idx+1])
					pass = true
				}
			case *int, *int8, *int16, *int32, *int64,
				*uint, *uint8, *uint16, *uint32, *uint64,
				*float32, *float64:
				if idx < lastIdx {
					if err := misc.ParseNumber(args[idx+1], v); err != nil {
						return fmt.Errorf("bad value of %s (%v)", arg, err)
					}
					pass = true
				}
			default:
				return fmt.Errorf("cannot parse argument for type of %T", entry.v)
			}
			parsedList[eidx] = true
		} else {
			UnknownArgs = append(UnknownArgs, arg)
		}
	}

	expectList := make([]*ArgumentEntry, 0, len(entries))
	for idx, parsed := range parsedList {
		if entries[idx].required && !parsed {
			expectList = append(expectList, entries[idx])
		}
	}

	if len(expectList) == 0 {
		return nil
	} else {
		return &ArgumentRequiredError{
			Requireds: expectList,
		}
	}
}

func ParseFromList(args []string) error {
	locker.Lock()
	defer locker.Unlock()

	return ParseFromListBaseOnEntries(args, argsEntries)
}

func Parse() error {
	return ParseFromList(os.Args)
}
