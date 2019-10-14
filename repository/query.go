package repository

const (
	stmtStory = `
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
		story.fileupdatetime AS created_at,
		story.last_publish_time AS updated_at,
		story.cbody AS body_cn,
        story.ebody AS body_en
	FROM cmstmp01.story AS story
		LEFT JOIN (
			cmstmp01.story_pic AS storyToPic
			INNER JOIN cmstmp01.picture AS picture
		)
		ON story.id = storyToPic.storyid 
		AND picture.id = storyToPic.picture_id
	WHERE story.id = ?
		AND story.publish_status = 'publish'`

	stmtVideo = `
    SELECT id AS id,
		cheadline AS title,
		clongleadbody AS standfirst,
		TRIM(CONCAT(cdescribe, ' ', cbyline)) AS byline,
		cc_vaddress AS ccId,
		cc_piclink AS posterUrl,
		fileupdatetime AS createdAt,
		pubdate AS updatedAt
	FROM cmstmp01.video_story
	WHERE id = ? 
		AND publish_status = 'publish'`

	stmtGalleryImages = `
    SELECT pic_url AS imageUrl,
		pbody AS caption
	FROM cmstmp01.photonews_picture
	WHERE photonewsid = ?
	ORDER BY orders`

	stmtGallery = `
    SELECT photonewsid AS id,
		cn_title AS title,
		shortlead AS standfirst,
		leadbody AS body,
		cover AS coverUrl,
		add_times AS updatedAt,
		tags AS tag
	FROM cmstmp01.photonews
	WHERE photonewsid = ?
	LIMIT 1`

	stmtChannel = `
    SELECT story.id AS id,
      story.cheadline AS titleCn,
      story.clongleadbody AS standfirst,
      story.tag AS tags,
      story.fileupdatetime AS createdAt,
      story.last_publish_time AS updatedAt
    FROM cmstmp01.channel_detail
      INNER JOIN cmstmp01.story
      ON channel_detail.id = story.id
    WHERE channel_detail.chaid = ?
      AND story.publish_status = 'publish'
    ORDER BY channel_detail.addtime DESC
	LIMIT ? OFFSET ?`

	stmtNav = `
    SELECT top.id AS id,
      top.code AS name,
      top.name AS title,
      GROUP_CONCAT(
        CONCAT_WS(':', sub.code, sub.name) ORDER BY sub.id ASC 
        SEPARATOR ';'
      ) AS children
    FROM cmstmp01.channel AS top
      INNER JOIN cmstmp01.channel AS sub
      ON top.id = sub.reid
    WHERE top.reid = 0
    GROUP BY top.id
	ORDER BY top.priority ASC`

	stmtNavMap = `
    SELECT id AS id,
      code AS name
    FROM cmstmp01.channel
    WHERE code IS NOT NULL
	ORDER BY id`
)
