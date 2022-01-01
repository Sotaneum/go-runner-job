package runnerjob

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	file "github.com/Sotaneum/go-json-file"
	"github.com/Sotaneum/go-runner"
)

type JobAdmin struct {
	Owner   string   `json:"owner"`
	Members []string `json:"members"`
}

type BaseJob struct {
	runner.JobInterface
	Active     bool     `json:"active"`
	ID         string   `json:"id"`
	Admin      JobAdmin `json:"admin"`
	CreateDate string   `json:"createDate"`
}

// IsRun : Job를 실행해야하는 타임인지 여부를 반환합니다.
func (job *BaseJob) IsRun(t time.Time) bool {
	return job.Active
}

// Save : Job를 파일로 저장합니다.
func (job *BaseJob) Save(path string) {
	f := file.File{Path: path, Name: job.ID + ".json"}
	f.Remove()
	f.SaveObject(job)
}

// Remove : Job를 삭제합니다.
func (job *BaseJob) Remove(path string) error {
	f := file.File{Path: path, Name: job.ID + ".json"}
	return f.Remove()
}

// GetDefaultOwner : 기본 관리자를 수정하려면, 임베딩을 사용하세요.
func (job *BaseJob) GetDefaultOwner() string {
	return "admin"
}

// Run : Job를 실행합니다. 임베딩을 사용하세요.
func (job *BaseJob) Run() interface{} {
	return "run"
}

// SetOwner : Job의 주인을 설정합니다.
func (job *BaseJob) SetOwner(member string) {
	if member == "" {
		job.SetOwner(job.GetDefaultOwner())
		return
	}
	job.Admin.Owner = member
	for _, user := range job.Admin.Members {
		if user == member {
			return
		}
	}
	job.Admin.Members = append(job.Admin.Members, member)
}

// GetOwner : owner 정보를 반환합니다.
func (job *BaseJob) GetOwner() string {
	return job.Admin.Owner
}

// GetID : Job ID를 반환합니다.
func (job *BaseJob) GetID() string {
	return job.ID
}

// CreateID : ID를 생성합니다.
func (job *BaseJob) CreateID() error {
	hash := sha256.New()
	data, err := json.Marshal(job)

	if err != nil {
		return err
	}

	hash.Write(data)

	job.ID = hex.EncodeToString(hash.Sum(nil))

	return nil
}

// HasAuthorization : 주어진 멤버가 이 Job에 권한이 있는지 여부를 반환합니다.
func (job *BaseJob) HasAuthorization(member string) bool {
	if job.HasAdminAuthorization(member) {
		return true
	}

	members := job.Admin.Members

	for _, m := range members {
		if m == member {
			return true
		}
	}
	return false
}

// HasAdminAuthorization : Job를 관리자 수준까지 권한이 있는지 여부를 반한합니다.
func (job *BaseJob) HasAdminAuthorization(member string) bool {
	return job.GetOwner() == member
}

// IsAvailability : 데이터가 유효성이 존재하는 지 여부를 반환합니다.
func (job *BaseJob) IsAvailability() bool {
	return job.GetOwner() != ""
}
