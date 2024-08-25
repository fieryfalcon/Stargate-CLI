package apod

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"stargate/pkg/models"
	"stargate/pkg/utils"
	"io"
	"runtime"
	"strings"
	"io/ioutil"
)

const apiURL = "https://api.nasa.gov/planetary/apod"

func FetchTodayAPOD() (*models.APODResponse, error) {
	url := fmt.Sprintf("%s?api_key=%s", apiURL, utils.GetAPIKey())
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apod models.APODResponse
	if err := json.NewDecoder(resp.Body).Decode(&apod); err != nil {
		return nil, err
	}
	return &apod, nil
}
func DownloadAPODAndSave(url, date string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create directory with the date
	dirPath := filepath.Join("apod_images", date)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return err
	}

	// Save the image to the directory
	filePath := filepath.Join(dirPath, "apod.jpg")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// OpenImage opens an image using the default image viewer
func OpenImage(filePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", filePath)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", filePath)
	case "linux":
		cmd = exec.Command("xdg-open", filePath)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}

// SetWallpaper sets an image as the wallpaper
func SetWallpaper(filePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("osascript", "-e", fmt.Sprintf("tell application \"Finder\" to set desktop picture to POSIX file \"%s\"", filePath))
	case "windows":
		cmd = exec.Command("reg", "add", "HKEY_CURRENT_USER\\Control Panel\\Desktop", "/v", "Wallpaper", "/t", "REG_SZ", "/d", filePath, "/f")
		cmd = exec.Command("RUNDLL32.EXE", "user32.dll,UpdatePerUserSystemParameters")
	case "linux":
		cmd = exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", "file://"+filePath)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Run()
}

// CopyToFavorites copies an image to the favorites folder
func CopyToFavorites(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func SetSlideshowFromFavorites(favoritesPath string) error {
	files, err := ioutil.ReadDir(favoritesPath)
	if err != nil {
		return err
	}

	var imagePaths []string
	for _, file := range files {
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".jpg") || strings.HasSuffix(file.Name(), ".png")) {
			imagePaths = append(imagePaths, filepath.Join(favoritesPath, file.Name()))
		}
	}

	if len(imagePaths) == 0 {
		return fmt.Errorf("no images found in the favorites folder")
	}

	switch runtime.GOOS {
	case "darwin":
		return setMacOSWallpaperSlideshow(imagePaths)
	case "windows":
		return setWindowsWallpaperSlideshow(imagePaths)
	case "linux":
		return setLinuxWallpaperSlideshow(imagePaths)
	default:
		return fmt.Errorf("unsupported platform")
	}
}

func setMacOSWallpaperSlideshow(imagePaths []string) error {
	// macOS specific implementation (not implemented here)
	return fmt.Errorf("macOS slideshow setting not implemented")
}

func setWindowsWallpaperSlideshow(imagePaths []string) error {
	// Windows specific implementation (not implemented here)
	return fmt.Errorf("Windows slideshow setting not implemented")
}

func setLinuxWallpaperSlideshow(imagePaths []string) error {
	// Linux specific implementation (not implemented here)
	return fmt.Errorf("Linux slideshow setting not implemented")
}