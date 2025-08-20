#include "process.h"
#include "common.h"
#include "llama.h"
#include "log.h"
#include "runner.h"
#include <iostream>
#include <sstream>
#include <string>
#include <vector>

// Forward declaration
int process(common_params &params);

static Runner *g_runner;
static int g_idx = 0;

// Global variables for memory-loaded model (NOT static so they can be accessed
// from runner.cpp)
const void *g_model_buffer = nullptr;
size_t g_model_buffer_size = 0;
bool g_use_mmap = false;

extern "C" {
int llama_start(const char *args, int async, const char *prompt) {
    if (g_runner != nullptr) {
        LOG("Delete last runner: id=%d\n", g_runner->getID());
        delete g_runner;
        g_runner = nullptr;
    }
    std::istringstream iss(args);
    std::vector<std::string> v_args;
    std::string v_a;
    while (iss >> v_a) {
        v_args.push_back(v_a);
    }

    g_runner = new Runner(g_idx, v_args, async > 0, std::string(prompt));
    g_idx++;
    if (g_runner->start()) {
        return EXIT_SUCCESS;
    }
    return EXIT_FAILURE;
}

int llama_stop() {
    if (g_runner == nullptr) {
        LOG("Runner is already delete\n");
        return EXIT_SUCCESS;
    }
    bool ret = g_runner->stop();
    LOG("Delete last runner: id=%d\n", g_runner->getID());
    delete g_runner;
    g_runner = nullptr;
    if (ret) {
        return EXIT_SUCCESS;
    }
    return EXIT_FAILURE;
}

const char *llama_gen(const char *prompt) {
    if (g_runner == nullptr) {
        LOG_ERR("Not init llama\n");
        return "";
    }
    std::string result = g_runner->generate(std::string(prompt));
    char *arr = new char[result.size() + 1];
    std::copy(result.begin(), result.end(), arr);
    arr[result.size()] = '\0';

    return arr;
}

const char *llama_chat(const char **roles, const char **contents, int size) {
    if (g_runner == nullptr) {
        LOG_ERR("Not init llama\n");
        return "";
    }
    std::vector<Message> msgs;

    for (int i = 0; i < size; i++) {
        Message msg;
        msg.role = roles[i];
        msg.content = contents[i];

        msgs.push_back(msg);
    }

    std::string result = g_runner->chat(msgs);
    char *arr = new char[result.size() + 1];
    std::copy(result.begin(), result.end(), arr);
    arr[result.size()] = '\0';

    return arr;
}
} // extern "C"

// Common function to run model from memory
static int llama_run_from_memory_internal(const void *buffer, size_t size,
                                          bool is_mmap, const char *args,
                                          int async, const char *prompt) {
    if (g_runner != nullptr) {
        LOG("Delete last runner: id=%d\n", g_runner->getID());
        delete g_runner;
        g_runner = nullptr;
    }

    // Store the memory buffer in global variables for Runner to access
    g_model_buffer = buffer;
    g_model_buffer_size = size;
    g_use_mmap = is_mmap;

    // Parse arguments into vector of strings for Runner
    std::vector<std::string> v_args;
    v_args.push_back("llama"); // dummy executable name

    // Parse user-provided arguments (but skip model path)
    std::istringstream iss(args);
    std::string arg;
    while (iss >> arg) {
        // Skip model path argument since we're loading from memory
        if (arg != "-m" && arg != "--model") {
            v_args.push_back(arg);
        } else {
            // Skip the next argument (model path)
            iss >> arg;
        }
    }

    // Create runner with the modified arguments
    std::string prompt_str = prompt ? std::string(prompt) : "";
    g_runner = new Runner(g_idx, v_args, async > 0, prompt_str);
    g_idx++;

    if (g_runner->start()) {
        return EXIT_SUCCESS;
    }

    return EXIT_FAILURE;
}

extern "C" int llama_start_from_memory(const void *model_data, size_t size,
                                       const char *args, int async,
                                       const char *prompt) {
    LOG("Starting llama from memory buffer (size=%zu bytes)\n", size);
    return llama_run_from_memory_internal(model_data, size, false, args, async,
                                          prompt);
}

extern "C" int llama_start_from_mmap(const void *addr, size_t size,
                                     const char *args, int async,
                                     const char *prompt) {
    LOG("Starting llama from mmap'd memory (addr=%p, size=%zu)\n", addr, size);
    return llama_run_from_memory_internal(addr, size, true, args, async,
                                          prompt);
}
