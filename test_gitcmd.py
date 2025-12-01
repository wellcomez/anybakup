#!/usr/bin/env python3
import ctypes
import os
import sys

# Load the shared library
lib_path = os.path.join(os.path.dirname(__file__), "build", "libgitcmd.so")
try:
    lib = ctypes.CDLL(lib_path)
except OSError as e:
    print(f"Error loading library: {e}")
    sys.exit(1)

# Define function prototypes
lib.GetFileLogC.argtypes = [ctypes.c_char_p]
lib.GetFileLogC.restype = ctypes.c_char_p

lib.AddFileC.argtypes = [ctypes.c_char_p]
lib.AddFileC.restype = ctypes.c_char_p

lib.FreeString.argtypes = [ctypes.c_char_p]
lib.FreeString.restype = None

def main():
    print("Testing libgitcmd.so from Python")
    print("=================================")

    # Test 1: GetFileLogC with non-existent file
    print("\n=== Test 1: GetFileLogC (non-existent file) ===")
    log_result = lib.GetFileLogC(b"/non/existent/file.txt")
    if log_result:
        print(f"Result: {log_result.decode('utf-8')}")
        lib.FreeString(log_result)
    else:
        print("Error: Function returned NULL")

    # Test 2: AddFileC with non-existent file
    print("\n=== Test 2: AddFileC (non-existent file) ===")
    add_result = lib.AddFileC(b"/non/existent/file.txt")
    if add_result:
        print(f"Result: {add_result.decode('utf-8')}")
        lib.FreeString(add_result)
    else:
        print("Error: Function returned NULL")

    # Test 3: AddFileC with NULL input
    print("\n=== Test 3: AddFileC (NULL input) ===")
    add_result = lib.AddFileC(None)
    if add_result:
        print(f"Result: {add_result.decode('utf-8')}")
        lib.FreeString(add_result)
    else:
        print("Error: Function returned NULL")

    print("\n=== All Python tests completed ===")

if __name__ == "__main__":
    main()