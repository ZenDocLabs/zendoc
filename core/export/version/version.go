package version

import (
	"encoding/json"
	"os"

	"github.com/dterbah/zendoc/core/export/helper"
)

type Versions struct {
	Versions []string `json:"versions"`
}

func UpdateVersions(versionPath string, currentVersion string) error {
	if !helper.IsFileExist(versionPath) {
		versions := Versions{
			Versions: []string{currentVersion},
		}
		saveVersions(versionPath, versions)
		return nil
	}

	data, err := os.ReadFile(versionPath)
	if err != nil {
		return err
	}

	var versions Versions
	err = json.Unmarshal(data, &versions)
	if err != nil {
		return err
	}

	for _, v := range versions.Versions {
		if v == currentVersion {
			return nil
		}
	}

	versions.Versions = append(versions.Versions, currentVersion)
	return saveVersions(versionPath, versions)
}

func saveVersions(versionPath string, versions Versions) error {
	data, err := json.MarshalIndent(versions, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(versionPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
