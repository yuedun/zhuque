package server

import "time"

type Server struct {
	ID        int       `json:"id"`
	Name  string    `json:"name"`
	IP     string    `json:"ip"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// 设置User的表名为`user`
func (Server) TableName() string {
	return "server"
}
