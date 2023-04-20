package pkg

const StmtVideo = `
SELECT id AS id,
    fileupdatetime AS created_date,
    pubdate AS updated_date,
    accessright AS access_right,
    cheadline AS title,
    clongleadbody AS long_lead_cn,
    cc_piclink AS poster_url,
    cc_vaddress AS cc_id,
    cdescribe AS byline_desc_cn,
    cbyline AS byline_cn
FROM cmstmp01.video_story
WHERE id = ?
    AND publish_status = 'publish'`
