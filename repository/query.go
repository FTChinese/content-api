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

	stmtOAuth = `
    SELECT access_token,
        is_active,
        expires_in,
        created_utc
    FROM oauth.access
    WHERE access_token = UNHEX(?)
    LIMIT 1`
)
