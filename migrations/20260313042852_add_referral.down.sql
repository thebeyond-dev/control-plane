DROP INDEX IF EXISTS idx_user_referrer_id;

ALTER TABLE "user" 
    DROP COLUMN IF EXISTS referrer_id,
    DROP COLUMN IF EXISTS referral_balance,
    DROP COLUMN IF EXISTS referral_commission_rate,
    DROP COLUMN IF EXISTS referrals_count;