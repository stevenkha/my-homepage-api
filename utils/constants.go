package utils

const (
	AnimeUrl       = "https://yugenanime.tv/mylist/"
	AnimeListClass = "list-entries"

	MangaUrl = "https://user.mngusr.com/bookmark_get_list_full"
)

type ItemInfo struct {
	Cover       string `json:"cover"`
	Title       string `json:"title"`
	Viewed      string `json:"viewed"`
	Current     string `json:"current"`
	CurrentLink string `json:"currentLink"`
}
