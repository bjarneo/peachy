package shared

import "testing"

func TestIsValidImage(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"photo.jpg", true},
		{"photo.JPG", true},
		{"photo.jpeg", true},
		{"photo.JPEG", true},
		{"photo.png", true},
		{"photo.PNG", true},
		{"photo.webp", true},
		{"photo.gif", true},
		{"photo.bmp", true},
		{"photo.tiff", true},
		{"/path/to/image.jpg", true},
		{"/path/to/image.PNG", true},
		{"photo.txt", false},
		{"photo.pdf", false},
		{"photo.mp4", false},
		{"photo", false},
		{"", false},
		{".jpg", true},
		{"photo.jpg.bak", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := IsValidImage(tt.path)
			if got != tt.want {
				t.Errorf("IsValidImage(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
