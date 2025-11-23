package userrepository

type UserModel struct {
    UserID   string `gorm:"primaryKey;column:user_id"`
    TeamName string `gorm:"column:team_name;index"`
    Username string
    IsActive bool
}

func (UserModel) TableName() string { return "users" }
