package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"code.google.com/p/freetype-go/freetype"
)

var white = color.Gray{0xFF}
var black = color.Gray{0x00}
var ftContext *freetype.Context

var studentIDString string
var formIDString string
var dataString string
var outFilename string

var SID []int
var FID []int
var DATA []int

const (
	BUBBLE_WIDTH     = 5.0
	BUBBLE_HEIGHT    = 1.0
	BUBBLE_HSEP      = 2.0
	BUBBLE_Q_HSEP    = 2.66
	BUBBLE_VSEP      = 1.3
	BUBBLE_Q_VSEP    = 1.5
	FONT_FILE        = "/Library/Fonts/Arial.ttf"
	BUBBLE_FONT_SIZE = 8
	HEADING_SPACE    = 2.5
	QUESTION_SECTOR  = iota
	FORM_SECTOR
	NUMBER_SECTOR
)

func init() {
	flag.StringVar(&studentIDString, "sid", "", "populate a student ID")
	flag.StringVar(&formIDString, "fid", "", "populate the form ID")
	flag.StringVar(&dataString, "data", "", "populate responses. eg: '1,2,8,4'")
	flag.StringVar(&outFilename, "out", "out.png", "where the output png will be saved")
}

func parseFlags() {
	if studentIDString != "" {
		SID = []int{}
		for i := 0; i < len(studentIDString); i++ {
			SID = append(SID, int(studentIDString[i])-48)
		}
	}

	if formIDString != "" {
		FID = []int{}
		a, _ := strconv.Atoi(formIDString)
		FID = append(FID, a)
	}

	if dataString != "" {
		DATA = []int{}
		nums := strings.Split(dataString, ",")
		for i := 0; i < len(nums); i++ {
			DATA = append(DATA, int(nums[i][0])-48)
		}
	}
}

func drawRect(img *image.Gray, c color.Gray, xPos, yPos, width, height float64) {
	for h := 0.0; h < height; h += 1.0 {
		for w := 0.0; w < width; w += 1.0 {
			img.SetGray(int(xPos+w), int(yPos+h), c)
		}
	}
}

func drawBubble(img *image.Gray, xPos, yPos, width, height float64, content string, filled bool) {
	gapBeg := xPos + width/4.0
	gapEnd := xPos + width - width/4.0

	if filled {
		drawRect(img, black, xPos, yPos, width, height)
		return
	}

	ftContext.DrawString(content, freetype.Pt(int(gapBeg+width/6), int(yPos+height)))
	for w := xPos; w < xPos+width; w += 1.0 {
		if w >= gapBeg && w <= gapEnd {
			continue
		}
		img.SetGray(int(w), int(yPos), color.Gray{0x22})
		img.SetGray(int(w), int(yPos+height), color.Gray{0x22})
	}

	for h := yPos; h < yPos+height; h += 1.0 {
		img.SetGray(int(xPos), int(h), color.Gray{0x22})
		img.SetGray(int(xPos+width), int(h), color.Gray{0x22})
	}
}

func drawSO(img *image.Gray, unitSize, xPos, yPos float64) {
	drawRect(img, black, xPos, yPos, unitSize*3.0, unitSize*3.0)
}

func drawFP(img *image.Gray, unitSize, xPos, yPos float64) {
	// Draw outter rect 9*unitSize
	drawRect(img, black, xPos, yPos, unitSize*9.0, unitSize*9.0)
	drawRect(img, white, xPos+unitSize, yPos+unitSize, unitSize*7.0, unitSize*7.0)
	drawRect(img, black, xPos+unitSize*2.0, yPos+unitSize*2.0, unitSize*5.0, unitSize*5.0)
	drawRect(img, white, xPos+unitSize*3.0, yPos+unitSize*3.0, unitSize*3.0, unitSize*3.0)
}

func drawAP(img *image.Gray, unitSize, xPos, yPos float64) {
	drawRect(img, black, xPos, yPos, unitSize*7.0, unitSize*7.0)
	drawRect(img, white, xPos+unitSize, yPos+unitSize, unitSize*5.0, unitSize*5.0)
	drawRect(img, black, xPos+unitSize*2.0, yPos+unitSize*2.0, unitSize*3.0, unitSize*3.0)
}

func drawSector(img *image.Gray,
	unitSize, xPos, yPos float64,
	rows, cols int,
	bubbleHSep, bubbleVSep,
	leftMargin float64,
	bottomless bool,
	stype int,
	numbered bool, offset int,
	heading string,
	data []int) (float64, float64) {

	origOffset := offset

	// Compute marker width
	halfUnit := unitSize / 2.0
	qUnit := unitSize / 4.0
	markerWidth := unitSize*BUBBLE_WIDTH + unitSize*bubbleHSep - unitSize

	// Draw the sector origin
	drawSO(img, unitSize, xPos, yPos)

	// Draw the top markers
	// Draw the startbox
	x := xPos + 3.0*unitSize + unitSize
	y := yPos + unitSize
	drawRect(img, black, x, y, unitSize, unitSize)
	// Draw the h markers
	x = xPos + 4.0*unitSize + halfUnit
	y = yPos + unitSize + qUnit
	drawRect(img, black, x, y, markerWidth-leftMargin, halfUnit)
	x += markerWidth - leftMargin + unitSize
	for i := 0; i < cols-1; i++ {
		drawRect(img, black, x, y, markerWidth, halfUnit)
		x += markerWidth + unitSize
	}

	// Draw the last top marker
	drawRect(img, black, x, y, markerWidth+unitSize*2, halfUnit)

	// Draw right markers
	x += markerWidth - unitSize
	//x -= halfUnit + unitSize + qUnit
	y = yPos + 3.0*unitSize + bubbleVSep*unitSize
	//if heading != "" {
	//	y += HEADING_SPACE * unitSize
	//}
	for i := 0; i < rows; i++ {
		drawRect(img, black, x, y, unitSize*3.0, halfUnit)
		y += BUBBLE_HEIGHT*unitSize + bubbleVSep*unitSize
	}

	// Draw the bottom markers
	if !bottomless {
		x = xPos + 3.0*unitSize + unitSize
		y -= qUnit
		drawRect(img, black, x, y, unitSize, unitSize)
		x = xPos + 4.0*unitSize + halfUnit
		y += qUnit

		drawRect(img, black, x, y, markerWidth-leftMargin, halfUnit)
		x += markerWidth - leftMargin + unitSize
		for i := 0; i < cols-1; i++ {
			drawRect(img, black, x, y, markerWidth, halfUnit)
			x += markerWidth + unitSize
		}
		// Draw the last bottom marker
		drawRect(img, black, x, y, markerWidth+unitSize*2, halfUnit)
	}

	// Draw left markers
	x = xPos
	y = yPos + 3.0*unitSize + bubbleVSep*unitSize
	//if heading != "" {
	//	y += HEADING_SPACE * unitSize
	//}
	for i := 0; i < rows; i++ {
		drawRect(img, black, x, y, unitSize*3.0, halfUnit)
		y += BUBBLE_HEIGHT*unitSize + bubbleVSep*unitSize
	}
	if !bottomless {
		drawRect(img, black, x, y, unitSize*3.0, halfUnit)
		y += BUBBLE_HEIGHT*unitSize + bubbleVSep*unitSize
	}

	// Draw the heading
	if heading != "" {
		x = xPos + 5.0*unitSize
		y = yPos - 3.0*unitSize + HEADING_SPACE*unitSize
		ftContext.SetFontSize(9)
		ftContext.SetSrc(image.NewUniform(color.Gray{0x00}))
		ftContext.DrawString(heading, freetype.Pt(int(x), int(y)))
	}

	// Draw numbers
	ftContext.SetFontSize(7)
	ftContext.SetSrc(image.NewUniform(color.Gray{0x00}))
	if numbered {
		x = xPos + 3.0*unitSize + unitSize + unitSize + markerWidth + qUnit
		y = yPos + 3.0*unitSize + bubbleVSep*unitSize + qUnit
		x -= BUBBLE_WIDTH*unitSize + halfUnit
		y += BUBBLE_HEIGHT * unitSize / 1.75
		if heading != "" {
			y += HEADING_SPACE * unitSize
		}
		for j := 0; j < rows; j++ {
			ftContext.DrawString(strconv.Itoa(offset), freetype.Pt(int(x), int(y)))
			offset += 1
			y += BUBBLE_HEIGHT*unitSize + bubbleVSep*unitSize
		}
	}

	// Draw bubbles
	ftContext.SetFontSize(6.5)
	ftContext.SetSrc(image.NewUniform(color.Gray{0xBB}))
	x = xPos + 3.0*unitSize + unitSize + unitSize + markerWidth - leftMargin
	y = yPos + 3.0*unitSize + bubbleVSep*unitSize + qUnit
	x -= BUBBLE_WIDTH * unitSize / 2.0
	y -= BUBBLE_HEIGHT * unitSize / 2.0
	//if heading != "" {
	//	y += HEADING_SPACE * unitSize
	//}

	xStart := x
	for j := 0; j < rows; j++ {
		x = xStart
		for i := 0; i < cols; i++ {
			if stype == QUESTION_SECTOR {
				if data != nil && len(data) > origOffset-1+j {
					if data[j+origOffset-1] == 1<<uint(i) {
						drawBubble(img, x, y, BUBBLE_WIDTH*unitSize, BUBBLE_HEIGHT*unitSize, string(i+65), true)
					} else {
						drawBubble(img, x, y, BUBBLE_WIDTH*unitSize, BUBBLE_HEIGHT*unitSize, string(i+65), false)
					}
				} else {
					drawBubble(img, x, y, BUBBLE_WIDTH*unitSize, BUBBLE_HEIGHT*unitSize, string(i+65), false)
				}
			} else if stype == FORM_SECTOR {
				if data != nil && len(data) == 1 {
					if j*10+i == data[0] {
						drawBubble(img, x, y, BUBBLE_WIDTH*unitSize, BUBBLE_HEIGHT*unitSize, string(j*cols+i+65), true)
					} else {
						drawBubble(img, x, y, BUBBLE_WIDTH*unitSize, BUBBLE_HEIGHT*unitSize, string(j*cols+i+65), false)
					}
				} else {
					drawBubble(img, x, y, BUBBLE_WIDTH*unitSize, BUBBLE_HEIGHT*unitSize, string(j*cols+i+65), false)
				}
			} else if stype == NUMBER_SECTOR {
				if data != nil {
					if i < len(data) {
						if data[i] == j {
							drawBubble(img, x, y, BUBBLE_WIDTH*unitSize, BUBBLE_HEIGHT*unitSize, strconv.Itoa(j), true)
						} else {
							drawBubble(img, x, y, BUBBLE_WIDTH*unitSize, BUBBLE_HEIGHT*unitSize, strconv.Itoa(j), false)
						}
					} else {
						drawBubble(img, x, y, BUBBLE_WIDTH*unitSize, BUBBLE_HEIGHT*unitSize, strconv.Itoa(j), false)
					}
				} else {
					drawBubble(img, x, y, BUBBLE_WIDTH*unitSize, BUBBLE_HEIGHT*unitSize, strconv.Itoa(j), false)
				}
			} else {
				drawBubble(img, x, y, BUBBLE_WIDTH*unitSize, BUBBLE_HEIGHT*unitSize, "", false)
			}
			x += BUBBLE_WIDTH*unitSize + bubbleHSep*unitSize
		}
		y += BUBBLE_HEIGHT*unitSize + bubbleVSep*unitSize
	}
	sWidth := 3.0*unitSize + halfUnit + unitSize + markerWidth*float64(cols+1) + unitSize*float64(cols+1) + unitSize + qUnit - leftMargin
	sHeight := 3.0*unitSize + // Sector block
		bubbleVSep*unitSize + // Space below sector block
		float64(rows)*(BUBBLE_HEIGHT*unitSize+bubbleVSep*unitSize) + // Rows //NOTE
		halfUnit + qUnit // Bottom markers
	if bottomless {
		sHeight -= halfUnit + qUnit
	}
	if heading != "" {
		sHeight += HEADING_SPACE * unitSize
	}
	return sWidth, sHeight
}

func toBase4(num int, size int) []int {
	result := make([]int, size)
	var i = 0
	for num != 0 {
		r := num % 4
		num = num / 4
		result[i] = r
		i++
	}
	return result
}

func csumGen(num int) int {
	digit := 9
	sum := 0
	for digit > 0 {
		x := num % 10 // Extract the digit at the current position
		sum += x * digit
		num /= 10
		digit--
	}
	return sum % 11
}

func drawBar(img *image.Gray, unitSize, xPos, yPos float64, k int) float64 {
	switch k {
	case 0:
		drawRect(img, black, xPos, yPos, unitSize/2.0, unitSize*3.0)
		return unitSize / 2.0
	case 1:
		drawRect(img, black, xPos, yPos, unitSize, unitSize*3.0)
		return unitSize
	case 2:
		drawRect(img, black, xPos, yPos-unitSize*1.5, unitSize/2.0, unitSize*6.0)
		return unitSize / 2
	case 3:
		drawRect(img, black, xPos, yPos-unitSize*1.5, unitSize, unitSize*6.0)
		return unitSize
	default:
		log.Fatal("Bad digit given")
		return 0.0
	}
}

func drawBarcode(img *image.Gray, unitSize, xPos, yPos float64, num int) {
	drawSO(img, unitSize, xPos, yPos)
	x := xPos + unitSize*4.0

	// Draw reference bar
	x += drawBar(img, unitSize, x, yPos, 0) + unitSize

	// Print the form number below the barcode
	ftContext.SetFontSize(7)
	ftContext.SetSrc(image.NewUniform(color.Black))
	ftContext.DrawString(fmt.Sprintf("%d", num), freetype.Pt(int(x), int(yPos+unitSize*7)))

	// Draw bars
	ds := toBase4(num, 9)
	for i := len(ds) - 1; i >= 0; i-- {
		x += drawBar(img, unitSize, x, yPos, ds[i]) + unitSize
	}

	// Draw checksum bars
	ckBase10 := csumGen(num)
	fmt.Println("ckBase10 = ", ckBase10)
	ck := toBase4(ckBase10, 2)
	fmt.Println("ck = ", ck)
	for i := len(ck) - 1; i >= 0; i-- {
		x += drawBar(img, unitSize, x, yPos, ck[i]) + unitSize
	}

}

func main() {
	flag.Parse()
	parseFlags()
	dpi := 300.0
	unitSize := (1.0 / 16.0) * dpi // Unit is 1/16th of in inch

	// Page setup
	pageWidth := 8.5 * dpi
	pageHeight := 11.0 * dpi
	img := image.NewGray(image.Rect(0, 0, int(pageWidth), int(pageHeight)))
	draw.Draw(img, img.Bounds(), image.White, image.ZP, draw.Src)

	// Init freetype context
	fontBytes, err := ioutil.ReadFile(FONT_FILE)
	if err != nil {
		log.Fatal(err)
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Fatal(err)
	}
	ftContext = freetype.NewContext()
	ftContext.SetDPI(dpi)
	ftContext.SetFont(font)
	ftContext.SetClip(img.Bounds())
	ftContext.SetDst(img)
	ftContext.SetSrc(image.NewUniform(color.Gray{0xBB}))

	drawFP(img, unitSize, 0, 0)
	drawFP(img, unitSize, pageWidth-unitSize*9.0, 0)
	drawFP(img, unitSize, 0, pageHeight-unitSize*9.0)
	drawAP(img, unitSize, pageWidth-unitSize*8.0, pageHeight-unitSize*8.0)
	drawBarcode(img, unitSize, unitSize*14.0, pageHeight-unitSize*13.0, 111427)

	// Initial Offset
	yPos := unitSize * 14.0
	xPos := unitSize * 14.0

	// Name:
	ftContext.SetFontSize(11)
	ftContext.SetSrc(image.NewUniform(color.Gray{0x00}))
	ftContext.DrawString("Name:_________________________________   Student ID:______________________", freetype.Pt(int(xPos+4*unitSize), int(unitSize*5)))

	// Student ID
	w, h := drawSector(img, unitSize, xPos, yPos, 10, 10, BUBBLE_HSEP, BUBBLE_VSEP, unitSize*3.0, false, NUMBER_SECTOR, false, 0, "Student ID", SID)
	//yPos += h
	// Form ID
	w, h = drawSector(img, unitSize, xPos+w-3*unitSize, yPos, 10, 2, BUBBLE_HSEP, BUBBLE_VSEP, unitSize*3, false, FORM_SECTOR, false, 0, "Form", FID)
	yPos += h + unitSize*2

	// Questions: Col 1
	qStartHeight := yPos
	w, h = drawSector(img, unitSize, xPos, yPos, 20, 5, BUBBLE_Q_HSEP, BUBBLE_Q_VSEP, 0, true, QUESTION_SECTOR, true, 1, "", DATA)
	yPos += h
	w, h = drawSector(img, unitSize, xPos, yPos, 20, 5, BUBBLE_Q_HSEP, BUBBLE_Q_VSEP, 0, false, QUESTION_SECTOR, true, 21, "", DATA)
	//yPos += h
	//w, h = drawSector(img, unitSize, xPos, yPos, 20, 5, false)

	// Questions: Col 2
	yPos = qStartHeight
	xPos += w - unitSize*3 //+ unitSize*8
	w, h = drawSector(img, unitSize, xPos, yPos, 20, 5, BUBBLE_Q_HSEP, BUBBLE_Q_VSEP, 0, true, QUESTION_SECTOR, true, 41, "", DATA)
	yPos += h
	w, h = drawSector(img, unitSize, xPos, yPos, 20, 5, BUBBLE_Q_HSEP, BUBBLE_Q_VSEP, 0, false, QUESTION_SECTOR, true, 61, "", DATA)
	//yPos += h
	//w, h = drawSector(img, unitSize, xPos, yPos, 20, 5, false)

	// GD Logo
	logoFile, err := os.Open("gradeDeckLogo.png")
	if err != nil {
		log.Fatal(err)
	}
	logo, err := png.Decode(logoFile)
	if err != nil {
		log.Fatal(err)
	}

	draw.Draw(img, logo.Bounds().Add(image.Point{int(pageWidth - unitSize*40), int(pageHeight - unitSize*14.0)}), logo, image.ZP, draw.Src)

	imgFile, err := os.Create(outFilename)
	if err != nil {
		log.Fatal("Failed to open file. ", err)
	}
	defer imgFile.Close()
	png.Encode(imgFile, img)
}
