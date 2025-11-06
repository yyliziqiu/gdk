package xutil

import (
	"path/filepath"
)

var _mimeTypes = map[string]string{
	// 文本类型
	".js":   "text/javascript",
	".txt":  "text/plain",
	".css":  "text/css",
	".xml":  "application/xml",
	".htm":  "text/html",
	".html": "text/html",
	".json": "application/json",

	// 图片类型
	".svg":  "image/svg+xml",
	".gif":  "image/gif",
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".webp": "image/webp",

	// 视频类型
	".mp4":  "video/mp4",
	".avi":  "video/x-msvideo",
	".mov":  "video/quicktime",
	".mkv":  "video/x-matroska",
	".flv":  "video/x-flv",
	".webm": "video/webm",

	// 音频类型
	".mp3": "audio/mpeg",
	".wav": "audio/wav",
	".aac": "audio/aac",
	".ogg": "audio/ogg",
	".m4a": "audio/mp4",

	// 文档类型
	".pdf":  "application/pdf",
	".doc":  "application/msword",
	".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	".xls":  "application/vnd.ms-excel",
	".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	".ppt":  "application/vnd.ms-powerpoint",
	".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	".rtf":  "application/rtf",
	".odt":  "application/vnd.oasis.opendocument.text",
	".ods":  "application/vnd.oasis.opendocument.spreadsheet",
	".odp":  "application/vnd.oasis.opendocument.presentation",

	// 压缩文件类型
	".zip": "application/zip",
	".rar": "application/x-rar-compressed",
	".tar": "application/x-tar",
	".bz2": "application/x-bzip2",
	".gz":  "application/gzip",
	".7z":  "application/x-7z-compressed",

	// 其他常见类型
	".csv":     "text/csv",
	".tsv":     "text/tab-separated-values",
	".torrent": "application/x-bittorrent",
}

func ParseMimeType(path string) string {
	if typ, ok := _mimeTypes[filepath.Ext(path)]; ok {
		return typ
	}
	return "application/octet-stream"
}
