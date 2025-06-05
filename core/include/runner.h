#pragma once

#include "sampling.h"

class Runner {
private:
    int id;
    const std::vector<std::string> args;

    llama_context           ** g_ctx;
    llama_model             ** g_model;
    common_sampler          ** g_smpl;
    common_params            * g_params;
    std::vector<llama_token> * g_input_tokens;
    std::ostringstream       * g_output_ss;
    std::vector<llama_token> * g_output_tokens;
    bool is_interacting  = false;
    bool need_insert_eot = false;

public:
    Runner(int id,const std::vector<std::string> args);
    ~Runner();
    bool Start();
    bool Stop();
    const std::string Chat(const std::string input_prompt);
    int GetID();
};