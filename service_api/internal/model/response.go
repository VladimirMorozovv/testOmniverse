package model

type Product struct {
	Id    string `json:"id"`
	Price int    `json:"price"`
}

type Error struct {
	Error string `json:"error"`
}
