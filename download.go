package main

import "google.golang.org/api/drive/v3"

var (
	ErrNeedRecursive = error.New("to download a folder, use the '-r' flag")
)

func (cfg *apiConfig) downloadFile(path string, targetPath string, flags map[string]string) (*drive.File, error) {
	// look for file at drive location
	dirId, filename, isFolder, err := cfg.ValidateUploadPath(path)
	if err != nil {
		return nil, err
	}

	if isFolder {
		if _, ok := flags["r"]; ok {
			return cfg.downloadFolder(path, targetPath, flags)
		}
		return nil, 
	}
	// if exists, download file

	return nil, nil
}

func (cfg *apiConfig) downloadFolder(path string, targetPath string, flags map[string]string) (*drive.File, error) {

}
