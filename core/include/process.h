#ifndef PROCESS_H
#define PROCESS_H

#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif

// Original functions
int llama_start(const char *args, int async, const char *prompt);
int llama_stop();
const char *llama_gen(const char *prompt);
const char *llama_chat(const char **roles, const char **contents, int size);

// Memory-based loading functions
int llama_start_from_memory(const void *model_data, size_t size,
                            const char *args, int async, const char *prompt);
int llama_start_from_mmap(const void *addr, size_t size, const char *args,
                          int async, const char *prompt);

#ifdef __cplusplus
}
#endif

#endif // PROCESS_H
