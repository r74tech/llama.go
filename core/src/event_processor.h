#pragma once
#include <queue>
#include <mutex>
#include <future>
#include "message.h"

class EventProcessor {
public:
    struct Event {
        std::vector<Message> data;
        std::promise<std::string> result;
    };

    std::string enqueue(const std::vector<Message>& data);

    bool dequeue(Event& event);

    void stop();

private:
    std::queue<Event> m_queue;
    std::mutex m_mtx;
    std::condition_variable m_cv;
    bool m_stop = false;
};