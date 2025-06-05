#include <iostream>
#include <thread>
#include <chrono>

#include "process.h"

int main() {
    std::string args="llama-cli -m ./qwen2.5-0.5b-q8_0.gguf -no-cnv --seed 0";
    bool ret=llama_start(args.c_str());
    if (!ret) {
        return 1;
    }
    std::string prompt="why the sky is blue";
    ret=llama_chat(prompt.c_str());
    if (!ret) {
        return 1;
    }
    int seconds = 2;
    std::cout << "sleep...:"<<seconds<<" seconds" << std::endl;
    std::this_thread::sleep_for(std::chrono::seconds(seconds));

    ret=llama_stop();
    if (!ret) {
        return 1;
    }
    return 0;
}