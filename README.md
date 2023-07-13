# Content API

A RESTful API for the contents of FTChinese

## Endpoints

* GET `/` show available endpoints
* GET `/front-page/latest` a list of front page article items.
* GET `/front-page/archives/{date}` a list of articles on a specified date.
* GET `/channels` a list of all channel names
* GET `/channels/{name}` a list of articles on a specified channel.
* GET `/contents/stories/{id}` a news story
* GET `/contents/interactive/{id}` an interactive story
* GET `/contents/videos/{id}` a video
* GET `/contents/galleries/{id}` an image gallery
* POST `/starred` add an article to a user's favorite collection
* GET `/starred/{storyID}` check if an article is in user's favorite collection
* DELETE `/starred/{storyID}` delete an article from user's favorite collection.
* GET `__version` Show this's app version information.
* GET `/__status`
* GET `/__status/channel_ids`
