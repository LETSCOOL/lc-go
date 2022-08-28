package lg

import (
	"github.com/joho/godotenv"
	"log"
)

const (
	EnvImageTag            = "IMAGE_TAG"
	EnvImageVersionMajor   = "IMAGE_VERSION_MAJOR"
	EnvImageVersionMinor   = "IMAGE_VERSION_MINOR"
	EnvImageCommitShortSha = "IMAGE_COMMIT_SHORT_SHA"
	EnvImageBuildDate      = "IMAGE_BUILD_DATE"
	EnvUnknownString       = "__unknown__"
)

var (
	MajorVersion   = 0
	MinorVersion   = 0
	ImageTag       = EnvUnknownString
	CommitShortSHA = EnvUnknownString
	ImageBuildDate = EnvUnknownString
)

func LoadVerEnv(filenames ...string) {
	err := godotenv.Load(filenames...)
	if err != nil {
		log.Println("Error loading image_version.env file")
	} else {
		MajorVersion = GetenvInt(EnvImageVersionMajor, 0)
		MinorVersion = GetenvInt(EnvImageVersionMinor, 0)
		ImageTag = GetenvStr(EnvImageTag, EnvUnknownString)
		CommitShortSHA = GetenvStr(EnvImageCommitShortSha, EnvUnknownString)
		ImageBuildDate = GetenvStr(EnvImageBuildDate, EnvUnknownString)
	}
}
