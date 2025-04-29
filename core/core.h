#pragma once

#ifdef __cplusplus
extern "C" {
#endif

int llama_app(const char * model_file,const char * input_prompt,int n_gpu_layers,int n);

#ifdef __cplusplus
}
#endif