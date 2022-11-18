// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package lg

import (
	"os"
	"strconv"
)

func GetenvInt(key string, defaultValue int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	return IfeI(err == nil, value, defaultValue)
}

func GetenvStr(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	return IfeS(exists, value, defaultValue)
}
