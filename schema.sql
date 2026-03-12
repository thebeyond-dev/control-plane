CREATE TABLE node(
    id TEXT PRIMARY KEY,
    address TEXT NOT NULL,
    load_percent INTEGER NOT NULL DEFAULT 0
        CHECK (load_percent BETWEEN 0 AND 100)
);

CREATE TABLE "user"(
    id UUID PRIMARY KEY,
    devices INTEGER NOT NULL,
    bandwidth INTEGER NOT NULL,
    subscription_expires_at TIMESTAMPTZ,
    language_code TEXT NOT NULL,
    currency_code TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE identity(
    provider TEXT,
    provider_id TEXT,
    user_id UUID REFERENCES "user",
    PRIMARY KEY (provider, provider_id)
);

CREATE TABLE device(
	user_id UUID NOT NULL REFERENCES "user" ON DELETE CASCADE,
	pubkey TEXT NOT NULL,
	node_id TEXT REFERENCES node ON DELETE SET NULL,
	name TEXT NOT NULL,
    PRIMARY KEY(user_id, pubkey)
);