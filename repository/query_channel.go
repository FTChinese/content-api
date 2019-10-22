package repository

import "fmt"

const (
	stmtFrontPage = `
    SELECT story.id AS id,
        story.tag AS tag,
        story.cheadline AS title,
        story.clongleadbody AS standfirst,
        picture.piclink  AS cover_url,
        FROM_UNIXTIME(story.fileupdatetime) AS created_utc,
		FROM_UNIXTIME(story.last_publish_time) AS updated_utc
    FROM cmstmp01.story
        LEFT JOIN (
			cmstmp01.story_pic AS storyToPic
			INNER JOIN cmstmp01.picture AS picture
		)
		ON story.id = storyToPic.storyid 
		AND picture.id = storyToPic.picture_id
    WHERE story.publish_status = 'publish'
        %s
    ORDER BY story.priority,
        story.fileupdatetime`
)

// Current front page.
var stmtFrontPageToday = fmt.Sprintf(stmtFrontPage, `AND story.pubdate = (
    SELECT pubdate
    FROM cmstmp01.story
    WHERE publish_status = 'publish'
    ORDER BY pubdate DESC
    LIMIT 1
)`)

// The front page on a certain date.
var stmtFrontPageArchive = fmt.Sprintf(stmtFrontPage, `AND FROM_UNIXTIME(story.pubdate, "%Y-%m-%d") = ?`)
