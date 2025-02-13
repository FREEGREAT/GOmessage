-- Створення ключового простору (якщо ще не створено)
CREATE KEYSPACE IF NOT EXISTS chat_app WITH REPLICATION = {
    'class': 'SimpleStrategy',
    'replication_factor': 1
};

USE chat_app;

-- Таблиця для зберігання повідомлень
CREATE TABLE IF NOT EXISTS messages (
    message_id UUID PRIMARY KEY,
    user_id1 UUID,
    user_id2 UUID,
    message TEXT,
    sent_time TIMESTAMP,
    is_edited BOOLEAN
);

-- Таблиця для зберігання чатів
CREATE TABLE IF NOT EXISTS chats (
    chat_id UUID PRIMARY KEY,
    user_id1 TEXT,
    user_id2 TEXT,
    action TEXT
);
