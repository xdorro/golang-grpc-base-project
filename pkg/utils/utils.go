package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// MakeDir creates a directory if it does not exist.
func MakeDir(dir string) error {
	dirConvert := filepath.Dir(dir)
	if !Exists(dirConvert) {
		err := os.Mkdir(dirConvert, 0o775)
		if err != nil {
			return err
		}
	}

	return nil
}

// TotalPage returns the total number of pages.
func TotalPage(total int64, pageSize int64) (totalPage int64) {
	if total%pageSize == 0 {
		totalPage = total / pageSize
	} else {
		totalPage = total/pageSize + 1
	}

	if totalPage == 0 {
		totalPage = 1
	}

	return
}

// CurrentPage returns the current page.
func CurrentPage(page int64, totalPages int64) int64 {
	if page <= 0 || totalPages < page {
		page = 1
	} else if page > totalPages {
		page = totalPages
	}

	return page
}

// StringCompareOrPassValue returns string compare or pass value
func StringCompareOrPassValue(a, b string) string {
	if strings.Compare(a, b) != 0 {
		return b
	}

	return a
}
