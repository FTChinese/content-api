package repository

const (
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
)
