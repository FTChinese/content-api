package pkg

const StmtListStarred = `
SELECT fav.userid       AS user_id,
	fav.storyid         AS id,
	fav.type            AS kind,
    fav.subtime         AS starred_at,
	story.cheadline     AS story_title,
	story.clongleadbody AS story_lead,
	inter.cheadline     AS interact_title,
	inter.clongleadbody AS interact_lead,
    video.cheadline     AS video_title,
    video.clongleadbody AS video_lead,
    photo.cn_title      AS photo_title,
    photo.shortlead     AS photo_lead
FROM user_db.user_favorites as fav
LEFT JOIN cmstmp01.story as story
	ON fav.storyid = story.id
LEFT JOIN cmstmp01.interactive_story as inter
	ON fav.storyid = inter.id
LEFT JOIN cmstmp01.video_story as video
    ON fav.storyid = video.id
LEFT JOIN cmstmp01.photonews as photo
    ON fav.storyid = photo.photonewsid
WHERE FIND_IN_SET(fav.userid, ?) > 0
ORDER BY fav.subtime DESC
LIMIT ? OFFSET ?`

const StmtCountStarred = `
SELECT COUNT(*) AS row_count
FROM user_db.user_favorites AS fav
WHERE FIND_IN_SET(fav.userid, ?) > 0`

const sqlCSTNow = "DATE_ADD(UTC_TIMESTAMP(), INTERVAL 8 HOUR)"

const StmtSaveStar = `
INSERT IGNORE INTO user_db.user_favorites
SET userid = :user_id,
    storyid = :id,
    type = :kind,
    subtime = ` + sqlCSTNow

const StmtUnstar = `
DELETE FROM user_db.user_favorites
WHERE FIND_IN_SET(userid, ?) > 0
    AND storyid = ?
LIMIT 1`

const StmtIsStarring = `
SELECT EXISTS (
    SELECT *
    FROM user_db.user_favorites
    WHERE FIND_IN_SET(userid, ?) > 0
        AND storyid = ?
) AS alreadyExists`
