package loglist

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
)

type Entry struct {
	Time    time.Time
	Level   logrus.Level
	Message string
	Data    logrus.Fields
}

const dateFormat = "[2006-01-02 15:04:05 MST]"

// TODO allow limiting log length by rotating entries off the front?

type ListHook struct {
	LogData []Entry
	levels  []logrus.Level
}

func (p *ListHook) Levels() []logrus.Level {
	return p.levels
}

func (p *ListHook) Fire(e *logrus.Entry) error {
	p.LogData = append(p.LogData, Entry{
		Time:    e.Time,
		Level:   e.Level,
		Message: e.Message,
		Data:    e.Data,
	})
	return nil
}

func (p *ListHook) List() *widget.List {
	list := widget.NewList(
		func() int { return len(p.LogData) },
		func() fyne.CanvasObject {
			// show an icon for level
			icon := widget.NewIcon(theme.InfoIcon())

			// timestamp
			ts := canvas.NewText(dateFormat, theme.ForegroundColor())
			ts.TextStyle.Monospace = true

			// level label
			level := canvas.NewText("[unknown]", theme.ForegroundColor())
			level.TextStyle.Bold = true
			level.TextStyle.Monospace = true

			// message
			msg := widget.NewLabel("messsage")

			return container.NewHBox(icon, ts, level, msg)
			// return canvas.NewText()
		},
		func(i int, co fyne.CanvasObject) {
			e := p.LogData[i]

			row := co.(*fyne.Container)

			icon := row.Objects[0].(*widget.Icon)
			ts := row.Objects[1].(*canvas.Text)
			level := row.Objects[2].(*canvas.Text)
			msg := row.Objects[3].(*widget.Label)

			ts.Text = e.Time.Format(dateFormat)
			level.Text = fmt.Sprintf("[%s]", e.Level)
			msg.SetText(e.Message)

			switch e.Level {
			case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
				icon.SetResource(theme.ErrorIcon())
				level.Color = theme.ErrorColor()
			case logrus.WarnLevel:
				icon.SetResource(theme.WarningIcon())
				level.Color = theme.PrimaryColorNamed(theme.ColorOrange)
			case logrus.DebugLevel:
				icon.SetResource(theme.QuestionIcon())
				level.Color = theme.PrimaryColorNamed(theme.ColorGreen)
			default:
				icon.SetResource(theme.InfoIcon())
				level.Color = theme.ForegroundColor()
			}

			co.Refresh()
		})
	return list
}

func NewListHook(levels ...logrus.Level) *ListHook {
	return &ListHook{
		LogData: []Entry{},
		levels:  levels,
	}
}
