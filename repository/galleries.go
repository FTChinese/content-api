package repository

import "gitlab.com/ftchinese/content-api/models"

type imageResult struct {
	success []models.GalleryItem
	err     error
}

type galleryResult struct {
	success models.Gallery
	err     error
}

func (env ContentEnv) RetrieveGalleryImages(id int64) ([]models.GalleryItem, error) {
	var items = make([]models.GalleryItem, 0)

	if err := env.db.Select(&items, stmtGalleryImages, id); err != nil {
		return []models.GalleryItem{}, err
	}

	return items, nil
}

func (env ContentEnv) RetrieveGalleryBody(id int64) (models.Gallery, error) {
	var data models.RawGallery

	if err := env.db.Get(&data, stmtGallery, id); err != nil {
		return models.Gallery{}, err
	}

	return models.NewGallery(&data), nil
}

func (env ContentEnv) RetrieveGallery(id int64) (models.Gallery, error) {
	imageChan := make(chan imageResult)
	bodyChan := make(chan galleryResult)

	go func() {
		images, err := env.RetrieveGalleryImages(id)
		imageChan <- imageResult{
			success: images,
			err:     err,
		}
	}()

	go func() {
		body, err := env.RetrieveGalleryBody(id)
		bodyChan <- galleryResult{
			success: body,
			err:     err,
		}
	}()

	iResult, gResult := <-imageChan, <-bodyChan

	if iResult.err != nil {
		return models.Gallery{}, iResult.err
	}

	if gResult.err != nil {
		return models.Gallery{}, gResult.err
	}

	gResult.success.Items = iResult.success

	return gResult.success, nil
}
