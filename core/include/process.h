#pragma once

#ifdef __cplusplus
extern "C" {
#endif

    int llama_start(const char * args,int async,const char * prompt);
    int llama_stop();
    const char * llama_gen(const char * prompt);

#ifdef __cplusplus
}
#endif