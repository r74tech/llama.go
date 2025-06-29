#pragma once

#ifdef __cplusplus
extern "C" {
#endif

    int llama_start(const char * args,int async,const char * prompt);
    int llama_stop();
    const char * llama_gen(const char * prompt);
    const char * llama_chat(const char **roles,const char **contents, int size);

#ifdef __cplusplus
}
#endif