package main

/*
#include <stdlib.h>
#include <string.h>

typedef struct {
	char* commit;
	char* author;
	char* date;
	char* message;
} GitChange;

typedef struct {
	GitChange* changes;
	int count;
} GitChangeArray;
*/
import "C"

import (
	"anybakup/cmd"
	"fmt"
	"unsafe"
)

// C-exportable wrapper for GetFileLog
//export GetFileLogC
func GetFileLogC(filePath *C.char) *C.GitChangeArray {
	if filePath == nil {
		return nil
	}

	goFilePath := C.GoString(filePath)
	logs, err := cmd.GetFileLog(goFilePath)
	if err != nil {
		return nil
	}

	if len(logs) == 0 {
		return nil
	}

	// Allocate C array
	array := (*C.GitChangeArray)(C.malloc(C.size_t(unsafe.Sizeof(C.GitChangeArray{}))))
	count := len(logs)
	*array = C.GitChangeArray{}
	array.count = C.int(count)

	// Allocate memory for the array of GitChange structs
	changesPtr := C.malloc(C.size_t(count) * C.size_t(unsafe.Sizeof(C.GitChange{})))
	array.changes = (*C.GitChange)(changesPtr)

	// Convert each log entry to C struct
	for i, log := range logs {
		change := (*C.GitChange)(unsafe.Pointer(uintptr(changesPtr) + uintptr(i)*unsafe.Sizeof(C.GitChange{})))
		change.commit = C.CString(log.Commit)
		change.author = C.CString(log.Author)
		change.date = C.CString(log.Date)
		change.message = C.CString(log.Message)
	}

	return array
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

// Helper function to free GitChangeArray
//export FreeGitChangeArray
func FreeGitChangeArray(array *C.GitChangeArray) {
	if array == nil {
		return
	}

	// Free each GitChange struct's strings
	if array.changes != nil {
		for i := 0; i < int(array.count); i++ {
			change := (*C.GitChange)(unsafe.Pointer(uintptr(unsafe.Pointer(array.changes)) + uintptr(i)*unsafe.Sizeof(C.GitChange{})))
			if change.commit != nil {
				C.free(unsafe.Pointer(change.commit))
			}
			if change.author != nil {
				C.free(unsafe.Pointer(change.author))
			}
			if change.date != nil {
				C.free(unsafe.Pointer(change.date))
			}
			if change.message != nil {
				C.free(unsafe.Pointer(change.message))
			}
		}
		// Free the array of GitChange structs
		C.free(unsafe.Pointer(array.changes))
	}

	// Free the GitChangeArray struct itself
	C.free(unsafe.Pointer(array))
}

func main() {
	// Main function is required for c-shared buildmode but won't be called when used as a library
}