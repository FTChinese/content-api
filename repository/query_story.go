package repository

const stmtStoryFrom = `
FROM cmstmp01.story AS story
    LEFT JOIN cmstmp01.story_pic AS storyToPic
        ON story.id = storyToPic.storyid 
    LEFT JOIN cmstmp01.picture AS picture
        ON picture.id = storyToPic.picture_id
`

const stmtStoryBase = `
SELECT story.id AS id,
    FROM_UNIXTIME(story.fileupdatetime) AS created_utc,
    FROM_UNIXTIME(story.last_publish_time) AS updated_utc,
    story.tag AS tag,
    story.cheadline AS title_cn,
    story.clongleadbody AS long_lead_cn,
    picture.piclink  AS cover_url
`

const stmtStoryTeaser = stmtStoryBase + `
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

const stmtFrontPageToday = stmtStoryBase + stmtStoryFrom + `
WHERE story.publish_status = 'publish'
    AND story.pubdate = (
        SELECT pubdate
        FROM cmstmp01.story
        WHERE publish_status = 'publish'
        ORDER BY pubdate DESC
        LIMIT 1
    )
ORDER BY story.priority,
    story.fileupdatetime`

// The front page on a certain date.
const stmtFrontPageArchive = stmtStoryBase + stmtStoryFrom + `
WHERE story.publish_status = 'publish'
    AND FROM_UNIXTIME(story.pubdate, "%Y-%m-%d") = ?
ORDER BY story.priority,
    story.fileupdatetime`

const stmtStory = stmtStoryBase + `,
    story.cheadline AS title_cn,
    story.eheadline AS title_en,
    story.cbyline_description AS byline_desc_cn,
    story.ebyline_description AS byline_desc_en,
    story.cauthor AS byline_author_cn,
    story.eauthor AS byline_author_en,
    story.cbyline_status AS byline_status_cn,
    story.ebyline_status AS byline_status_en,
    story.accessright AS access_right,
    story.area AS area,
    story.genre AS genre,
    story.industry AS industry,
    story.topic AS topic,
    story.cbody AS body_cn,
    story.ebody AS body_en
FROM cmstmp01.story AS story
    LEFT JOIN cmstmp01.story_pic AS storyToPic
        ON story.id = storyToPic.storyid 
    LEFT JOIN cmstmp01.picture AS picture
        ON picture.id = storyToPic.picture_id
WHERE story.id = ?
    AND story.publish_status = 'publish'`
