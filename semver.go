package search

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type semver struct {
	Major int
	Minor int
	Patch int
}

func parseSemVer(str string) (v semver, err error) {
	splits := strings.Split(str, ".")
	if len(splits) != 3 {
		err = fmt.Errorf("invalid version: %s", str)
	} else if v.Major, err = strconv.Atoi(splits[0]); err != nil {
		err = fmt.Errorf("invalid version: %s", str)
	} else if v.Minor, err = strconv.Atoi(splits[1]); err != nil {
		err = fmt.Errorf("invalid version: %s", str)
	} else if v.Patch, err = strconv.Atoi(splits[2]); err != nil {
		err = fmt.Errorf("invalid version: %s", str)
	}
	return
}

// canLoad returns true if `v` can use data from `other`.
func (v *semver) canLoad(other semver) bool {
	if v.Major != other.Major {
		return false
	}
	if v.Minor == other.Minor {
		return v.Patch >= other.Patch
	}
	return v.Minor >= other.Minor
}

func (v *semver) MarshalJSON() ([]byte, error) {
	joined := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	return json.Marshal(joined)
}

func (v *semver) UnmarshalJSON(data []byte) (err error) {
	var str string
	if err = json.Unmarshal(data, &str); err == nil {
		*v, err = parseSemVer(str)
	}
	return
}
