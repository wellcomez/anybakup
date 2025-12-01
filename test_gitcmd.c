#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "build/libgitcmd.h"

void test_function(const char* func_name, char* result) {
    printf("\n=== Testing %s ===\n", func_name);
    if (result) {
        printf("Result: %s\n", result);
        FreeString(result);
    } else {
        printf("Error: Function returned NULL\n");
    }
}

int main() {
    printf("Testing libgitcmd.so dynamic library\n");
    printf("=====================================\n");

    // Test 1: GetFileLogC with a non-existent file
    printf("\n=== Test 1: GetFileLogC (non-existent file) ===\n");
    char* log_result = GetFileLogC("/non/existent/file.txt");
    test_function("GetFileLogC", log_result);

    // Test 2: AddFileC with a non-existent file
    printf("\n=== Test 2: AddFileC (non-existent file) ===\n");
    char* add_result = AddFileC("/non/existent/file.txt");
    test_function("AddFileC", add_result);

    // Test 3: GetFileC with invalid parameters
    printf("\n=== Test 3: GetFileC (NULL target) ===\n");
    char* get_result = GetFileC("somefile.txt", "HEAD", NULL);
    test_function("GetFileC", get_result);

    // Test 4: GetFileC with NULL file path
    printf("\n=== Test 4: GetFileC (NULL file path) ===\n");
    char* get_result2 = GetFileC(NULL, "HEAD", "/tmp/output.txt");
    test_function("GetFileC", get_result2);

    // Test 5: AddFileC with NULL input
    printf("\n=== Test 5: AddFileC (NULL input) ===\n");
    char* add_result2 = AddFileC(NULL);
    test_function("AddFileC", add_result2);

    // Test 6: GetFileLogC with NULL input
    printf("\n=== Test 6: GetFileLogC (NULL input) ===\n");
    char* log_result2 = GetFileLogC(NULL);
    test_function("GetFileLogC", log_result2);

    printf("\n=== All tests completed ===\n");
    return 0;
}