package lg

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
)

const (
	EnvImageTag            = "IMAGE_TAG"
	EnvImageVersionMajor   = "IMAGE_VERSION_MAJOR"
	EnvImageVersionMinor   = "IMAGE_VERSION_MINOR"
	EnvImageVersionPatch   = "IMAGE_VERSION_PATCH"
	EnvImageCommitShortSha = "IMAGE_COMMIT_SHORT_SHA"
	EnvImageBuildDate      = "IMAGE_BUILD_DATE"
	EnvUnknownString       = "__unknown__"
)

var (
	MajorVersion   = 0
	MinorVersion   = 0
	PatchVersion   = 0
	ImageTag       = EnvUnknownString
	CommitShortSHA = EnvUnknownString
	ImageBuildDate = EnvUnknownString
)

// Version 表示程式碼實作的版本
type Version struct {
	Major          int    `json:"major"`
	Minor          int    `json:"minor"`
	Patch          int    `json:"patch"`
	BuildDate      string `json:"buildDate"`
	ImageTag       string `json:"imageTag"`
	CommitShortSHA string `json:"commitShortSHA"`
	Route          string `json:"route"`
}

func LoadVerEnv(filenames ...string) (*Version, error) {
	err := godotenv.Load(filenames...)
	if err != nil {
		log.Println("Error loading image_version.env file")
		return nil, errors.New("loading image_version.env file failed")
	} else {
		MajorVersion = GetenvInt(EnvImageVersionMajor, 0)
		MinorVersion = GetenvInt(EnvImageVersionMinor, 0)
		PatchVersion = GetenvInt(EnvImageVersionPatch, 0)
		ImageTag = GetenvStr(EnvImageTag, EnvUnknownString)
		CommitShortSHA = GetenvStr(EnvImageCommitShortSha, EnvUnknownString)
		ImageBuildDate = GetenvStr(EnvImageBuildDate, EnvUnknownString)
		return &Version{
			Major:          MajorVersion,
			Minor:          MinorVersion,
			Patch:          PatchVersion,
			BuildDate:      ImageBuildDate,
			ImageTag:       ImageTag,
			CommitShortSHA: CommitShortSHA,
		}, nil
	}
}
