#pragma once

#ifdef __cplusplus
extern "C" {
#endif

int llama_interactive(const char * model_file,int n_gpu_layers,int ctx_size);

#ifdef __cplusplus
}
#endif