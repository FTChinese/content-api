package repository

const stmtStory = `
    SELECT story.id AS story_id,
		story.cheadline AS title_cn,
        story.eheadline AS title_en,
		story.clongleadbody AS standfirst,
        picture.piclink AS cover_url,
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
		story.tag AS tag,
        story.topic AS topic,
		FROM_UNIXTIME(story.fileupdatetime) AS created_at,
		FROM_UNIXTIME(story.last_publish_time) AS updated_at,
		story.cbody AS body_cn,
        story.ebody AS body_en
	FROM cmstmp01.story AS story
		LEFT JOIN cmstmp01.story_pic AS storyToPic
            ON story.id = storyToPic.storyid 
        LEFT JOIN cmstmp01.picture AS picture
		    ON picture.id = storyToPic.picture_id
	WHERE story.id = ?
		AND story.publish_status = 'publish'`
