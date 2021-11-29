package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/fogleman/ease"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"
	"github.com/nleeper/goment"
	"github.com/spf13/viper"
)

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

	progressFullChar  = "█"
	progressEmptyChar = "░"

	maxWidth = 80
)

var (
	cyanItalic = color.New(color.FgCyan, color.Italic).SprintFunc()
	cyanColor  = color.New(color.FgCyan).SprintFunc()
	greenColor = color.New(color.FgGreen, color.Italic).SprintFunc()
	pad        = strings.Repeat(" ", 2)

	term          = termenv.ColorProfile()
	progressEmpty = termenv.Style{}.Foreground(term.Color("241")).Styled(progressEmptyChar)
	ramp          []string
	speedup       int
	ch            = make(chan struct{}, 1)
)

func loadEnvStr(k, a string) string {
	s := os.Getenv(k)
	if s != "" {
		return s
	}

	s = viper.GetString(k)
	if s == "" {
		return a
	}
	return s
}

func loadEnvDate(k, a string) time.Time {
	s := os.Getenv(k)
	t, err := time.Parse(timeLayout, s)
	if err == nil {
		return t
	}

	t, err = time.Parse(timeLayout, viper.GetString(k))
	if err != nil {
		t, _ = time.Parse(timeLayout, a)
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

func getSpeedup() int {
	return speedup + (speedup/2+speedup+int(math.Log2(float64(speedup))))*speedup
}

type Clock struct {
	date   *Date
	frames []int
	loads  []bool
}

const (
	ProgLife = iota
	ProgWork
	ProgDay
	ProgWeek
	ProgMonth
	ProgYear
)

var globalNow = &goment.Goment{}

func updateGlobalNow() {
	globalNow, _ = goment.New()
	globalNow.Add(time.Hour * time.Duration(getSpeedup()))
}

func Now() *goment.Goment {
	return globalNow
}

func New() *Clock {
	return &Clock{
		date: &Date{
			birthday: loadEnvDate(envBirthday, defaultBirthday),
			passAway: loadEnvDate(envPassAway, defaultPassAway),
		},
		frames: make([]int, 6),
		loads:  make([]bool, 6),
	}
}

type Date struct {
	birthday time.Time
	passAway time.Time
}

func (d *Date) life() (string, float64) {
	now := Now()
	y := d.passAway.Year() - now.Year()
	m := humanize.Comma((now.ToUnix() - d.birthday.Unix()) / 3600)
	s := humanize.Comma(d.passAway.Unix() - now.ToUnix())
	percent := float64(d.passAway.Unix()-now.ToUnix()) / float64(d.passAway.Unix()-d.birthday.Unix())
	return fmt.Sprintf("你的 %s 还剩下 %s 年 已经走过 %s 小时 距离终点还有 %s 秒", cyanColor("人生"), cyanItalic(y), cyanItalic(m), cyanItalic(s)), percent
}

func (d *Date) work() (string, float64) {
	blessings := 35
	for {
		if d.birthday.AddDate(blessings, 0, 0).Unix()-Now().ToUnix() > 0 {
			return d.warning(blessings)
		}
		blessings += 5
	}
}

func (d *Date) warning(age int) (string, float64) {
	oh := d.birthday.AddDate(age, 0, 0).Unix()
	n := oh - Now().ToUnix()
	m := n / 60
	s := n - m*60

	percent := float64(n) / float64(oh-d.birthday.Unix())
	return fmt.Sprintf("%s 距离你 %s 岁生日还有 %s 分钟 %s 秒", cyanColor("别紧张"), cyanItalic(age), cyanItalic(humanize.Comma(m)), cyanItalic(s)), percent
}

func (d *Date) day() (string, float64) {
	now := Now()
	h := now.Hour()
	if h > 0 {
		h--
	}
	m := now.Minute()
	escaped := float64(h*60 + m)
	s := 60 - (Now().ToUnix() % 60)
	percent := (1440 - escaped) / 1440
	return fmt.Sprintf("%s 还剩下 %s 小时 %s 分钟 %s 秒", cyanColor("今天"), cyanItalic(23-now.Hour()), cyanItalic(59-m), cyanItalic(s)), percent
}

func (d *Date) week() (string, float64) {
	w := Now().Weekday()
	n := 7 - w
	percent := float64(n) / 7
	return fmt.Sprintf("%s 还剩下 %s 天", cyanColor("这周"), cyanItalic(n)), percent
}

func (d *Date) month() (string, float64) {
	m := Now().Date()
	days := Now().DaysInMonth()
	percent := float64(days-m) / float64(days)
	return fmt.Sprintf("%s 还余下 %s 天", cyanColor("本月"), cyanItalic(days-m)), percent
}

func (d *Date) year() (string, float64) {
	m := Now().Month()
	percent := float64(12-m) / 12
	return fmt.Sprintf("%s 年还余下 %s 个月", cyanItalic(Now().Year()), cyanItalic(12-m)), percent
}

func (c *Clock) Who() string {
	return line(fmt.Sprintf("%s 现在是 %s", greenColor("Hi "+loadEnvStr(envWho, defaultWho)), cyanItalic(Now().Format("YYYY-MM-DD HH:mm:ss"))))
}

func (c *Clock) Sigh() string {
	return line("来都来了 给个面子 就这样吧 都不容易 是个孩子 大过年的 都是朋友 习惯就好")
}

func (c *Clock) Stop() string {
	return line("你已经走完了生命的全部旅程 你无法超越时间", "珍惜当下吧 这个世界没有什么好畏惧的 反正我们只来一次")
}

func (c *Clock) Help() string {
	text := "长按 <s> 为生命加速；按 <p> 回到当下；按 <q> 退出"
	return line(termenv.Style{}.Foreground(term.Color("241")).Styled(text))
}

func (c *Clock) Loaded() bool {
	for _, b := range c.loads {
		if !b {
			return false
		}
	}
	return true
}

func (c *Clock) render(s string, p float64, index int, f func(float64) float64) string {
	if !c.loads[index] {
		c.frames[index]++
		val := f(float64(c.frames[index]) / float64(100))
		if val >= p {
			val = p
			c.loads[index] = true
		}
		return line(s, c.progressbar(val))
	}
	return line(s, c.progressbar(p))
}

func (c *Clock) Life() string {
	s, p := c.date.life()
	return c.render(s, p, ProgLife, ease.InOutQuad)
}

func (c *Clock) Work() string {
	s, p := c.date.work()
	return c.render(s, p, ProgWork, ease.InOutCubic)
}

func (c *Clock) Day() string {
	s, p := c.date.day()
	return c.render(s, p, ProgDay, ease.InOutQuart)
}

func (c *Clock) Week() string {
	s, p := c.date.week()
	return c.render(s, p, ProgWeek, ease.InOutSine)
}

func (c *Clock) Month() string {
	s, p := c.date.month()
	return c.render(s, p, ProgMonth, ease.InOutCirc)
}

func (c *Clock) Year() string {
	s, p := c.date.year()
	return c.render(s, p, ProgYear, ease.InOutBounce)
}

func (c *Clock) Init() tea.Cmd { return tickFastCmd() }
func (c *Clock) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "q", "Q":
			return c, tea.Quit
		case "p", "P":
			speedup = 0
			updateGlobalNow()
		case "s", "S":
			speedup++
			updateGlobalNow()
		}
	}

	if c.Loaded() {
		return c, tickSlowCmd()
	}
	return c, tickFastCmd()
}

func (c *Clock) View() string {
	if Now().ToUnix() >= c.date.passAway.Unix() {
		globalNow, _ = goment.New(c.date.passAway)
		return "\n" + c.Who() + c.Stop() + c.Help()
	}
	return "\n" + c.Who() + c.Life() + c.Work() + c.Day() + c.Week() + c.Month() + c.Year() + c.Sigh() + c.Help()
}

func (c *Clock) progressbar(percent float64) string {
	w := float64(maxWidth)

	fullSize := int(math.Round(w * percent))
	var fullCells string
	for i := 0; i < fullSize; i++ {
		fullCells += termenv.String(progressFullChar).Foreground(term.Color(ramp[i])).String()
	}

	emptySize := int(w) - fullSize
	emptyCells := strings.Repeat(progressEmpty, emptySize)

	return fmt.Sprintf("%s%s %3.0f%%", fullCells, emptyCells, math.Round(percent*100))
}

type tickFastMsg struct{}

func tickFastCmd() tea.Cmd {
	return tea.Tick(time.Second/15, func(t time.Time) tea.Msg {
		return tickFastMsg{}
	})
}

type tickSlowMsg struct{}

func tickSlowCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickSlowMsg{}
	})
}

func makeRamp(colorA, colorB string, steps float64) (s []string) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0.0; i < steps; i++ {
		c := cA.BlendLuv(cB, i/steps)
		s = append(s, colorToHex(c))
	}
	return
}

func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}

func main() {
	viper.SetConfigName(".life_clock")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath("/etc")
	_ = viper.ReadInConfig()

	defer func() {
		ch <- struct{}{}
	}()

	ticker := time.Tick(time.Second)
	go func() {
		for {
			select {
			case <-ticker:
				updateGlobalNow()
			case <-ch:
				return
			}
		}
	}()

	startColor := loadEnvStr(envStartColor, defaultStartColor)
	endColor := loadEnvStr(envEndColor, defaultEndColor)
	ramp = makeRamp(startColor, endColor, maxWidth)
	updateGlobalNow()

	if err := tea.NewProgram(New()).Start(); err != nil {
		fmt.Println("Oh shit!", err)
		os.Exit(1)
	}
}
