package repository

import "fmt"

const (
	stmtStoryTeaser = `
    SELECT story.id AS id,
        story.tag AS tag,
        story.cheadline AS title,
        story.clongleadbody AS standfirst,
        picture.piclink  AS cover_url,
        FROM_UNIXTIME(story.fileupdatetime) AS created_utc,
		FROM_UNIXTIME(story.last_publish_time) AS updated_utc`

	stmtChannelContent = stmtStoryTeaser + `
    FROM cmstmp01.channel_detail AS ch_story
        LEFT JOIN cmstmp01.story 
            ON ch_story.id = story.id
        LEFT JOIN cmstmp01.story_pic
            ON story.id = story_pic.storyid
        LEFT JOIN cmstmp01.picture AS picture
            ON story_pic.picture_id = picture.id
    WHERE ch_story.chaid = ?
        AND story.publish_status = 'publish'
    ORDER BY ch_story.addtime DESC
	LIMIT ? OFFSET ?`

	stmtFrontPage = stmtStoryTeaser + `
    FROM cmstmp01.story
        LEFT JOIN cmstmp01.story_pic
            ON story.id = story_pic.storyid 
        LEFT JOIN cmstmp01.picture AS picture
            ON picture.id = story_pic.picture_id
    WHERE story.publish_status = 'publish'
        %s
    ORDER BY story.priority,
        story.fileupdatetime`
)

// Current front page.
var stmtFrontPageToday = fmt.Sprintf(stmtFrontPage, `
AND story.pubdate = (
    SELECT pubdate
    FROM cmstmp01.story
    WHERE publish_status = 'publish'
    ORDER BY pubdate DESC
    LIMIT 1
)`)

// The front page on a certain date.
var stmtFrontPageArchive = fmt.Sprintf(stmtFrontPage, `
AND FROM_UNIXTIME(story.pubdate, "%Y-%m-%d") = ?`)

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
WHERE full_path IS NOT NULL
    AND is_active = 1
ORDER BY id`

const stmtChannelSetting = stmtChannelItem + `
WHERE full_path = ?
    AND is_active = 1
LIMIT 1`
