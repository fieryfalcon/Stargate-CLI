package apod

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"stargate/pkg/models"
	"stargate/pkg/utils"
	"strings"
	"syscall"
	"unsafe"

	"image/jpeg"

	_ "image/png"
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

func DetermineFileType(url string) string {
	// Check for common image extensions
	if strings.HasSuffix(url, ".jpg") || strings.HasSuffix(url, ".jpeg") || strings.HasSuffix(url, ".png") || strings.HasSuffix(url, ".gif") || strings.HasSuffix(url, ".bmp") {
		return "image"
	}

	// Check for common video platforms or extensions
	if strings.Contains(url, "youtube.com") || strings.Contains(url, "vimeo.com") || strings.HasSuffix(url, ".mp4") || strings.HasSuffix(url, ".mov") {
		return "video"
	}

	// If the file type is not obvious, try to fetch the content type from the URL
	resp, err := http.Head(url)
	if err != nil {
		return "unknown"
	}
	contentType := resp.Header.Get("Content-Type")

	// Check the content type header
	if strings.HasPrefix(contentType, "image/") {
		return "image"
	}
	if strings.HasPrefix(contentType, "video/") {
		return "video"
	}

	return "unknown"
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

var (
	procSystemParametersInfo = syscall.NewLazyDLL("user32.dll").NewProc("SystemParametersInfoW")
)

const (
	SPI_SETDESKWALLPAPER = 20
)

// ConvertToJPEG converts the image to JPEG format
func ConvertToJPEG(src, dst string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	outFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return jpeg.Encode(outFile, img, nil)
}

func SetWallpaper(filePath string) error {
	// Ensure the image is in a supported format (JPEG)
	jpegPath := filepath.Join(filepath.Dir(filePath), "wallpaper.jpg")
	if err := ConvertToJPEG(filePath, jpegPath); err != nil {
		fmt.Println("Failed to convert image to JPEG:", err)
		return err
	}

	// Use absolute path for the wallpaper file
	absolutePath, err := filepath.Abs(jpegPath)
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return err
	}

	// Set the wallpaper using SystemParametersInfoW
	_, _, err = procSystemParametersInfo.Call(
		uintptr(SPI_SETDESKWALLPAPER),
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(absolutePath))),
		0,
	)
	if err != syscall.Errno(0) {
		fmt.Println("Failed to set wallpaper:", err)
		return err
	}

	// Update the Windows registry to ensure the change is reflected
	cmd := exec.Command("reg", "add", "HKEY_CURRENT_USER\\Control Panel\\Desktop", "/v", "Wallpaper", "/t", "REG_SZ", "/d", absolutePath, "/f")
	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to update registry:", err)
		return err
	}

	// Refresh the desktop to apply the change
	cmd = exec.Command("RUNDLL32.EXE", "user32.dll,UpdatePerUserSystemParameters", "1", "True")
	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to refresh desktop:", err)
		return err
	}

	fmt.Println("Wallpaper set successfully to:", absolutePath)
	return nil
}

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
