package model

import "gorm.io/gorm"

// Comment 评论表
type Comment struct {
	gorm.Model
	UserId   uint64 `gorm:"type:bigint(20);not null;unique" json:"user_id"`
	Content  string `gorm:"type:text;not null" json:"content"`
	Type     uint   `gorm:"type:tinyint(4);not null" json:"type"`
	TargetId uint   `gorm:"type:bigint(20);not null" json:"target_id"`
	ParentId uint   `gorm:"type:bigint(20);default null" json:"parent_id"`
	Likes    uint   `gorm:"type:int(11);default 0" json:"likes"`

	// 添加评论与用户的多对一关系
	User User `gorm:"foreignKey:UserId" json:"-"`
	// 添加评论的自引用关系（父子评论）
	Children []Comment `gorm:"foreignKey:ParentId" json:"-"`
}

//'1:歌曲,2:歌单,3:专辑'
