package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/gosimple/slug"
	"github.com/nrednav/cuid2"
	"github.com/pelletier/go-toml/v2"
)

var CreateId = CreateIdGenerator(32)
var CreateSmallId = CreateIdGenerator(8)

var CreateUserId = CreateIdGenerator(10)

var CreateArtistId = CreateIdGenerator(10)
var CreateAlbumId = CreateIdGenerator(16)
var CreateTrackId = CreateIdGenerator(32)
var CreateTrackMediaId = CreateIdGenerator(32)

var CreatePlaylistFilterId = CreateIdGenerator(8)
var CreateTrackFilterId = CreateIdGenerator(8)

var CreateVirtualPlaylistId = CreateIdGenerator(16)

var CreateApiTokenId = CreateIdGenerator(32)

var CreateUserListeningEventId = CreateIdGenerator(32)

func CreateIdGenerator(length int) func() string {
	res, err := cuid2.Init(cuid2.WithLength(length))
	if err != nil {
		// TODO(patrik): Change
		log.Fatal("Failed to create id generator", "err", err)
	}

	return res
}

// TODO(patrik): Move to ImageService
func CreateSquareImage(src, dest string) error {
	cmd := exec.Command(
		"magick", src,
		"-gravity", "Center",
		"-extent", "%[fx:min(w,h)]x%[fx:min(w,h)]",
		dest,
	)
	// TODO(patrik): Make this configureble
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// TODO(patrik): Move to ImageService
func CreateResizedImage(src string, dest string, width, height int) error {
	args := []string{
		src,
		"-resize", fmt.Sprintf("%dx%d^", width, height),
		"-gravity", "Center",
		"-extent", fmt.Sprintf("%dx%d", width, height),
		dest,
	}

	cmd := exec.Command("magick", args...)
	// TODO(patrik): Make this configureble
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func Slug(s string) string {
	return slug.Make(s)
}

func ExtractNumber(s string) int {
	n := ""
	for _, c := range s {
		if unicode.IsDigit(c) {
			n += string(c)
		} else {
			break
		}
	}

	if len(n) == 0 {
		return 0
	}

	i, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		return 0
	}

	return int(i)
}

var validImageExts = []string{
	".png",
	".jpg",
	".jpeg",
}

func IsValidImageExt(ext string) bool {
	return slices.Contains(validImageExts, ext)
}

func IsValidImage(p string) (bool, error) {
	cmd := exec.Command("magick", "identify", "-format", "%m", p)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, nil
	}

	format := strings.TrimSpace(string(output))
	return format != "", nil
}

func GetImageFormat(p string) (string, error) {
	cmd := exec.Command("magick", "identify", "-format", "%m", p)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to identify image: %w", err)
	}

	format := strings.TrimSpace(string(output))
	if format == "" {
		return "", fmt.Errorf("could not determine image format")
	}

	return format, nil
}

func ImageFormatToExt(format string) string {
	switch strings.ToLower(format) {
	case "png":
		return ".png"
	case "jpeg", "jpg":
		return ".jpeg"
	default:
		return "." + strings.ToLower(format)
	}
}

func FindImagesInDir(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var images []string
	for _, e := range entries {
		if e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}

		ext := filepath.Ext(e.Name())
		if !IsValidImageExt(ext) {
			continue
		}

		p := filepath.Join(dir, e.Name())
		ok, _ := IsValidImage(p)
		if ok {
			images = append(images, p)
		}
	}

	return images, nil
}

// TODO(patrik): Update this
var validTrackExts []string = []string{
	".wav",
	".flac",
	".opus",
}

func IsValidTrackExt(ext string) bool {
	return slices.Contains(validTrackExts, ext)
}

func CreateDirectories(dirs []string) error {
	for _, dir := range dirs {
		err := os.Mkdir(dir, 0755)
		if err != nil && !errors.Is(err, os.ErrExist) {
			return err
		}
	}

	return nil
}

func ReadToml[T any](p string) (T, error) {
	var res T

	d, err := os.ReadFile(p)
	if err != nil {
		return res, fmt.Errorf("read file: %w", err)
	}

	err = toml.Unmarshal(d, &res)
	if err != nil {
		return res, fmt.Errorf("unmarshal: %w", err)
	}

	return res, nil
}

func ReadJson[T any](p string) (T, error) {
	var res T

	d, err := os.ReadFile(p)
	if err != nil {
		return res, fmt.Errorf("read file: %w", err)
	}

	err = json.Unmarshal(d, &res)
	if err != nil {
		return res, fmt.Errorf("unmarshal: %w", err)
	}

	return res, nil
}
