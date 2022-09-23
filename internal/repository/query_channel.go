package repository

const stmtChannelItem = `
SELECT id,
    parent_id,
    key_name,
    title,
    description,
    created_utc,
    updated_utc
FROM cmstmp01.channel`

const stmtListChannels = stmtChannelItem + `
WHERE key_name IS NOT NULL
    AND is_active = 1
ORDER BY id`

const stmtChannelSetting = stmtChannelItem + `
WHERE key_name = ?
    AND is_active = 1
LIMIT 1`
