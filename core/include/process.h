#pragma once

#ifdef __cplusplus
extern "C" {
#endif

    bool llama_start(const char * args);
    bool llama_stop();
    const char * llama_gen(const char * prompt);

#ifdef __cplusplus
}
#endif