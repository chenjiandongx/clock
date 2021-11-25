package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbletea"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/nleeper/goment"
)

type Clock struct {
	date     *Date
	progress progress.Model
}

const (
	envBirthday   = "CLOCK_BIRTHDAY"
	envPassAway   = "CLOCK_PASS_AWAY"
	envStartColor = "CLOCK_START_COLOR"
	envEndColor   = "CLOCK_END_COLOR"
	envWho        = "CLOCK_WHO"

	timeLayout        = "2006-01-02"
	defaultBirthday   = "1996-04-12"
	defaultPassAway   = "2086-04-12"
	defaultStartColor = "#F3C8ED"
	defaultEndColor   = "#B2F6EF"
	defaultWho        = "Coder"

	padding  = 2
	maxWidth = 80
)

var (
	cyanItalic = color.New(color.FgCyan, color.Italic).SprintFunc()
	cyanColor  = color.New(color.FgCyan).SprintFunc()
	greenColor = color.New(color.FgGreen, color.Italic).SprintFunc()
	pad        = strings.Repeat(" ", 2)
)

func loadEnvStr(k, a string) string {
	s := os.Getenv(k)
	if s == "" {
		return a
	}
	return s
}

func loadEnvDate(k, a string) time.Time {
	s := os.Getenv(k)
	t, err := time.Parse(timeLayout, s)
	if err != nil {
		t, _ = time.Parse(timeLayout, a)
		return t
	}
	return t
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
	g, _ := goment.New()
	return &Clock{
		date: &Date{
			date:     g,
			birthday: loadEnvDate(envBirthday, defaultBirthday),
			passAway: loadEnvDate(envPassAway, defaultPassAway),
		},
		progress: progress.NewModel(progress.WithScaledGradient(
			loadEnvStr(envStartColor, defaultStartColor),
			loadEnvStr(envEndColor, defaultEndColor)),
		),
	}
}

type Date struct {
	date     *goment.Goment
	birthday time.Time
	passAway time.Time
}

func (d *Date) life() (string, float64) {
	now := time.Now()
	y := d.passAway.Year() - now.Year()
	m := humanize.Comma((now.Unix() - d.birthday.Unix()) / 3600)
	s := humanize.Comma(d.passAway.Unix() - now.Unix())
	percent := float64(d.passAway.Unix()-now.Unix()) / float64(d.passAway.Unix())
	return fmt.Sprintf("你的 %s 还剩下大约 %s 年 已经走过 %s 分钟 距离终点还有 %s 秒", cyanColor("人生"), cyanItalic(y), cyanItalic(m), cyanItalic(s)), percent
}

func (d *Date) day() (string, float64) {
	now := time.Now()
	h := now.Hour()
	if h > 0 {
		h--
	}
	m := now.Minute()
	escaped := float64(h*60 + m)
	s := 59 - (time.Now().Unix() % 60)
	percent := (1440 - escaped) / 1440
	return fmt.Sprintf("%s 还剩下 %s 小时 %s 分钟 %s 秒", cyanColor("今天"), cyanItalic(23-d.date.Hour()), cyanItalic(59-m), cyanItalic(s)), percent
}

func (d *Date) week() (string, float64) {
	w := d.date.Weekday()
	n := 7 - w
	percent := float64(n) / 7
	return fmt.Sprintf("%s 还剩下 %s 天", cyanColor("这周"), cyanItalic(n)), percent
}

func (d *Date) month() (string, float64) {
	m := d.date.Date()
	days := d.date.DaysInMonth()
	percent := float64(days-m) / float64(days)
	return fmt.Sprintf("%s 还余下 %s 天", cyanColor("本月"), cyanItalic(days-m)), percent
}

func (d *Date) year() (string, float64) {
	m := d.date.Month()
	percent := float64(12-m) / 12
	return fmt.Sprintf("%s 年还余下 %s 个月", cyanItalic(d.date.Year()), cyanItalic(12-m)), percent
}

func (c *Clock) Who() string {
	return line(greenColor("Hi " + loadEnvStr(envWho, defaultWho)))
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

func (c *Clock) Sigh() string {
	return line("来都来了 给个面子 都不容易 大过年的 习惯就好")
}

func (c *Clock) Init() tea.Cmd { return tickCmd() }
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
		return c, tickCmd()
	}
}

func (c *Clock) View() string {
	return "\n" + c.Who() + c.Life() + c.Day() + c.Week() + c.Month() + c.Year() + c.Sigh()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return t
	})
}

func main() {
	if err := tea.NewProgram(New()).Start(); err != nil {
		fmt.Println("Oh shit!", err)
		os.Exit(1)
	}
}
