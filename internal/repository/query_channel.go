package repository

//const stmtNav = `
//    SELECT top.id AS id,
//      top.code AS name,
//      top.name AS title,
//      GROUP_CONCAT(
//        CONCAT_WS(':', sub.code, sub.name) ORDER BY sub.id ASC
//        SEPARATOR ';'
//      ) AS children
//    FROM cmstmp01.channel AS top
//      INNER JOIN cmstmp01.channel AS sub
//      ON top.id = sub.reid
//    WHERE top.reid = 0
//    GROUP BY top.id
//	ORDER BY top.priority ASC`

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
