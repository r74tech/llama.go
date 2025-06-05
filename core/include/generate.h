#pragma once

#ifdef __cplusplus
extern "C" {
#endif

const char * llama_generate(const char * model_file,const char * input_prompt,int n_gpu_layers,int n);

#ifdef __cplusplus
}
#endif