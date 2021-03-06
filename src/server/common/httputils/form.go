package httputils

import (
	"fmt"
	"net/http"
	"path/filepath"
	"server/common/tool"
	"strconv"
	"strings"
)

// BoolValue transforms a form value in different formats into a boolean type.
func BoolValue(r *http.Request, k string) bool {
	s := strings.ToLower(strings.TrimSpace(r.FormValue(k)))
	return !(s == "" || s == "0" || s == "no" || s == "false" || s == "none")
}

// BoolValueOrDefault returns the default bool passed if the query param is
// missing, otherwise it's just a proxy to boolValue above
func BoolValueOrDefault(r *http.Request, k string, d bool) bool {
	if _, ok := r.Form[k]; !ok {
		return d
	}
	return BoolValue(r, k)
}

// Int64ValueOrZero parses a form value into an int64 type.
// It returns 0 if the parsing fails.
func Int64ValueOrZero(r *http.Request, k string) int64 {
	val, err := Int64ValueOrDefault(r, k, 0)
	if err != nil {
		return 0
	}
	return val
}

// Int64ValueOrDefault parses a form value into an int64 type. If there is an
// error, returns the error. If there is no value returns the default value.
func Int64ValueOrDefault(r *http.Request, field string, def int64) (int64, error) {
	if r.Form.Get(field) != "" {
		value, err := strconv.ParseInt(r.Form.Get(field), 10, 64)
		if err != nil {
			return value, err
		}
		return value, nil
	}
	return def, nil
}

// IntValueOrZero parses a form value into an int type.
// It returns 0 if the parsing fails.
func IntValueOrZero(r *http.Request, k string) int {
	val, err := IntValueOrDefault(r, k, 0)
	if err != nil {
		return 0
	}
	return val
}

// IntValueOrDefault parses a form value into an int type. If there is an
// error, returns the error. If there is no value returns the default value.
func IntValueOrDefault(r *http.Request, field string, def int) (int, error) {
	//		fmt.Println(r.FormValue("offset"))
	//	fmt.Println(field, r.Form.Get(field))
	if r.Form.Get(field) != "" {
		value, err := strconv.Atoi(r.Form.Get(field))
		if err != nil {
			return value, err
		}
		return value, nil
	}
	return def, nil
}

// IntValueOrDefault parses a form value into an int type. If there is an
// error, returns the error. If there is no value returns the default value.
func IntValueOrDefaultIfErrZero(r *http.Request, field string, def int) int {
	if r.Form.Get(field) != "" {
		value, err := strconv.Atoi(r.Form.Get(field))
		if err != nil {
			return 0
		}
		return value
	}
	return def
}

//func IntArray(r *http.Request, field string, def int) []int {
//	if r.Form.Get(field) != "" {
//		value, err := strconv.Atoi(r.Form.Get(field))
//		if err != nil {
//			return 0
//		}
//		return value
//	}
//	return def
//}
func StringArray(r *http.Request, field string) []string {
	arrstring := r.FormValue(field)
	if arrstring != "" {
		arr := strings.Split(arrstring, ",")
		if len(arr) != 0 {
			return arr
		} else {
			return nil
		}
	}
	return nil

}

// ArchiveOptions stores archive information for different operations.
type ArchiveOptions struct {
	Name string
	Path string
}

// ArchiveFormValues parses form values and turns them into ArchiveOptions.
// It fails if the archive name and path are not in the request.
func ArchiveFormValues(r *http.Request, vars map[string]string) (ArchiveOptions, error) {
	if err := ParseForm(r); err != nil {
		return ArchiveOptions{}, err
	}

	name := vars["name"]
	path := filepath.FromSlash(r.Form.Get("path"))

	switch {
	case name == "":
		return ArchiveOptions{}, fmt.Errorf("bad parameter: 'name' cannot be empty")
	case path == "":
		return ArchiveOptions{}, fmt.Errorf("bad parameter: 'path' cannot be empty")
	}

	return ArchiveOptions{name, path}, nil
}

func GetPagination(r *http.Request) tool.Meta {
	offset := IntValueOrZero(r, "offset")
	//	fmt.Println("offset", offset)
	limit, _ := IntValueOrDefault(r, "limit", 100)
	total := IntValueOrZero(r, "total")
	return tool.Meta{Offset: offset, Limit: limit, Total: total}
}
