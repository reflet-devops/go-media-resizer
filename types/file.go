package types

import "slices"

const (
	MimeTypeAVIF = "image/avif"
	MimeTypeWEBP = "image/webp"
	MimeTypeJPEG = "image/jpeg"
	MimeTypePNG  = "image/png"
	MimeTypeGIF  = "image/gif"
	MimeTypeTIFF = "image/tiff"
	MimeTypeSVG  = "image/svg+xml"

	MimeTypeText = "text/plain"
	MimeTypeHTML = "text/html"
	MimeTypeXML  = "text/xml"

	MimeTypeDefault = "application/octet-stream"
	MimeTypeJSON    = "application/json"
	MimeTypePDF     = "application/pdf"

	MimeTypeMP4  = "video/mp4"
	MimeTypeWEBM = "video/webm"
	MimeTypeMEPG = "video/mpeg"

	TypeAVIF = "avif"
	TypeWEBP = "webp"
	TypeJPEG = "jpeg"
	TypePNG  = "png"
	TypeGIF  = "gif"
	TypeTIFF = "tiff"
	TypeSVG  = "svg"

	TypeText = "plain"
	TypeHTML = "html"
	TypeXML  = "xml"

	TypeDefault = "default"
	TypeJSON    = "json"
	TypePDF     = "pdf"

	TypeMP4  = "mp4"
	TypeWEBM = "webm"
	TypeMEPG = "mpeg"

	ExtensionAVIF = ".avif"
	ExtensionWEBP = ".webp"
	ExtensionJPEG = ".jpeg"
	ExtensionJPG  = ".jpg"
	ExtensionPNG  = ".png"
	ExtensionGIF  = ".gif"
	ExtensionTIFF = ".tiff"
	ExtensionSVG  = ".svg"

	ExtensionText = ".txt"
	ExtensionHTML = ".html"
	ExtensionXML  = ".xml"

	ExtensionJSON = ".json"
	ExtensionPDF  = ".pdf"

	ExtensionMP4  = ".mp4"
	ExtensionWEBM = ".webm"
	ExtensionMEPG = ".mpeg"

	TypeFormatAuto = "auto"
)

func GetMimeType(code string) string {
	switch code {
	case TypeAVIF:
		return MimeTypeAVIF
	case TypeWEBP:
		return MimeTypeWEBP
	case TypeJPEG:
		return MimeTypeJPEG
	case TypePNG:
		return MimeTypePNG
	case TypeGIF:
		return MimeTypeGIF
	case TypeTIFF:
		return MimeTypeTIFF
	case TypeSVG:
		return MimeTypeSVG
	case TypeText:
		return MimeTypeText
	case TypeHTML:
		return MimeTypeHTML
	case TypeXML:
		return MimeTypeXML
	case TypeJSON:
		return MimeTypeJSON
	case TypePDF:
		return MimeTypePDF
	case TypeMP4:
		return MimeTypeMP4
	case TypeWEBM:
		return MimeTypeWEBM
	case TypeMEPG:
		return MimeTypeMEPG
	default:
		return MimeTypeDefault
	}
}

func GetType(extension string) string {
	switch extension {
	case ExtensionAVIF:
		return TypeAVIF
	case ExtensionWEBP:
		return TypeWEBP
	case ExtensionJPEG:
		return TypeJPEG
	case ExtensionJPG:
		return TypeJPEG
	case ExtensionPNG:
		return TypePNG
	case ExtensionGIF:
		return TypeGIF
	case ExtensionTIFF:
		return TypeTIFF
	case ExtensionSVG:
		return TypeSVG
	case ExtensionText:
		return TypeText
	case ExtensionHTML:
		return TypeHTML
	case ExtensionXML:
		return TypeXML
	case ExtensionJSON:
		return TypeJSON
	case ExtensionPDF:
		return TypePDF
	case ExtensionMP4:
		return TypeMP4
	case ExtensionWEBM:
		return TypeWEBM
	case ExtensionMEPG:
		return TypeMEPG
	default:
		return TypeDefault
	}
}

func ValidateType(fileType string, acceptedTypes []string) bool {
	return slices.Contains(acceptedTypes, fileType)
}
