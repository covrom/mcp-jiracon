package utils

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"
)

// AttachmentMaxBytes is the inline attachment size cap (50 MiB).
// Mirrors utils.media.ATTACHMENT_MAX_BYTES.
const AttachmentMaxBytes = 50 * 1024 * 1024

// imageMimeTypes is the explicit set of recognized image MIME types.
var imageMimeTypes = map[string]bool{
	"image/png":     true,
	"image/jpeg":    true,
	"image/gif":     true,
	"image/webp":    true,
	"image/svg+xml": true,
	"image/bmp":     true,
}

// imageExts is the fallback extension set used when the MIME is
// application/octet-stream.
var imageExts = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true,
	".webp": true, ".svg": true, ".bmp": true,
}

// IsImageAttachment reports whether the (mediaType, filename) pair is
// recognized as an image. Returns (isImage, resolvedMIME).
func IsImageAttachment(mediaType, filename string) (bool, string) {
	mt := strings.ToLower(strings.TrimSpace(mediaType))
	if mt != "" && imageMimeTypes[mt] {
		return true, mt
	}
	if mt == "" || mt == "application/octet-stream" || mt == "application/binary" {
		guess := guessImageMIME(filename)
		if guess != "" {
			return true, guess
		}
	}
	return false, "application/octet-stream"
}

// guessImageMIME returns the image MIME type guessed from the filename
// extension. Returns "" if no match.
func guessImageMIME(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if !imageExts[ext] {
		return ""
	}
	mimeType := mime.TypeByExtension(ext)
	if mimeType != "" {
		return mimeType
	}
	// Fallback to canonical names.
	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".bmp":
		return "image/bmp"
	}
	return ""
}

// FetchAndEncodeAttachment fetches a URL, base64-encodes its bytes,
// and returns (base64, mimeType, bytesRead). On size-overrun returns
// (nil, nil, totalBytes).
func FetchAndEncodeAttachment(fetch func() (io.ReadCloser, error), mimeType string, maxBytes ...int64) (string, string, int, error) {
	max := int64(AttachmentMaxBytes)
	if len(maxBytes) > 0 && maxBytes[0] > 0 {
		max = maxBytes[0]
	}
	rc, err := fetch()
	if err != nil {
		return "", "", 0, err
	}
	defer rc.Close()
	data, err := io.ReadAll(io.LimitReader(rc, max+1))
	if err != nil {
		return "", "", 0, err
	}
	if int64(len(data)) > max {
		return "", "", len(data), fmt.Errorf("attachment exceeds %d bytes", max)
	}
	resolved := mimeType
	if resolved == "" {
		resolved = guessImageMIME("") // empty; caller passes filename for fallback
	}
	return base64.StdEncoding.EncodeToString(data), resolved, len(data), nil
}
