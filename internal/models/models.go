package models

type JSONReq struct {
	URL         string `json:"url,omitempty"`
	CorrID      string `json:"correlation_id,omitempty"`
	OriginalURL string `json:"original_url,omitempty"`
}

type JSONRes struct {
	Result      string `json:"result,omitempty"`
	CorrID      string `json:"correlation_id,omitempty"`
	ShortURL    string `json:"short_url,omitempty"`
	OriginalURL string `json:"-"`
}

type JSONUserRes struct {
	ShortURL    string `json:"short_url,omitempty"`
	OriginalURL string `json:"original_url,omitempty"`
	UserID      string `json:"-"`
}

type KeysToDelete struct {
	Keys   []string
	UserID string
}
