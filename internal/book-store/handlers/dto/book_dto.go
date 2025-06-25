package dto

type CreateBookRequest struct {
	Title    string `json:"title"`
	Author   string `json:"author"`
	Genre    string `json:"genre"`
	Year     int    `json:"year"`
	Quantity int    `json:"quantity"`
}

type UpdateBookRequest struct {
	Title    *string `json:"title,omitempty"`
	Author   *string `json:"author,omitempty"`
	Genre    *string `json:"genre,omitempty"`
	Year     *int    `json:"year,omitempty"`
	Quantity *int    `json:"quantity,omitempty"`
}

type BookFilterRequest struct {
	Title  string
	Author string
	Genre  string
	Year   *int
	Limit  int
	Offset int
}
