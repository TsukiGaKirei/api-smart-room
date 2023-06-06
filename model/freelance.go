package model

import (
	"time"

	"api-smart-room/database"

	"gorm.io/gorm"
)

type FreelanceData struct {
	IdFreelance  int `gorm:"primaryKey;autoIncrement;"`
	IdUser       string
	IsTrainee    bool
	Rating       float64
	JobDone      int
	DateJoin     time.Time
	Address      string
	AddressLong  float64
	AddressLat   float64
	IsMale       bool
	Dob          time.Time
	Nik          string
	ProfilePict  string
	Points       float64
	JobChildCode string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (FreelanceData) TableName() string {
	return "public.freelance_data"
}

type KeahlianFreelance struct {
	Keahlian string
}

type BidangFreelance struct {
	Bidang string
}

func (fd *FreelanceData) FindNlpTag() (*FreelancerNlp, error) {
	var nlpTags FreelancerNlp
	db := database.GetDBInstance()
	if err := db.Select("nlp_tag1, nlp_tag2, nlp_tag3, nlp_tag4, nlp_tag5").
		Find(&nlpTags, "id_freelance = ?", fd.IdFreelance).Error; err != nil {
		return nil, err
	}

	return &nlpTags, nil
}

func (fd *FreelanceData) FindFreelanceKeahlian() (string, error) {
	db := database.GetDBInstance()

	var data KeahlianFreelance
	res := db.Model(&FreelanceData{}).Select("job_child_code.job_child_name as keahlian").
		Joins(`left join public.job_child_code on job_child_code.job_child_code = public.freelance_data.job_child_code`).
		Where(`public.freelance_data.id_freelance = ?`, fd.IdFreelance).Scan(&data)

	if res.Error != nil {
		return "", res.Error
	}

	if res.RowsAffected == 0 {
		return "", gorm.ErrRecordNotFound
	}

	return data.Keahlian, nil
}

func (fd *FreelanceData) FindFreelanceBidang() (string, error) {
	db := database.GetDBInstance()

	var data BidangFreelance
	res := db.Model(&FreelanceData{}).Select("job_code.job_category as bidang").
		Joins(`left join public.job_child_code job_child_code on job_child_code.job_child_code = public.freelance_data.job_child_code`).
		Joins(`left join public.job_code job_code on job_code.job_code = job_child_code.job_code`).
		Where(`public.freelance_data.id_freelance = ?`, fd.IdFreelance).Scan(&data)

	if res.Error != nil {
		return "", res.Error
	}

	if res.RowsAffected == 0 {
		return "", gorm.ErrRecordNotFound
	}

	return data.Bidang, nil
}
