cmake_minimum_required(VERSION 3.14)
project(llama_core)

set(CMAKE_C_STANDARD 11)
set(CMAKE_C_STANDARD_REQUIRED true)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED true)

set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)
set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY_DEBUG ${CMAKE_BINARY_DIR}/lib)
set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY_RELEASE ${CMAKE_BINARY_DIR}/lib)
set(EXECUTABLE_OUTPUT_PATH ${CMAKE_BINARY_DIR}/bin)

# Apply memory loading patch if not already applied
if(NOT EXISTS "${CMAKE_CURRENT_SOURCE_DIR}/llama.cpp/.memory_patch_applied")
    message(STATUS "Applying memory loading patch...")

    # Check if we're on Windows
    if(WIN32)
        # Try git apply first (more commonly available on Windows)
        execute_process(
            COMMAND git apply --check ${CMAKE_CURRENT_SOURCE_DIR}/patches/memory-loading.patch
            WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/llama.cpp
            RESULT_VARIABLE PATCH_CHECK_RESULT
            OUTPUT_QUIET
            ERROR_QUIET
        )

        if(PATCH_CHECK_RESULT EQUAL 0)
            # Patch can be applied
            execute_process(
                COMMAND git apply ${CMAKE_CURRENT_SOURCE_DIR}/patches/memory-loading.patch
                WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/llama.cpp
                RESULT_VARIABLE PATCH_RESULT
                OUTPUT_VARIABLE PATCH_OUTPUT
                ERROR_VARIABLE PATCH_ERROR
            )
        else()
            # Try patch command if available
            find_program(PATCH_EXECUTABLE patch)
            if(PATCH_EXECUTABLE)
                execute_process(
                    COMMAND ${PATCH_EXECUTABLE} -p1 -N -i ${CMAKE_CURRENT_SOURCE_DIR}/patches/memory-loading.patch
                    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/llama.cpp
                    RESULT_VARIABLE PATCH_RESULT
                    OUTPUT_VARIABLE PATCH_OUTPUT
                    ERROR_VARIABLE PATCH_ERROR
                )
            else()
                message(STATUS "Neither git apply nor patch command available - assuming patch already applied")
                set(PATCH_RESULT 1)
            endif()
        endif()
    else()
        # Unix/Linux/macOS - use patch command
        execute_process(
            COMMAND patch -p1 -N
            INPUT_FILE ${CMAKE_CURRENT_SOURCE_DIR}/patches/memory-loading.patch
            WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/llama.cpp
            RESULT_VARIABLE PATCH_RESULT
            OUTPUT_VARIABLE PATCH_OUTPUT
            ERROR_VARIABLE PATCH_ERROR
        )
    endif()

    if(PATCH_RESULT EQUAL 0)
        file(WRITE "${CMAKE_CURRENT_SOURCE_DIR}/llama.cpp/.memory_patch_applied" "patch applied")
        message(STATUS "Memory loading patch applied successfully")
    else()
        message(STATUS "Patch already applied or failed (this is okay if already patched)")
    endif()
endif()

# llama.cpp
set(GGML_USE_OPENMP OFF)
set(GGML_OPENMP OFF)
set(GGML_STATIC ON)
set(BUILD_SHARED_LIBS OFF)
set(LLAMA_BUILD_COMMON ON)
set(LLAMA_CURL OFF)
add_subdirectory(llama.cpp)

# core
set(SRCS src/generate.cpp src/interactive.cpp src/process.cpp src/runner.cpp src/event_processor.cpp src/embedding.cpp)
set(TARGET llama_core)

include_directories(./include)
include_directories(./llama.cpp/include)
include_directories(./llama.cpp/common)

link_directories(${CMAKE_BINARY_DIR}/lib)

add_library(${TARGET} STATIC ${SRCS})
target_link_libraries(${TARGET} PRIVATE common llama ${CMAKE_THREAD_LIBS_INIT})

# test
option(BUILD_TEST "Build the testing tree." OFF)

if(BUILD_TEST)
    include(CTest)
    add_subdirectory(tests)
endif()