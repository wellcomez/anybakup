package main

/*
#include <stdlib.h>
#include <string.h>
#include <stdint.h>

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

typedef struct {
	int64_t id;
	char* src_file;
	char* dest_file;
	int is_file;  // 1 for true, 0 for false
	int revcount;
	int sub;
	char* tag;
	char* add_time;
	char* update_time;
} FileOperationC;

typedef struct {
	FileOperationC* operations;
	int count;
} FileOperationArray;
*/
import "C"

import (
	"fmt"
	"unsafe"

	"anybakup/cmd"
	"anybakup/util"
	// "fmt"
)

// C-exportable wrapper for GetFileLog
//
//export GetFileLogC
func GetFileLogC(profilename *C.char, filePath *C.char) *C.GitChangeArray {
	if filePath == nil {
		return nil
	}
	goProfilename := C.GoString(profilename)
	goFilePath := C.GoString(filePath)
	g := cmd.NewGitCmd(goProfilename)
	logs, err := g.GetFileLog(util.RepoPath(goFilePath))
	if err != nil {
		return nil
	}

	if len(logs) == 0 {
		return nil
	}

	// Allocate C array
	array := (*C.GitChangeArray)(C.malloc(C.size_t(unsafe.Sizeof(C.GitChangeArray{}))))
	if array == nil {
		return nil
	}
	count := len(logs)
	*array = C.GitChangeArray{}
	array.count = C.int(count)

	// Allocate memory for the array of GitChange structs
	changesPtr := C.malloc(C.size_t(count) * C.size_t(unsafe.Sizeof(C.GitChange{})))
	if changesPtr == nil {
		C.free(unsafe.Pointer(array))
		return nil
	}
	array.changes = (*C.GitChange)(changesPtr)

	// Convert each log entry to C struct
	for i, log := range logs {
		change := (*C.GitChange)(unsafe.Pointer(uintptr(changesPtr) + uintptr(i)*unsafe.Sizeof(C.GitChange{})))
		change.commit = C.CString(log.Commit)
		change.author = C.CString(log.Author)
		change.date = C.CString(log.Date)
		change.message = C.CString(log.Message)

		// Check if any CString allocation failed
		if change.commit == nil || change.author == nil || change.date == nil || change.message == nil {
			// Free already allocated strings and memory
			for j := range i {
				prevChange := (*C.GitChange)(unsafe.Pointer(uintptr(changesPtr) + uintptr(j)*unsafe.Sizeof(C.GitChange{})))
				if prevChange.commit != nil {
					C.free(unsafe.Pointer(prevChange.commit))
				}
				if prevChange.author != nil {
					C.free(unsafe.Pointer(prevChange.author))
				}
				if prevChange.date != nil {
					C.free(unsafe.Pointer(prevChange.date))
				}
				if prevChange.message != nil {
					C.free(unsafe.Pointer(prevChange.message))
				}
			}
			C.free(unsafe.Pointer(changesPtr))
			C.free(unsafe.Pointer(array))
			return nil
		}
	}

	return array
}

// C-exportable wrapper for BackupOptAdd
//
//export BackupOptAddC
// func BackupOptAddC(srcFile *C.char, destFile *C.char, isFile C.int) C.int {
// 	if srcFile == nil || destFile == nil {
// 		return -1
// 	}

// 	goSrcFile := C.GoString(srcFile)
// 	goDestFile := C.GoString(destFile)
// 	goIsFile := isFile != 0

// 	err := cmd.BakupOptAdd(goSrcFile, goDestFile, goIsFile)
// 	if err != nil {
// 		return -2
// 	}

// 	return 0
// }

// C-exportable wrapper for BackupOptRm
//
//export BackupOptRmC
// func BackupOptRmC(file *C.char) C.int {
// 	if file == nil {
// 		return -1
// 	}

// 	goFile := C.GoString(file)
// 	err := cmd.BakupOptRm(goFile)
// 	if err != nil {
// 		return -2
// 	}

// 	return 0
// }

// C-exportable wrapper for GetAllOpt
//
//export GetAllOptC
func GetAllOptC(profilename *C.char) *C.FileOperationArray {
	g := cmd.NewGitCmd(C.GoString(profilename))
	operations, err := cmd.GetAllOpt(g.C)
	if err != nil {
		return nil
	}

	if len(operations) == 0 {
		return nil
	}

	// Allocate C array
	array := (*C.FileOperationArray)(C.malloc(C.size_t(unsafe.Sizeof(C.FileOperationArray{}))))
	if array == nil {
		return nil
	}
	count := len(operations)
	*array = C.FileOperationArray{}
	array.count = C.int(count)

	// Allocate memory for the array of FileOperationC structs
	operationsPtr := C.malloc(C.size_t(count) * C.size_t(unsafe.Sizeof(C.FileOperationC{})))
	if operationsPtr == nil {
		C.free(unsafe.Pointer(array))
		return nil
	}
	array.operations = (*C.FileOperationC)(operationsPtr)

	// Convert each FileOperation to a FileOperationC struct
	for i, op := range operations {
		cOp := (*C.FileOperationC)(unsafe.Pointer(uintptr(operationsPtr) + uintptr(i)*unsafe.Sizeof(C.FileOperationC{})))
		cOp.id = C.int64_t(op.ID)
		cOp.src_file = C.CString(op.SrcFile)
		cOp.dest_file = C.CString(op.DestFile)
		if op.IsFile {
			cOp.is_file = 1
		} else {
			cOp.is_file = 0
		}
		cOp.revcount = C.int(op.RevCount)
		// New field added for sub
		if op.Sub {
			cOp.sub = 1
		} else {
			cOp.sub = 0
		}
		cOp.tag = C.CString(op.Tag)
		cOp.add_time = C.CString(op.AddTime.Format("2006-01-02 15:04:05"))
		cOp.update_time = C.CString(op.UpdateTime.Format("2006-01-02 15:04:05"))

		// Check if any CString allocation failed
		if cOp.src_file == nil || cOp.dest_file == nil || cOp.tag == nil || cOp.add_time == nil || cOp.update_time == nil {
			// Free already allocated strings and memory
			for j := range i {
				prevOp := (*C.FileOperationC)(unsafe.Pointer(uintptr(operationsPtr) + uintptr(j)*unsafe.Sizeof(C.FileOperationC{})))
				if prevOp.src_file != nil {
					C.free(unsafe.Pointer(prevOp.src_file))
				}
				if prevOp.dest_file != nil {
					C.free(unsafe.Pointer(prevOp.dest_file))
				}
				if prevOp.tag != nil {
					C.free(unsafe.Pointer(prevOp.tag))
				}
				if prevOp.add_time != nil {
					C.free(unsafe.Pointer(prevOp.add_time))
				}
				if prevOp.update_time != nil {
					C.free(unsafe.Pointer(prevOp.update_time))
				}
			}
			C.free(unsafe.Pointer(operationsPtr))
			C.free(unsafe.Pointer(array))
			return nil
		}
	}

	return array
}

// C-exportable wrapper for AddFile
//
//export RmFileC
func RmFileC(profilename *C.char, filePath *C.char) C.int {
	if filePath == nil {
		return -1
	}
	goFilePath := C.GoString(filePath)
	g := cmd.NewGitCmd(C.GoString(profilename))
	err := g.RmFile(util.RepoPath(goFilePath))
	if err != nil {
		fmt.Printf("RmFileC failed %v err=%v", goFilePath, err)
		return -2
	}
	fmt.Printf("RmFileC success %v", goFilePath)
	return 0
}

// C-exportable wrapper for AddFile
//
//export AddFileC
func AddFileC(profilename *C.char, filePath *C.char) C.int {
	if filePath == nil {
		return 3
	}
	goFilePath := C.GoString(filePath)
	g := cmd.NewGitCmd(C.GoString(profilename))
	result := g.AddFile(goFilePath)
	if result.Err != nil {
		return 2
	}
	if result.Result == util.GitResultTypeNochange {
		return 1
	}
	return 0
}

// C-exportable wrapper for GitInitC
//
//export GitInitC
func GitInitC(profilename *C.char, filePath *C.char) C.int {
	if filePath == nil {
		return -1
	}
	goFilePath := C.GoString(filePath)
	goProfileName := C.GoString(profilename)
	if _, err := cmd.GitInitProfile(goProfileName, goFilePath); err == nil {
		return 0
	} else {
		return -1
	}
}

// C-exportable wrapper for GetFile
//
//export GetFileC
func GetFileC(profilename *C.char, filePath *C.char, commit *C.char, target *C.char) C.int {
	if filePath == nil || target == nil {
		return -1
	}
	goFilePath := C.GoString(filePath)
	goCommit := C.GoString(commit)
	goTarget := C.GoString(target)
	g := cmd.NewGitCmd(C.GoString(profilename))
	err := g.GetFile(util.RepoPath(goFilePath), goCommit, goTarget)
	if err != nil {
		return -2
	}
	return 0
}

// Helper function to free C strings
//
//export FreeString
func FreeString(str *C.char) {
	if str != nil {
		C.free(unsafe.Pointer(str))
	}
}

// Helper function to free GitChangeArray
//
//export FreeGitChangeArray
func FreeGitChangeArray(array *C.GitChangeArray) {
	if array == nil {
		return
	}

	// Free each GitChange struct's strings
	if array.changes != nil {
		for i := range int(array.count) {
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

// Helper function to free FileOperationArray
//
//export FreeFileOperationArray
func FreeFileOperationArray(array *C.FileOperationArray) {
	if array == nil {
		return
	}

	// Free each FileOperationC struct's strings
	if array.operations != nil {
		for i := range int(array.count) {
			op := (*C.FileOperationC)(unsafe.Pointer(uintptr(unsafe.Pointer(array.operations)) + uintptr(i)*unsafe.Sizeof(C.FileOperationC{})))
			if op.src_file != nil {
				C.free(unsafe.Pointer(op.src_file))
			}
			if op.dest_file != nil {
				C.free(unsafe.Pointer(op.dest_file))
			}
			if op.tag != nil {
				C.free(unsafe.Pointer(op.tag))
			}
			if op.add_time != nil {
				C.free(unsafe.Pointer(op.add_time))
			}
			if op.update_time != nil {
				C.free(unsafe.Pointer(op.update_time))
			}
		}
		// Free the array of FileOperationC structs
		C.free(unsafe.Pointer(array.operations))
	}

	// Free the FileOperationArray struct itself
	C.free(unsafe.Pointer(array))
}

// C-exportable wrapper for SetFileTag
//
//export SetFileTagC
func SetFileTagC(profilename *C.char, filePath *C.char, tag *C.char) C.int {
	if filePath == nil || tag == nil {
		return -1
	}
	goFilePath := C.GoString(filePath)
	goTag := C.GoString(tag)
	g := cmd.NewGitCmd(C.GoString(profilename))

	err := cmd.SetFileTag(util.RepoPath(goFilePath), goTag, g.C)
	if err != nil {
		fmt.Printf("SetFileTagC failed %v err=%v", goFilePath, err)
		return -2
	}
	fmt.Printf("SetFileTagC success %v", goFilePath)
	return 0
}

// C-exportable wrapper for GetFileTag
//
//export GetFileTagC
func GetFileTagC(profilename *C.char, filePath *C.char) *C.char {
	if filePath == nil {
		return nil
	}
	goFilePath := C.GoString(filePath)
	g := cmd.NewGitCmd(C.GoString(profilename))

	tag, err := cmd.GetFileTag(util.RepoPath(goFilePath), g.C)
	if err != nil {
		fmt.Printf("GetFileTagC failed %v err=%v", goFilePath, err)
		return nil
	}
	fmt.Printf("GetFileTagC success %v tag=%v", goFilePath, tag)
	return C.CString(tag)
}

func main() {
	// Main function is required for c-shared buildmode but won't be called when used as a library
}
