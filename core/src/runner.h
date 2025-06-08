#pragma once

#include "safe_queue.h"
#include "sampling.h"

class Runner {
private:
    int m_id;
    const std::vector<std::string> m_args;
    SafeQueue m_queue;
    std::atomic<bool> m_running;
    bool m_async;

    llama_context           * m_ctx;
    llama_model             * m_model;
    common_sampler          * m_smpl;
    common_params           * m_params;

    std::vector<llama_token> * m_input_tokens;
    std::ostringstream       * m_output_ss;
    std::vector<llama_token> * m_output_tokens;

public:
    Runner(int id,const std::vector<std::string>& args);
    ~Runner();
    bool start();
    bool stop();
    const std::string generate(const std::string& prompt);
    int getID();
    bool isRunning();

    bool getPrompt(std::string& prompt);
};