package ferror

import "strings"

func ErrIsImportant(e error) bool {
	if strings.Contains(e.Error(),"part of a ReadProcessMemory or WriteProcessMemory request was completed") {
		return false // error occurs when we read a dead player
	} else if strings.Contains(e.Error(),"The operation completed successfully") {
		return false
	} else if e != nil {
		return true
	}
	return false
}