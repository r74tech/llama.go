cmake_minimum_required(VERSION 3.14)
project(llama_core)

set(CMAKE_C_STANDARD 11)
set(CMAKE_C_STANDARD_REQUIRED true)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED true)

set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)
set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY_DEBUG ${CMAKE_BINARY_DIR}/lib)
set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY_RELEASE ${CMAKE_BINARY_DIR}/lib)
set(CMAKE_DEBUG_POSTFIX "_d")

# llama.cpp
set(GGML_USE_OPENMP OFF)
set(GGML_OPENMP OFF)
set(GGML_STATIC ON)
set(BUILD_SHARED_LIBS OFF)
add_subdirectory(llama.cpp)

# core
set(SRCS core.cpp)
set(TARGET llama_core)

include_directories(./)

link_directories(./../build/bin)

add_library(${TARGET} STATIC ${SRCS})
target_link_libraries(${TARGET} PRIVATE llama ${CMAKE_THREAD_LIBS_INIT})
