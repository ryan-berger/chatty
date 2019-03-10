CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE conversation
(
  id     UUID PRIMARY KEY,
  name   TEXT             DEFAULT NULL,
  direct BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE conversant
(
  id          UUID PRIMARY KEY,
  external_id UUID UNIQUE NOT NULL
);

CREATE TABLE conversant_conversation
(
  conversation_id UUID REFERENCES conversation NOT NULL,
  conversant_id   UUID REFERENCES conversant   NOT NULL,
  PRIMARY KEY (conversation_id, conversant_id)
);

CREATE TABLE chat_message
(
  id           UUID PRIMARY KEY,
  message      TEXT                         NOT NULL,
  sender       UUID REFERENCES conversant   NOT NULL,
  conversation UUID REFERENCES conversation NOT NULL
);