CREATE TABLE user_chat
(
    chat_id uuid REFERENCES chats,
    user_id integer REFERENCES users,
    PRIMARY KEY (chat_id, user_id)
);