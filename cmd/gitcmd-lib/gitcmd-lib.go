package main

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"

import (
	"anybakup/cmd"
	"fmt"
	"unsafe"
)

// C-exportable wrapper for GetFileLog
//export GetFileLogC
func GetFileLogC(filePath *C.char) *C.char {
	if filePath == nil {
		return C.CString("error: file path is nil")
	}

	goFilePath := C.GoString(filePath)
	logs, err := cmd.GetFileLog(goFilePath)
	if err != nil {
		return C.CString(fmt.Sprintf("error: %v", err))
	}

	// Convert logs to JSON string for C compatibility
	result := fmt.Sprintf("logs: %d files", len(logs))
	return C.CString(result)
}

// C-exportable wrapper for AddFile
//export AddFileC
func AddFileC(filePath *C.char) *C.char {
	if filePath == nil {
		return C.CString("error: file path is nil")
	}

	goFilePath := C.GoString(filePath)
	result := cmd.AddFile(goFilePath)
	if result.Err != nil {
		return C.CString(fmt.Sprintf("error: %v", result.Err))
	}

	return C.CString(fmt.Sprintf("success: dest=%s, result=%v", result.Dest, result.Result))
}

// C-exportable wrapper for GetFile
//export GetFileC
func GetFileC(filePath *C.char, commit *C.char, target *C.char) *C.char {
	if filePath == nil || target == nil {
		return C.CString("error: file path or target is nil")
	}

	goFilePath := C.GoString(filePath)
	goCommit := C.GoString(commit)
	goTarget := C.GoString(target)

	err := cmd.GetFile(goFilePath, goCommit, goTarget)
	if err != nil {
		return C.CString(fmt.Sprintf("error: %v", err))
	}

	return C.CString("success: file retrieved")
}

// Helper function to free C strings
//export FreeString
func FreeString(str *C.char) {
	if str != nil {
		C.free(unsafe.Pointer(str))
	}
}

func main() {
	// Main function is required for c-shared buildmode but won't be called when used as a library
}