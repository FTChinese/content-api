package repository

import "fmt"

const stmtInteractiveSource = `
FROM cmstmp01.interactive_story AS story
    LEFT JOIN cmstmp01.interactive_pic AS a_p
        ON story.id = a_p.interactive_id
    LEFT JOIN cmstmp01.picture AS picture
        ON picture.id = a_p.picture_id
`

const stmtAudioBase = `
SELECT story.id,
    FROM_UNIXTIME(story.fileupdatetime) AS created_utc,
    FROM_UNIXTIME(story.last_publish_time) AS updated_utc,
    story.tag,
    story.cheadline AS title,
    story.clongleadbody AS standfirst,
    picture.piclink AS cover_url
    %s
` + stmtInteractiveSource + `
WHERE %s`

var stmtAudioTeasers = fmt.Sprintf(stmtAudioBase, "", `
find_in_set(?, story.tag)
ORDER BY story.fileupdatetime DESC
LIMIT ? OFFSET ?`)

var stmtAudioArticle = fmt.Sprintf(stmtAudioBase, `,
story.cshortleadbody AS audio_url,
CONCAT(story.cbyline_description, story.cauthor) AS byline,
story.cbody AS raw_body_json,
story.ebody AS raw_body_xml`,

	`story.id = ?
LIMIT 1`)

const stmtSpeedReadingBase = `
SELECT story.id,
    FROM_UNIXTIME(story.fileupdatetime) AS created_utc,
    FROM_UNIXTIME(story.last_publish_time) AS updated_utc,
    story.tag,
    story.cheadline AS title,
    story.cshortleadbody AS standfirst,
    picture.piclink AS cover_url
    %s
` + stmtInteractiveSource + `
WHERE %s`

var stmtSpeedReadingTeasers = fmt.Sprintf(stmtSpeedReadingBase, "", `find_in_set('速读', story.tag)
ORDER BY story.fileupdatetime DESC
LIMIT ? OFFSET ?`)

var stmtSpeedReadingContent = fmt.Sprintf(stmtSpeedReadingBase, `,
story.clongleadbody AS raw_vocab,
story.eheadline AS title_en,
story.ebody AS raw_body,
story.cbody AS quiz`,

	`story.id = ?
    AND find_in_set('速读', story.tag)
LIMIT 1`)
