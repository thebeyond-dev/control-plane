ALTER TABLE "user"
    ADD COLUMN referrer_id UUID REFERENCES "user"(id) ON DELETE SET NULL,
    ADD COLUMN referral_balance INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN referral_commission_rate INTEGER NOT NULL DEFAULT 33,
    ADD COLUMN referrals_count INTEGER NOT NULL DEFAULT 0;

CREATE INDEX idx_user_referrer_id ON "user"(referrer_id);