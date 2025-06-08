#include <iostream>
#include <thread>
#include <chrono>
#include <cstdlib>
#include <future>
#include <unistd.h>

#include "process.h"

int main() {
    const char* env_model = "LLAMA_TEST_MODEL";
    const char* model = std::getenv(env_model);

    if (model == nullptr) {
        std::cerr << "error：can't find " << env_model << std::endl;
        return EXIT_FAILURE;
    }

    std::cout << "env: " << env_model << "=" << model << std::endl;

    std::stringstream ss;
    ss << "test_runner_gen -m " << model << " -i --seed 0";

    std::future<void> ll_main = std::async(std::launch::async, [&ss](){
        bool ret = llama_start(ss.str().c_str(), true);
        std::cout<<"Result0:"<<ret<<std::endl;
    });

    std::this_thread::sleep_for(std::chrono::seconds(2));

    std::future<void> ll_gen = std::async(std::launch::async, [](){
        std::string prompt="why sky is blue";
        std::string content = llama_gen(prompt.c_str());
        if (content.empty()) {
            return;
        }

        std::cout<<"Response:"<<content<<std::endl;

        bool ret =llama_stop();
        std::cout<<"Result1:"<<ret<<std::endl;
        });


    ll_gen.wait();
    ll_main.wait();


    std::cout<<"success"<<std::endl;

    return EXIT_SUCCESS;
}