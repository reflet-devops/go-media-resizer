package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetMimeType_Success(t *testing.T) {

	tests := []struct {
		name         string
		searchedType string
		want         string
	}{
		{
			name:         TypeAVIF,
			searchedType: TypeAVIF,
			want:         MimeTypeAVIF,
		},
		{
			name:         TypeWEBP,
			searchedType: TypeWEBP,
			want:         MimeTypeWEBP,
		},
		{
			name:         TypeJPEG,
			searchedType: TypeJPEG,
			want:         MimeTypeJPEG,
		},
		{
			name:         TypePNG,
			searchedType: TypePNG,
			want:         MimeTypePNG,
		},
		{
			name:         TypeGIF,
			searchedType: TypeGIF,
			want:         MimeTypeGIF,
		},
		{
			name:         TypeTIFF,
			searchedType: TypeTIFF,
			want:         MimeTypeTIFF,
		},
		{
			name:         TypeSVG,
			searchedType: TypeSVG,
			want:         MimeTypeSVG,
		},
		{
			name:         TypeText,
			searchedType: TypeText,
			want:         MimeTypeText,
		},
		{
			name:         TypeHTML,
			searchedType: TypeHTML,
			want:         MimeTypeHTML,
		},
		{
			name:         TypeXML,
			searchedType: TypeXML,
			want:         MimeTypeXML,
		},
		{
			name:         TypeJSON,
			searchedType: TypeJSON,
			want:         MimeTypeJSON,
		},
		{
			name:         TypePDF,
			searchedType: TypePDF,
			want:         MimeTypePDF,
		},
		{
			name:         TypeMP4,
			searchedType: TypeMP4,
			want:         MimeTypeMP4,
		},
		{
			name:         TypeWEBM,
			searchedType: TypeWEBM,
			want:         MimeTypeWEBM,
		},
		{
			name:         TypeMEPG,
			searchedType: TypeMEPG,
			want:         MimeTypeMEPG,
		},
		{
			name:         TypeDefault,
			searchedType: "unknown",
			want:         MimeTypeDefault,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetMimeType(tt.searchedType)
			assert.Equal(t, tt.want, got)
		})
	}
}
func Test_GetType(t *testing.T) {

	tests := []struct {
		name              string
		searchedExtension string
		want              string
	}{
		{
			name:              TypeAVIF,
			want:              TypeAVIF,
			searchedExtension: ExtensionAVIF,
		},
		{
			name:              TypeWEBP,
			want:              TypeWEBP,
			searchedExtension: ExtensionWEBP,
		},
		{
			name:              TypeJPEG,
			want:              TypeJPEG,
			searchedExtension: ExtensionJPEG,
		},
		{
			name:              TypeJPEG,
			want:              TypeJPEG,
			searchedExtension: ExtensionJPG,
		},
		{
			name:              TypePNG,
			want:              TypePNG,
			searchedExtension: ExtensionPNG,
		},
		{
			name:              TypeGIF,
			want:              TypeGIF,
			searchedExtension: ExtensionGIF,
		},
		{
			name:              TypeTIFF,
			want:              TypeTIFF,
			searchedExtension: ExtensionTIFF,
		},
		{
			name:              TypeSVG,
			want:              TypeSVG,
			searchedExtension: ExtensionSVG,
		},
		{
			name:              TypeText,
			want:              TypeText,
			searchedExtension: ExtensionText,
		},
		{
			name:              TypeHTML,
			want:              TypeHTML,
			searchedExtension: ExtensionHTML,
		},
		{
			name:              TypeXML,
			want:              TypeXML,
			searchedExtension: ExtensionXML,
		},
		{
			name:              TypeJSON,
			want:              TypeJSON,
			searchedExtension: ExtensionJSON,
		},
		{
			name:              TypePDF,
			want:              TypePDF,
			searchedExtension: ExtensionPDF,
		},
		{
			name:              TypeMP4,
			want:              TypeMP4,
			searchedExtension: ExtensionMP4,
		},
		{
			name:              TypeWEBM,
			want:              TypeWEBM,
			searchedExtension: ExtensionWEBM,
		},
		{
			name:              TypeMEPG,
			want:              TypeMEPG,
			searchedExtension: ExtensionMEPG,
		},
		{
			name:              TypeDefault,
			want:              TypeDefault,
			searchedExtension: ".unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetType(tt.searchedExtension)
			assert.Equal(t, tt.want, got)
		})
	}
}
