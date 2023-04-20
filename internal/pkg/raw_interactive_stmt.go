package pkg

const stmtInteractiveFrom = `
FROM cmstmp01.interactive_story AS story
    LEFT JOIN cmstmp01.interactive_pic AS a_p
        ON story.id = a_p.interactive_id
    INNER JOIN cmstmp01.picture AS picture
        ON picture.id = a_p.picture_id
`

const StmtInteractiveTeaser = `
SELECT story.id,
    story.fileupdatetime AS created_at,
    story.last_publish_time AS updated_at,
    story.cheadline AS title_cn,
    story.clongleadbody AS long_lead_cn,
    story.cshortleadbody AS short_lead_cn,
    story.tag,
    picture.piclink AS cover_url
` + stmtInteractiveFrom + `
WHERE find_in_set(?, story.tag)
    AND story.publish_status = 'queue'
ORDER BY story.fileupdatetime DESC
LIMIT ? OFFSET ?`

const StmtInteractiveContent = `
SELECT story.id AS id,
    story.fileupdatetime AS created_utc,
    story.last_publish_time AS updated_utc,
    story.cheadline AS title_cn,
    story.eheadline AS title_en,
    story.clongleadbody AS long_lead_cn,
    story.elongleadbody AS long_lead_en,
    story.cshortleadbody AS short_lead_cn,
    CONCAT(story.cbyline_description, story.cauthor) AS byline_cn,
    story.tag AS tag,
    picture.piclink AS cover_url,
    story.cbody AS body_cn,
    story.ebody AS body_en
` + stmtInteractiveFrom + `
WHERE story.id = ?
    AND story.publish_status = 'queue'
LIMIT 1`
