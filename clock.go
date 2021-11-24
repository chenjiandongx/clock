package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/nleeper/goment"
)

type Clock struct {
	date     *Date
	progress progress.Model
}

const (
	envAge        = "CLOCK_AGE"
	envExpected   = "CLOCK_EXPECTED"
	envStartColor = "CLOCK_START_COLOR"
	envEndColor   = "CLOCK_END_COLOR"
	envWho        = "CLOCK_WHO"

	defaultAge        = 25
	defaultExpected   = 90
	defaultStartColor = "#F3C8ED"
	defaultEndColor   = "#B2F6EF"
	defaultWho        = "Coder"

	padding  = 2
	maxWidth = 80
)

var (
	blueColor   = color.New(color.FgBlue, color.Italic).SprintFunc()
	cyanColor   = color.New(color.FgCyan).SprintFunc()
	yellowColor = color.New(color.FgGreen, color.Italic).SprintFunc()
	pad         = strings.Repeat(" ", 2)
)

func loadEnvInt(k string, n int) int {
	i, err := strconv.Atoi(os.Getenv(k))
	if err != nil {
		return n
	}
	return i
}

func loadEnvStr(k string, c string) string {
	s := os.Getenv(k)
	if s == "" {
		return c
	}
	return s
}

func line(text ...string) string {
	buf := &bytes.Buffer{}
	for _, t := range text {
		buf.WriteString(pad)
		buf.WriteString(t)
		buf.WriteString("\n\n")
	}
	return buf.String()
}

func New() *Clock {
	d, _ := goment.New()
	return &Clock{
		date: &Date{
			date:     d,
			age:      loadEnvInt(envAge, defaultAge),
			excepted: loadEnvInt(envExpected, defaultExpected)},
		progress: progress.NewModel(progress.WithScaledGradient(
			loadEnvStr(envStartColor, defaultStartColor),
			loadEnvStr(envEndColor, defaultEndColor)),
		),
	}
}

type Date struct {
	date     *goment.Goment
	age      int
	excepted int
}

func (d *Date) life() (string, float64) {
	n := d.excepted - d.age
	return fmt.Sprintf("你的 %s 还剩下大约 %s 年", cyanColor("人生"), blueColor(n)), float64(n) / float64(d.excepted)
}

func (d *Date) day() (string, float64) {
	h := d.date.Hour()
	if h > 0 {
		h--
	}
	m := d.date.Minute()
	escaped := float64(h*60 + m)
	return fmt.Sprintf("%s 还剩下 %s 小时 %s 分钟", cyanColor("今天"), blueColor(23-d.date.Hour()), blueColor(59-m)), (1440 - escaped) / 1440
}

func (d *Date) week() (string, float64) {
	w := d.date.Weekday()
	n := 7 - w
	return fmt.Sprintf("%s 还剩下 %s 天", cyanColor("这周"), blueColor(n)), float64(n) / 7
}

func (d *Date) month() (string, float64) {
	m := d.date.Date()
	days := d.date.DaysInMonth()
	return fmt.Sprintf("%s 还余下 %s 天", cyanColor("本月"), blueColor(days-m)), float64(days-m) / float64(days)
}

func (d *Date) year() (string, float64) {
	m := d.date.Month()
	return fmt.Sprintf("%s 年还余下 %s 个月", blueColor(d.date.Year()), blueColor(12-m)), float64(12-m) / 12
}

func (c *Clock) Who() string {
	return line(yellowColor("Hi " + loadEnvStr(envWho, defaultWho)))
}

func (c *Clock) Life() string {
	s, p := c.date.life()
	return line(s, c.progress.ViewAs(p))
}

func (c *Clock) Day() string {
	s, p := c.date.day()
	return line(s, c.progress.ViewAs(p))
}

func (c *Clock) Week() string {
	s, p := c.date.week()
	return line(s, c.progress.ViewAs(p))
}

func (c *Clock) Month() string {
	s, p := c.date.month()
	return line(s, c.progress.ViewAs(p))
}

func (c *Clock) Year() string {
	s, p := c.date.year()
	return line(s, c.progress.ViewAs(p))
}

func (c *Clock) Init() tea.Cmd { return nil }
func (c *Clock) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return c, tea.Quit

	case tea.WindowSizeMsg:
		c.progress.Width = msg.Width - padding*2 - 4
		if c.progress.Width > maxWidth {
			c.progress.Width = maxWidth
		}
		return c, nil

	default:
		return c, nil
	}
}

func (c *Clock) View() string {
	return "\n" + c.Who() + c.Life() + c.Day() + c.Week() + c.Month() + c.Year()
}

func main() {
	if err := tea.NewProgram(New()).Start(); err != nil {
		fmt.Println("Oh shit!", err)
		os.Exit(1)
	}
}
