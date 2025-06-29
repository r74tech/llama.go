#include "event_processor.h"

std::string EventProcessor::enqueue(const std::vector<Message>& data) {
    Event event;
    event.data = data;

    std::future<std::string> resultFuture = event.result.get_future();

    {
        std::lock_guard<std::mutex> lock(m_mtx);
        m_queue.push(std::move(event));
    }

    m_cv.notify_one();
    return resultFuture.get();
}

bool EventProcessor::dequeue(Event& event) {
    std::unique_lock<std::mutex> lock(m_mtx);
    m_cv.wait(lock, [this]() {
        return !m_queue.empty() || m_stop;
    });

    if (m_stop && m_queue.empty())
        return false;

    event = std::move(m_queue.front());
    m_queue.pop();
    return true;
}

void EventProcessor::stop() {
    {
        std::lock_guard<std::mutex> lock(m_mtx);
        m_stop = true;
    }
    m_cv.notify_all();
}