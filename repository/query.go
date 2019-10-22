package repository

const (
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

	stmtOAuth = `
    SELECT access_token,
        is_active,
        expires_in,
        created_utc
    FROM oauth.access
    WHERE access_token = UNHEX(?)
    LIMIT 1`
)
