#pragma once

#include "chat.h"

struct Message {
    std::string role;
    std::string content;

    void fillMessage(common_chat_msg& msg) const {
        msg.content=content;
        msg.role=role;
    }
};