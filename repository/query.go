package repository

const (
	stmtVideo = `
    SELECT id AS id,
        fileupdatetime AS created_at,
		pubdate AS updated_at,
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

	stmtGalleryImages = `
    SELECT TRIM(BOTH FROM pic_url) AS image_url,
		TRIM(BOTH FROM pbody) AS caption
	FROM cmstmp01.photonews_picture
	WHERE photonewsid = ?
	ORDER BY orders`

	stmtGallery = `
    SELECT photonewsid AS id,
        FROM_UNIXTIME(add_times) AS created_utc,
        accessright AS access_right,
		cn_title AS title_cn,
		TRIM(BOTH FROM shortlead) AS long_lead_cn,
		leadbody AS body,
		cover AS cover_url
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
