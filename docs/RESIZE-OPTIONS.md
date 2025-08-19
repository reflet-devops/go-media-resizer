# Image Transformation Options

This document describes all available options for image transformation in the `ResizeOption` structure. These options can be used in endpoint configurations, CDN-CGI URLs, and default resize settings.

## Quick Reference

| Parameter | Type | Description | Default | CDN-CGI Support |
|-----------|------|-------------|---------|-----------------|
| `width` | Integer | Image width in pixels | 0 (original) | ✅ |
| `height` | Integer | Image height in pixels | 0 (original) | ✅ |
| `quality` | Integer | JPEG compression quality (1-100) | 0 (default) | ✅ |
| `format` | String | Output image format | `"auto"` | ✅ |
| `fit` | String | Resize method | `"resize"` | ✅ |
| `blur` | Float | Blur radius | 0 (no blur) | ✅ |
| `brightness` | Float | Brightness adjustment | 0 (no change) | ✅ |
| `contrast` | Float | Contrast adjustment | 0 (no change) | ✅ |
| `saturation` | Float | Saturation adjustment | 0 (no change) | ✅ |
| `sharpen` | Float | Sharpening amount | 0 (no sharpening) | ✅ |
| `gamma` | Float | Gamma correction | 0 (no correction) | ✅ |

## Detailed Parameters

### Width
**Type:** Integer  
**Range:** 1-9999 pixels  
**Default:** 0 (original width)  
**CDN-CGI:** `width=500`

Controls the output image width in pixels.

```yaml
# Configuration
default_resize:
  width: 800

# CDN-CGI
  /cdn-cgi/image/width=800/source.jpg

  # URL Pattern
  /resize/800/image.jpg        # width=800
  /resize/800x600/image.jpg    # width=800, height=600
```

**Behavior:**
- `0`: Keep original width
- `> 0`: Resize to specified width
- When only width is specified, height is calculated proportionally (unless `fit: "crop"`)
- Maximum supported width: 9999 pixels

---

### Height
**Type:** Integer  
**Range:** 1-9999 pixels  
**Default:** 0 (original height)  
**CDN-CGI:** `height=300`

Controls the output image height in pixels.

```yaml
# Configuration
default_resize:
  height: 600

# CDN-CGI
  /cdn-cgi/image/height=600/source.jpg

  # URL Pattern
  /resize/x600/image.jpg       # height=600
  /resize/800x600/image.jpg    # width=800, height=600
```

**Behavior:**
- `0`: Keep original height
- `> 0`: Resize to specified height
- When only height is specified, width is calculated proportionally (unless `fit: "crop"`)
- Maximum supported height: 9999 pixels

---

### Quality
**Type:** Integer  
**Range:** 1-100  
**Default:** 95
**CDN-CGI:** `quality=95`

Controls JPEG compression quality. Only applies to JPEG output format.

```yaml
# Configuration
default_resize:
  quality: 90

# CDN-CGI
/cdn-cgi/image/quality=90,format=jpeg/source.png

# URL Pattern
/resize/800x600-90/image.jpg  # quality=90
```

**Quality Guidelines:**
- `95-100`: Highest quality, largest file size
- `85-95`: High quality, good for detailed images
- `75-85`: Good quality, balanced file size (recommended)
- `60-75`: Acceptable quality, smaller file size
- `1-60`: Lower quality, smallest file size

**Behavior:**
- `0`: Uses encoder default quality
- Only affects JPEG output
- Ignored for PNG, WebP, AVIF (they use their own compression)

---

### Format
**Type:** String  
**Values:** `"auto"`, `"jpeg"`, `"png"`, `"webp"`, `"avif"`  
**Default:** `"auto"`  
**CDN-CGI:** `format=webp`

Controls the output image format.

```yaml
# Configuration
default_resize:
  format: "webp"

# CDN-CGI
/cdn-cgi/image/format=avif/source.jpg

# URL Pattern
/resize/800x600-90-webp/image.jpg  # format=webp
```

#### Format Options

**`auto`** (Recommended)
- Automatically selects the best format based on client's `Accept` header
- Priority: AVIF > WebP > Original format
- Provides optimal file size and quality

**`jpeg`**
- JPEG format with quality control
- Best for photographs with many colors
- Supports `quality` parameter
- Universal browser support

**`png`**
- PNG format, lossless compression
- Best for images with transparency or few colors
- Ignores `quality` parameter
- Larger file sizes than JPEG

**`webp`**
- Modern format with excellent compression
- Good browser support (>95%)
- Smaller than JPEG/PNG with similar quality
- Supports both lossy and lossless compression

**`avif`**
- Newest format with best compression
- ~50% smaller than JPEG with same quality
- Limited browser support (~90%)
- Requires `libaom-dev` system dependency

#### Auto Format Selection

When `format: "auto"`, the service selects format based on the client's `Accept` header:

1. If client accepts AVIF → serves AVIF
2. Else if client accepts WebP → serves WebP
3. Else → serves original format

```http
# Client sends
Accept: image/avif,image/webp,image/*,*/*;q=0.8

# Service responds with AVIF (best compression)
Content-Type: image/avif
```

---

### Fit
**Type:** String  
**Values:** `"resize"`, `"crop"`  
**Default:** `""`  
**CDN-CGI:** `fit=crop`

Controls how the image is resized when both width and height are specified.

```yaml
# Configuration
default_resize:
  fit: "crop"

# CDN-CGI
/cdn-cgi/image/width=400,height=400,fit=crop/source.jpg
```

#### Fit Options

**`Default`** (empty value)
- Proportional resizing using Lanczos algorithm
- Maintains aspect ratio
- Image fits within specified dimensions
- May not fill exact width×height if aspect ratios differ

```yaml
# Example: 1000x800 image → 400x400
fit: "resize"  # Result: 400x320 (maintains ratio)
```

**`crop`**
- Crops image to exact dimensions
- Centers the crop area
- Always produces exact width×height output
- May cut off parts of the image

```yaml
# Example: 1000x800 image → 400x400  
fit: "crop"    # Result: 400x400 (crops center)
```

#### Visual Examples

| Original | `fit: ""` | `fit: "crop"` |
|----------|-----------------|---------------|
| 1200×800 → 400×400 | 400×267 (maintains ratio) | 400×400 (crops center) |
| 800×1200 → 400×400 | 267×400 (maintains ratio) | 400×400 (crops center) |

---

### Blur
**Type:** Float  
**Range:** 0.0-10.0  
**Default:** 0 (no blur)  
**CDN-CGI:** `blur=2.5`

Applies Gaussian blur to the image with the specified radius.

```yaml
# Configuration
default_resize:
  blur: 2.5

# CDN-CGI
/cdn-cgi/image/blur=2.5/source.jpg

# URL Pattern (if supported)
/resize/blur-2.5/image.jpg
```

**Blur Guidelines:**
- `0`: No blur effect
- `0.5-1.0`: Subtle blur, slight softening
- `1.0-3.0`: Moderate blur, good for backgrounds
- `3.0-5.0`: Heavy blur, strong effect
- `5.0+`: Extreme blur

---

### Brightness
**Type:** Float  
**Range:** -100.0 to 100.0  
**Default:** 0 (no change)  
**CDN-CGI:** `brightness=20`

Adjusts image brightness. Positive values make the image brighter, negative values make it darker.

```yaml
# Configuration
default_resize:
  brightness: 15

# CDN-CGI
/cdn-cgi/image/brightness=15/source.jpg
```

**Brightness Guidelines:**
- `-100`: Completely black
- `-50`: Very dark
- `-20 to -5`: Slightly darker
- `0`: No change (default)
- `5 to 20`: Slightly brighter
- `50`: Very bright
- `100`: Completely white

---

### Contrast
**Type:** Float  
**Range:** -100.0 to 100.0  
**Default:** 0 (no change)  
**CDN-CGI:** `contrast=25`

Adjusts image contrast. Positive values increase contrast, negative values decrease it.

```yaml
# Configuration
default_resize:
  contrast: 20

# CDN-CGI
/cdn-cgi/image/contrast=20/source.jpg
```

**Contrast Guidelines:**
- `-100`: No contrast (gray)
- `-50`: Very low contrast
- `-20 to -5`: Slightly less contrast
- `0`: No change (default)
- `5 to 20`: Slightly more contrast
- `50`: High contrast
- `100`: Maximum contrast

---

### Saturation
**Type:** Float  
**Range:** -100.0 to 100.0  
**Default:** 0 (no change)  
**CDN-CGI:** `saturation=30`

Adjusts color saturation. Positive values make colors more vivid, negative values make them more muted.

```yaml
# Configuration
default_resize:
  saturation: 25

# CDN-CGI
/cdn-cgi/image/saturation=25/source.jpg
```

**Saturation Guidelines:**
- `-100`: Grayscale (no color)
- `-50`: Very muted colors
- `-20 to -5`: Slightly less saturated
- `0`: No change (default)
- `5 to 30`: More vivid colors
- `50`: Very saturated
- `100`: Maximum saturation

---

### Sharpen
**Type:** Float  
**Range:** 0.0-5.0  
**Default:** 0 (no sharpening)  
**CDN-CGI:** `sharpen=1.5`

Applies unsharp mask sharpening to enhance image details.

```yaml
# Configuration
default_resize:
  sharpen: 1.2

# CDN-CGI
/cdn-cgi/image/sharpen=1.2/source.jpg
```

**Sharpening Guidelines:**
- `0`: No sharpening
- `0.5-1.0`: Subtle sharpening, good for web images
- `1.0-2.0`: Moderate sharpening, good for photos
- `2.0-3.0`: Strong sharpening, for very soft images
- `3.0+`: Extreme sharpening (may cause artifacts)

---

### Gamma
**Type:** Float  
**Range:** 0.1-5.0  
**Default:** 0 (no correction)  
**CDN-CGI:** `gamma=1.2`

Applies gamma correction to adjust the image's luminance curve. Values > 1.0 brighten midtones, values < 1.0 darken them.

```yaml
# Configuration
default_resize:
  gamma: 1.4

# CDN-CGI
/cdn-cgi/image/gamma=1.4/source.jpg
```

**Gamma Guidelines:**
- `0.5-0.8`: Darken midtones, increase contrast
- `0.8-1.0`: Slightly darken midtones
- `1.0`: No correction (linear)
- `1.0-1.5`: Slightly brighten midtones
- `1.5-2.5`: Brighten midtones, reduce contrast
- `2.5+`: Very bright midtones

**Note:** When `gamma` is set to 0, no gamma correction is applied. Use `gamma=1.0` for linear correction.

---

### Source
**Type:** String  
**CDN-CGI:** Not applicable (part of URL)

Specifies the source file path in the storage backend. This parameter is automatically extracted from URL patterns and cannot be set manually in resize options.

```yaml
# Extracted from URL patterns
regex: '^/images/(?<source>.*)'
# URL: /images/products/photo.jpg → source: "products/photo.jpg"
```

## Usage Examples

### Basic Resizing

```yaml
# Resize to 800px width, proportional height
endpoints:
  - regex: '^/resize/(?<width>[0-9]+)/(?<source>.*)'
    default_resize:
      format: "auto"

# URL: /resize/800/product.jpg
```

### Fixed Size with Crop

```yaml
# Create square thumbnails
endpoints:
  - regex: '^/thumb/(?<source>.*)'
    default_resize:
      width: 300
      height: 300
      fit: "crop"
      format: "webp"
      quality: 85

# URL: /thumb/product.jpg → 300x300 WebP
```

### Flexible Resize Pattern

```yaml
# Support width, height, quality, and format in URL
endpoints:
  - regex: '^/img/((?<width>[0-9]+)?(x(?<height>[0-9]+))?(-(?<quality>[0-9]{1,2}))?(-(?<format>webp|avif|jpeg|png))?/)?(?<source>.*)'
    default_resize:
      format: "auto"
      quality: 85

# Examples:
# /img/product.jpg                    → auto format, quality 85
# /img/800/product.jpg                → 800px width, auto format
# /img/800x600/product.jpg            → 800x600, auto format  
# /img/800x600-95/product.jpg         → 800x600, quality 95
# /img/800x600-95-webp/product.jpg    → 800x600, quality 95, WebP
```

### Image Enhancement with Adjustments

```yaml
# Create enhanced thumbnails with sharpening and saturation boost
endpoints:
  - regex: '^/enhanced/(?<source>.*)'
    default_resize:
      width: 400
      height: 400
      fit: "crop"
      format: "webp"
      quality: 90
      sharpen: 1.2      # Enhance details
      saturation: 15    # More vivid colors
      contrast: 10      # Slight contrast boost

# URL: /enhanced/product.jpg → Enhanced 400x400 WebP
```

### Background Blur Effect

```yaml
# Create background images with blur
endpoints:
  - regex: '^/bg/(?<source>.*)'
    default_resize:
      width: 1920
      height: 1080
      fit: "crop"
      format: "webp"
      quality: 75
      blur: 3.0         # Strong blur for backgrounds
      brightness: -10   # Slightly darker

# URL: /bg/hero.jpg → Blurred background image
```

### CDN-CGI Compatible

```yaml
resize_cgi:
  enabled: true
  default_resize:
    format: "auto"
    quality: 85

# CDN-CGI URLs:
# Basic resize:
# /cdn-cgi/image/width=500/https://example.com/image.jpg

# With quality and format:
# /cdn-cgi/image/width=500,height=300,quality=95,format=webp/https://example.com/image.jpg

# With image adjustments:
# /cdn-cgi/image/width=500,sharpen=1.2,saturation=20/https://example.com/image.jpg
# /cdn-cgi/image/blur=2.5,brightness=10,contrast=15/https://example.com/image.jpg
```

## Performance Considerations

### Format Selection
- **AVIF**: Best compression (~50% smaller than JPEG) but slower encoding
- **WebP**: Good compression (~25-30% smaller than JPEG) with fast encoding
- **JPEG**: Fast encoding, good browser support
- **PNG**: Use only for images requiring transparency

### Quality Settings
- **High quality (90-100)**: Use for hero images, detailed photos
- **Medium quality (75-85)**: Recommended for most use cases
- **Low quality (60-75)**: Use for thumbnails, non-critical images

### Image Adjustments Impact
- **Blur**: Minimal performance impact, fast operation
- **Brightness/Contrast/Saturation**: Low performance impact, efficient operations
- **Sharpen**: Moderate performance impact, more computationally intensive
- **Gamma**: Low performance impact, efficient operation
- **Multiple adjustments**: Cumulative impact, applied sequentially

### Size Recommendations
- **Thumbnails**: 150×150 to 300×300 pixels
- **Gallery images**: 400×400 to 800×600 pixels
- **Hero images**: 1200×800 to 1920×1080 pixels
- **Maximum**: 4K (3840×2160) for performance reasons

### Optimization Tips
- Combine resize and adjust operations in single request rather than chaining
- Use moderate adjustment values to avoid over-processing
- Consider caching heavily-processed images
- Test adjustment combinations for optimal visual-to-performance ratio

## Error Handling

### Invalid Parameters
- Width/height > 9999: Clamped to 9999
- Quality > 100: Clamped to 100
- Quality < 1: Clamped to 1
- Invalid format: Falls back to `origin_format`
- Blur > 10.0: Clamped to 10.0
- Brightness/Contrast/Saturation > 100: Clamped to 100
- Brightness/Contrast/Saturation < -100: Clamped to -100
- Sharpen > 5.0: Clamped to 5.0
- Gamma > 5.0: Clamped to 5.0
- Gamma < 0.1: Clamped to 0.1

### Unsupported Operations
- Format conversion from/to SVG (served as-is)
- Resizing of animated GIFs (served as-is)
- Quality parameter on non-JPEG formats (ignored)

### Processing Limits
Images exceeding system memory limits will return an error. Consider implementing size restrictions in your URL patterns for production use.

```yaml
# Recommended: Limit maximum dimensions
regex: '^/resize/((?<width>[0-9]{1,4})?(x(?<height>[0-9]{1,4}))?/)?(?<source>.*)'
# Allows max 9999x9999 pixels
```
