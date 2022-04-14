CREATE TABLE message
(
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_datetime timestamp,
    message_text     TEXT,
    message_chat_id  uuid,
    message_user_id  integer
);