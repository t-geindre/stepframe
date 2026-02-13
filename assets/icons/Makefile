SVG_DIR := src
PNG_DIR := pngs
SHEET   := ../../ui/theme/icons.png

SIZE    := 64
COLS    := 8

SVGS      := $(wildcard $(SVG_DIR)/*.svg)
RAW_PNGS  := $(wildcard $(SVG_DIR)/*.png)

# PNGs generated from SVGs
PNGS_FROM_SVG := $(patsubst $(SVG_DIR)/%.svg,$(PNG_DIR)/%.png,$(SVGS))
# PNGs copied as-is from src/
PNGS_FROM_RAW := $(patsubst $(SVG_DIR)/%.png,$(PNG_DIR)/%.png,$(RAW_PNGS))

ALL_PNGS := $(PNGS_FROM_SVG) $(PNGS_FROM_RAW)

.PHONY: all icons raw clean

all: raw icons clean

# Copy raw PNGs from src/ to pngs/
raw:
	@mkdir -p $(PNG_DIR)
	@if ls $(SVG_DIR)/*.png >/dev/null 2>&1; then \
		cp $(SVG_DIR)/*.png $(PNG_DIR)/ ; \
	fi

# Build sprite sheet from both converted SVGs and copied PNGs
icons: $(ALL_PNGS)
	@mkdir -p $(dir $(SHEET))
	montage $(ALL_PNGS) -background none -tile $(COLS)x -geometry +0+0 $(SHEET)

# Convert SVG -> PNG
$(PNG_DIR)/%.png: $(SVG_DIR)/%.svg
	@mkdir -p $(PNG_DIR)
	inkscape $< -w $(SIZE) -h $(SIZE) --export-type=png --export-filename=$@

clean:
	rm -rf $(PNG_DIR)
