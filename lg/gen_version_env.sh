#!/bin/bash
# Copyright 2022 Yuchi Chen. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# 執行的目錄與與產生的 env 檔案相同
# 複製來源：lc-go/lg/gen_version_env.sh

COMMIT_TAGS=$(git tag --contains)
if [ -z "$COMMIT_TAGS" ]
then
  echo "No available tag."
  exit 0
fi

for COMMIT_TAG in $COMMIT_TAGS
do
  #echo "check $COMMIT_TAG"
  if [[ "$COMMIT_TAG" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]
  then
    #echo "good"
    break
  fi
done

if ! [[ "$COMMIT_TAG" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]
then
  echo "Incorrect tag format. ($COMMIT_TAG)"
  echo "All tags: $COMMIT_TAGS"
  exit 0
fi

VERSION=$(echo $COMMIT_TAG | tr "v" ".")
#echo $VERSION
MAJOR=$(echo $VERSION | cut -d. -f2)
MINOR=$(echo $VERSION | cut -d. -f3)
PATCH=$(echo $VERSION | cut -d. -f4)
#echo "Major: $MAJOR"
#echo "Minor: $MINOR"
#echo "Patch: $PATCH"

COMMIT_SHA=${CI_COMMIT_SHA:-$(git rev-parse HEAD)}
if [ -z "$COMMIT_SHA" ]
then
  echo "No available commit."
  exit 1
fi
COMMIT_SHORT_SHA=${COMMIT_SHA:0:8}

#echo $COMMIT_SHA
#echo $COMMIT_SHORT_SHA

BUILD_DATE_STR=`date +%y%m%d`
#echo $BUILD_DATE_STR

IMAGE_TAG="$MAJOR.$MINOR.$PATCH.$BUILD_DATE_STR-$COMMIT_SHORT_SHA"
echo "Preparing image with tag name: $IMAGE_TAG"

echo "# 這個檔案會自動產生，不需要編輯。\n# 有沒有commit也無所謂。" > image_version.env
echo "IMAGE_TAG=${IMAGE_TAG}" >> image_version.env
echo "IMAGE_VERSION_MAJOR=$MAJOR" >> image_version.env
echo "IMAGE_VERSION_MINOR=$MINOR" >> image_version.env
echo "IMAGE_VERSION_PATCH=$PATCH" >> image_version.env
echo "IMAGE_COMMIT_SHORT_SHA=${COMMIT_SHORT_SHA}" >> image_version.env
echo "IMAGE_BUILD_DATE=${BUILD_DATE_STR}" >> image_version.env

