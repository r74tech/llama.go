#include <iostream>
#include <thread>
#include <chrono>
#include <cstdlib>

#include "embedding.h"

int main() {
    const char* env_model = "LLAMA_TEST_MODEL";
    const char* model = std::getenv(env_model);

    if (model == nullptr) {
        std::cerr << "error：can't find " << env_model << std::endl;
        return EXIT_FAILURE;
    }

    std::cout << "env: " << env_model << "=" << model << std::endl;

    std::stringstream ss;
    ss << "test_embedding -m " << model << " --pooling mean";

    int ret=llama_embedding(ss.str().c_str(),std::string("Hello World").c_str());
    if (ret != 0) {
        return EXIT_FAILURE;
    }
    return EXIT_SUCCESS;
}