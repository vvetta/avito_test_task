package teamrepository

type TeamModel struct {
	TeamName string `gorm:"primaryKey"`
}

func (TeamModel) TableName() string {
	return "teams"
}


