CREATE TABLE chats
(
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_avatar    VARCHAR(120),
    chat_password  CHAR(64),
    user_chat_user integer
);


