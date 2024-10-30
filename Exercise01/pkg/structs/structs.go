package structs

// import "time"

type Data1 struct {
	// Example string
	Nodes      []Article
	Previous   string
	Next       string
	Page       int
	TotalPages int
}

type Article struct {
	ID   int    `gorm:"id"`
	Data string `gorm:"data"`
	Node string `gorm:"node"`
}
