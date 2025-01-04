-- Створення ключового простору (якщо ще не створено)
CREATE KEYSPACE IF NOT EXISTS chat_app WITH REPLICATION = {
    'class': 'SimpleStrategy',
    'replication_factor': 1
};

USE chat_app;

-- Таблиця для зберігання чатів користувачів
CREATE TABLE IF NOT EXISTS user_chats (
    user_id TEXT,
    chat_id TEXT,
    chat_name TEXT,
    PRIMARY KEY (user_id, chat_id)
);

-- Таблиця для зберігання повідомлень
CREATE TABLE IF NOT EXISTS messages (
    message_id TEXT PRIMARY KEY,
    user_id1 TEXT,
    user_id2 TEXT,
    message TEXT,
    sent_time TIMESTAMP,
    is_edited BOOLEAN
);

-- Таблиця для зберігання статусу повідомлень
CREATE TABLE IF NOT EXISTS message_status (
    chat_id TEXT,
    message_id TEXT,
    user_id TEXT,
    status TEXT,
    PRIMARY KEY (chat_id, message_id, user_id)
);

-- Таблиця для зберігання чатів
CREATE TABLE IF NOT EXISTS chats (
    chat_id TEXT PRIMARY KEY,
    user_id1 TEXT,
    user_id2 TEXT,
    action TEXT
);

-- Таблиця для зберігання повідомлень чату
CREATE TABLE IF NOT EXISTS chat_messages (
    chat_id TEXT,
    message_id TEXT,
    sender_id TEXT,
    message_text TEXT,
    sent_at TIMESTAMP,
    PRIMARY KEY (chat_id, message_id)
);