#include "safe_queue.h"

void SafeQueue::enqueue(std::string task) {
    std::lock_guard<std::mutex> lock(m_mtx);
    m_queue.push(std::move(task));
    m_cv.notify_one();
}

bool SafeQueue::dequeue(std::string& task) {
    std::unique_lock<std::mutex> lock(m_mtx);
    m_cv.wait(lock, [this]() {
        return !m_queue.empty() || m_stop;
    });

    if (m_stop && m_queue.empty())
        return false;

    task = std::move(m_queue.front());
    m_queue.pop();
    return true;
}

void SafeQueue::stop() {
    {
        std::lock_guard<std::mutex> lock(m_mtx);
        m_stop = true;
    }
    m_cv.notify_all();
}