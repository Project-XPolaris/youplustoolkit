package youlibrary

type MatchSubjectResponse struct {
	Success bool    `json:"success"`
	Data    Subject `json:"data"`
}
type GetSubjectResponse struct {
	Success bool    `json:"success"`
	Data    Subject `json:"data"`
}
type Subject struct {
	Id         uint         `json:"id"`
	Title      string       `json:"title"`
	Cover      string       `json:"cover"`
	Characters []*Character `json:"characters"`
	Release    string       `json:"release"`
	MalId      string       `json:"malId"`
	BangumiId  string       `json:"bangumiId"`
	TmdbId     string       `json:"tmdbId"`
}

type Person struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Pic       string `json:"pic"`
	MalId     string `json:"malId"`
	BangumiId string `json:"bangumiId"`
	TmdbId    string `json:"tmdbId"`
}

type Character struct {
	Name      string  `json:"name"`
	Pic       string  `json:"pic"`
	Person    *Person `json:"person"`
	MalId     string  `json:"malId"`
	BangumiId string  `json:"bangumiId"`
	TmdbId    string  `json:"tmdbId"`
}
