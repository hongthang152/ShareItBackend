package models

type File struct {
	ID   int64  `json:"id" binding:"required"`
	name string `json:"name" binding:"required"`
	PIN  int64  `json:"pin" binding:"required"`
	url  string `json:"url" binding:"required"`
}
