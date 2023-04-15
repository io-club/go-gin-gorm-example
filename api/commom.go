package api


type Pageable struct {
	Page *int `form:"page" json:"page" binding:"required"`
	Size *int `form:"size" json:"size" binding:"required"`
}
