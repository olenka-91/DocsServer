package entity

type LimitedDocsListInput struct {
	Token string `json:"token"`
	Login string `json:"login"` //опционально — если не указан — то список своих
	Key   string `json:"key"`   //имя колонки для фильтрации
	Value string `json:"value"` //- значение фильтра
	Limit int    `json:"limit"` //кол-во документов в списке
}

type Document struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Mime    string   `json:"mime"`
	File    bool     `json:"file"`
	Public  bool     `json:"public"`
	Created string   `json:"created_at"`
	Grant   []string `json:"grant"`
}

type DocsData struct {
	Docs []Document `json:"docs"`
}

type DocsResponse struct {
	Data DocsData `json:"data"`
}
