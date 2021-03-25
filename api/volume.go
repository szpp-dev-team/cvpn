package api

import (
	"errors"
	"strings"
)

const (
	VolumeIDFSShare = "resource_1423533946.487706.3"
	VolumeIDFS      = "resource_1389773645.177066.2," // resource_1389773645.177066.2,2020 とか
)

const (
	VolumeNameFSShare = "fsshare"
	VolumeNameFS      = "fs"
)

var VolumeMap = map[string]string{
	VolumeNameFSShare: VolumeIDFSShare,
	VolumeNameFS:      VolumeIDFS,
}

func GetVolumeIDFromName(name string) (string, error) {
	if name == VolumeNameFSShare {
		return VolumeIDFSShare, nil
	}

	if len(name) >= 2 && name[:2] == VolumeNameFS {
		tokens := strings.Split(name, "/")
		if len(tokens) < 2 {
			return "", errors.New("invalid fs volume format. example: -v fs/2020")
		}

		return VolumeIDFS + tokens[1], nil
	}

	return "", errors.New("no such volume. please try fs/{any} or fsshare")
}
