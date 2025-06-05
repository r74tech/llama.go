enable_testing()

set(SRCS test_runner.cpp)
set(TARGET test_runner)

add_executable(${TARGET} ${SRCS})
target_link_libraries(${TARGET} PRIVATE common llama ${CMAKE_THREAD_LIBS_INIT} llama_core)

add_test(NAME RunnerTest COMMAND ${TARGET})