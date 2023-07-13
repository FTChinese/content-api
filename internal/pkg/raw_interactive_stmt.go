package pkg

const stmtInteractiveFrom = `
FROM cmstmp01.interactive_story AS story
    LEFT JOIN cmstmp01.interactive_pic AS a_p
        ON story.id = a_p.interactive_id
    INNER JOIN cmstmp01.picture AS picture
        ON picture.id = a_p.picture_id
`

const stmtSelectInteractiveBase = `
SELECT story.id,
    story.fileupdatetime AS created_at,
    story.last_publish_time AS updated_at,
    story.cheadline AS title_cn,
    story.clongleadbody AS long_lead_cn,
    story.cshortleadbody AS short_lead_cn,
    story.tag,
    picture.piclink AS cover_url
`

const StmtInteractiveTeaser = stmtSelectInteractiveBase + stmtInteractiveFrom + `
WHERE find_in_set(?, story.tag)
    AND story.publish_status = 'queue'
ORDER BY story.fileupdatetime DESC
LIMIT ? OFFSET ?`

const StmtInteractiveContent = stmtSelectInteractiveBase + `,
    story.eheadline AS title_en,
    story.elongleadbody AS long_lead_en,
    CONCAT(story.cbyline_description, story.cauthor) AS byline_cn,
    story.cbody AS body_cn,
    story.ebody AS body_en
` + stmtInteractiveFrom + `
WHERE story.id = ?
    AND story.publish_status = 'queue'
LIMIT 1`
