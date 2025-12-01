#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/stat.h>
#include "build/libgitcmd.h"

void create_test_file(const char* filename, const char* content) {
    FILE* file = fopen(filename, "w");
    if (file) {
        fprintf(file, "%s", content);
        fclose(file);
        printf("Created test file: %s\n", filename);
    } else {
        printf("Failed to create test file: %s\n", filename);
    }
}

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
    const char* test_file = "/tmp/test_anybakup.txt";
    const char* output_file = "/tmp/test_output.txt";

    printf("Testing libgitcmd.so with real files\n");
    printf("=====================================\n");

    // Create a test file
    create_test_file(test_file, "Hello, this is a test file for anybakup!\n");

    // Test 1: Try to get file log for our test file
    printf("\n=== Test 1: GetFileLogC (real file) ===\n");
    char* log_result = GetFileLogC(test_file);
    test_function("GetFileLogC", log_result);

    // Test 2: Try to add the test file
    printf("\n=== Test 2: AddFileC (real file) ===\n");
    char* add_result = AddFileC(test_file);
    test_function("AddFileC", add_result);

    // Test 3: Try to get the file (this might fail if not in repo)
    printf("\n=== Test 3: GetFileC (real file) ===\n");
    char* get_result = GetFileC(test_file, "HEAD", output_file);
    test_function("GetFileC", get_result);

    // Check if output file was created
    if (access(output_file, F_OK) == 0) {
        printf("Output file %s was created successfully\n", output_file);

        // Read and display the output file content
        FILE* file = fopen(output_file, "r");
        if (file) {
            printf("File content:\n");
            char buffer[256];
            while (fgets(buffer, sizeof(buffer), file)) {
                printf("%s", buffer);
            }
            fclose(file);
        }
        // Clean up output file
        unlink(output_file);
    } else {
        printf("Output file %s was not created (expected if file not in repo)\n", output_file);
    }

    // Clean up test file
    unlink(test_file);
    printf("\nCleaned up test files\n");

    printf("\n=== All tests completed ===\n");
    return 0;
}