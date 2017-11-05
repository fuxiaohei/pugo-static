package models

// Archive is list of archive by year
type Archive struct {
	Year  int
	Posts []*Post
}

// NewArchives return archives list with posts
func NewArchives(posts []*Post) []*Archive {
	var (
		archives       []*Archive
		currentArchive *Archive
	)
	for _, p := range posts {
		if currentArchive != nil && currentArchive.Year != p.Created().Year() {
			archives = append(archives, currentArchive)
			currentArchive = nil
		}
		if currentArchive == nil {
			currentArchive = &Archive{
				Year:  p.Created().Year(),
				Posts: []*Post{p},
			}
			continue
		}
		currentArchive.Posts = append(currentArchive.Posts, p)
	}
	if currentArchive != nil {
		archives = append(archives, currentArchive)
	}
	return archives
}
