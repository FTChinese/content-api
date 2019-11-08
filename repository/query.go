package repository

const (
	stmtVideo = `
    SELECT id AS id,
        fileupdatetime AS created_at,
		pubdate AS updated_at,
        accessright AS access_right,
		cheadline AS title,
		clongleadbody AS standfirst,
        cc_piclink AS poster_url,
		cc_vaddress AS cc_id,
		TRIM(CONCAT(cdescribe, ' ', cbyline)) AS byline
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

	stmtOAuth = `
    SELECT access_token,
        is_active,
        expires_in,
        created_utc
    FROM oauth.access
    WHERE access_token = UNHEX(?)
    LIMIT 1`
)
