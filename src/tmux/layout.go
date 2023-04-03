package tmux

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/mdcli/cmd"
)

var (
	checksumPatternString = `[a-f0-9]{4}`
	checksumPattern       = regexp.MustCompile(checksumPatternString)

	dimensionPatternString = `(?P<width>\d+)x(?P<height>\d+),\d+,\d+`
	dimensionPattern       = regexp.MustCompile(dimensionPatternString)

	ignoreDimensionPatternString = `\d+x\d+,\d+,\d+`
	ignoreDimensionPattern       = regexp.MustCompile(ignoreDimensionPatternString)

	sidePanesPatternString = `\[(?P<sidePanes>.*)\]`
	sidePanesPattern       = regexp.MustCompile(sidePanesPatternString)

	sidePanePatternString = `\d+x\d+,\d+,\d+,(?P<sidePaneId>\d+)`
	sidePanePattern       = regexp.MustCompile(sidePanePatternString)

	singlePanePattern = regexp.MustCompile(fmt.Sprintf("^%v,%v,%v", checksumPatternString, dimensionPatternString, `(?P<mainPaneId>\d+)`))
	twoPanePattern    = regexp.MustCompile(fmt.Sprintf("^%v,%v{%v,%v,%v,%v}", checksumPatternString, dimensionPatternString, ignoreDimensionPatternString, `(?P<mainPaneId>\d+)`, ignoreDimensionPatternString, `(?P<sidePaneId>\d+)`))
	multiPanePattern  = regexp.MustCompile(fmt.Sprintf("^%v,%v{%v,%v,%v%v}", checksumPatternString, dimensionPatternString, ignoreDimensionPatternString, `(?P<mainPaneId>\d+)`, ignoreDimensionPatternString, sidePanesPatternString))
)

type tmuxLayout struct {
	csum   string
	layout string
}

type baseLayout struct {
	mainPaneId  string
	sidePaneIds []string
	width       int
	height      int
}

func (l *tmuxLayout) String() string {
	return fmt.Sprintf("%v,%v", l.csum, l.layout)
}

func layoutChecksum(layout string) string {
	csum := 0
	for _, c := range layout {
		csum = (csum >> 1) + ((csum & 1) << 15)
		csum += int(c)
	}

	return fmt.Sprintf("%04x", csum)
}

func newTmuxLayout(layout string) *tmuxLayout {
	return &tmuxLayout{
		csum:   layoutChecksum(layout),
		layout: layout,
	}
}

func parseLayout(layout *tmuxLayout) (*baseLayout, error) {
	if singlePanePattern.Match([]byte(layout.String())) {
		matches := singlePanePattern.FindStringSubmatch(layout.String())

		mainPaneId := matches[singlePanePattern.SubexpIndex("mainPaneId")]

		width, err := strconv.Atoi(matches[singlePanePattern.SubexpIndex("width")])
		if err != nil {
			return &baseLayout{}, err
		}

		height, err := strconv.Atoi(matches[singlePanePattern.SubexpIndex("height")])
		if err != nil {
			return &baseLayout{}, err
		}

		return &baseLayout{
			mainPaneId:  mainPaneId,
			sidePaneIds: []string{},
			width:       width,
			height:      height,
		}, nil
	} else if twoPanePattern.Match([]byte(layout.String())) {
		matches := twoPanePattern.FindStringSubmatch(layout.String())

		mainPaneId := matches[twoPanePattern.SubexpIndex("mainPaneId")]
		sidePaneId := matches[twoPanePattern.SubexpIndex("sidePaneId")]

		width, err := strconv.Atoi(matches[twoPanePattern.SubexpIndex("width")])
		if err != nil {
			return &baseLayout{}, err
		}

		height, err := strconv.Atoi(matches[twoPanePattern.SubexpIndex("height")])
		if err != nil {
			return &baseLayout{}, err
		}

		return &baseLayout{
			mainPaneId:  mainPaneId,
			sidePaneIds: []string{sidePaneId},
			width:       width,
			height:      height,
		}, nil
	} else if multiPanePattern.Match([]byte(layout.String())) {
		matches := multiPanePattern.FindStringSubmatch(layout.String())

		mainPaneId := matches[multiPanePattern.SubexpIndex("mainPaneId")]

		width, err := strconv.Atoi(matches[multiPanePattern.SubexpIndex("width")])
		if err != nil {
			return &baseLayout{}, err
		}

		height, err := strconv.Atoi(matches[multiPanePattern.SubexpIndex("height")])
		if err != nil {
			return &baseLayout{}, err
		}

		sidePanes := matches[multiPanePattern.SubexpIndex("sidePanes")]
		sidePaneIds := make([]string, 0)
		for _, match := range sidePanePattern.FindAllStringSubmatch(sidePanes, -1) {
			sidePaneId := match[sidePanePattern.SubexpIndex("sidePaneId")]
			sidePaneIds = append(sidePaneIds, sidePaneId)
		}

		return &baseLayout{
			mainPaneId:  mainPaneId,
			sidePaneIds: sidePaneIds,
			width:       width,
			height:      height,
		}, nil
	}

	return &baseLayout{}, errors.New("Could not parse layout")
}

func defaultLayout(layout *baseLayout) (*tmuxLayout, error) {
	mainWidth := int(math.Round(0.65 * float64(layout.width)))
	sideWidth := layout.width - mainWidth

	if len(layout.sidePaneIds) == 0 {
		return newTmuxLayout(fmt.Sprintf("%vx%v,0,0,%v", layout.width, layout.height, layout.mainPaneId)), nil
	} else if len(layout.sidePaneIds) == 1 {
		mainLayout := fmt.Sprintf("%vx%v,0,0,%v", mainWidth, layout.height, layout.mainPaneId)
		sideLayout := fmt.Sprintf("%vx%v,%v,0,%v", sideWidth, layout.height, mainWidth+1, layout.mainPaneId)

		return newTmuxLayout(fmt.Sprintf("%vx%v,0,0{%v,%v}", layout.width, layout.height, mainLayout, sideLayout)), nil
	} else {
		mainLayout := fmt.Sprintf("%vx%v,0,0,%v", mainWidth, layout.height, layout.mainPaneId)

		sideLayout := fmt.Sprintf("%vx%v,%v,0", layout.width-mainWidth, layout.height, mainWidth+1)
		sideLayouts := make([]string, 0)
		numSidePanes := len(layout.sidePaneIds)
		for i, paneId := range layout.sidePaneIds {
			sideStart := 0
			if i != 0 {
				sideStart = int(math.Round(float64(i)*float64(layout.height)/float64(numSidePanes))) + 1
			}

			sideEnd := int(math.Round(float64(1+i) * float64(layout.height) / float64(numSidePanes)))
			sideHeight := sideEnd - sideStart

			sideLayout := fmt.Sprintf("%vx%v,%v,%v,%v", sideWidth, sideHeight, mainWidth+1, sideStart, paneId)
			sideLayouts = append(sideLayouts, sideLayout)
		}

		return newTmuxLayout(fmt.Sprintf(
			"%vx%v,0,0{%v,%v[%v]}",
			layout.width, layout.height,
			mainLayout,
			sideLayout, strings.Join(sideLayouts, ","),
		)), nil
	}
}

func getLayout(session string, window string) (*tmuxLayout, error) {
	output, err := cmd.CaptureCommand(
		"tmux", "display-message",
		"-t", fmt.Sprintf("%v:%v", session, window),
		"-p", "#{window_layout}",
	)
	if err != nil {
		return &tmuxLayout{}, err
	}

	layoutStr := strings.TrimSpace(output)
	csum, layout, _ := strings.Cut(layoutStr, ",")

	return &tmuxLayout{
		csum:   csum,
		layout: layout,
	}, nil
}

func setLayout(session string, window string, layout *tmuxLayout) error {
	layoutString := fmt.Sprintf("%v,%v", layout.csum, layout.layout)
	err := cmd.RunCommand(
		"tmux", "select-layout",
		"-t", fmt.Sprintf("%v:%v", session, window),
		layoutString,
	)

	return err
}

func getDefaultLayout(session string, window string) (*tmuxLayout, error) {
	layout, err := getLayout(session, window)
	if err != nil {
		return &tmuxLayout{}, err
	}

	parsedLayout, err := parseLayout(layout)
	if err != nil {
		return &tmuxLayout{}, err
	}

	newLayout, err := defaultLayout(parsedLayout)
	if err != nil {
		return &tmuxLayout{}, err
	}

	return newLayout, nil
}

func setDefaultLayout(session string, window string) error {
	layout, err := getDefaultLayout(session, window)
	if err != nil {
		return err
	}

	return setLayout(session, window, layout)
}
