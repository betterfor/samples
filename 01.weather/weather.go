/**
 * @Author: betterfor
 * @File: weather
 * @Date: 2021/09/13 9:55
 * @Description:
 */

package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/smtp"
	"strings"
	"time"
)

const (
	url = "https://restapi.amap.com/v3/weather/weatherInfo?"
)

var (
	key string
)

func init() {
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	key = viper.GetString("secretKey")
}

func main() {
	fmt.Println(GetWeather())
	//TimeSettle()
}

func GetWeather() (subjects, bodys []string, err error) {
	cities := viper.GetStringSlice("cityCode")
	for _, city := range cities {
		subject, body, err := getWeather(city)
		if err != nil {
			log.Fatal(err)
		}
		subjects = append(subjects, subject)
		bodys = append(bodys, body)
	}
	return
}

// 获取天气信息
func getWeather(city string) (subject, body string, err error) {
	var data Weather
	var str string
	urlInfo := url + "city=" + city + "&key=" + key + "&extensions=" + viper.GetString("weatherType")

	// 调用高德地图获取天气预报
	resp, err := http.Get(urlInfo)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	// 解析数据
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "json数据解析失败", "", err
	}
	if len(data.Forecasts) < 1 {
		return "", "", fmt.Errorf("weather: %s forecasts has no data", key)
	}
	fore := data.Forecasts[0]
	output := fore.Province + fore.City + "预报时间：" + fore.Reporttime + "\n"
	for i := 0; i < len(fore.Casts); i++ {
		cast := fore.Casts[i]
		str += fmt.Sprintf(`
日期 %s 星期 %s
白天：【天气: %s	温度: %s		风向: %s		风力: %s】
夜晚：【天气: %s	温度: %s		风向: %s		风力: %s】`,
			cast.Date, numToStr(cast.Week),
			cast.Dayweather, cast.Daytemp, cast.Daywind, cast.Daypower,
			cast.Nightweather, cast.Nighttemp, cast.Nightwind, cast.Nightpower)
	}
	subject = verify(fore.Casts[0].Dayweather, fore.Casts[0].Nightweather)
	return subject, output + str, nil
}

func verify(dayWeather, nightWeather string) string {
	var sub string
	rain := "雨"
	snow := "雪"
	sub = "今日天气预报"
	if strings.Contains(dayWeather, rain) || strings.Contains(nightWeather, rain) {
		sub = sub + "今天将降雨，出门请别忘带伞"
	}
	if strings.Contains(dayWeather, snow) || strings.Contains(nightWeather, snow) {
		sub = sub + "    下雪了"
	}
	return sub
}

// 发送邮件
func sendToMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])

	var contentType string
	if mailtype == "html" {
		contentType = "Content_Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		contentType = "Content_Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To:" + to + "\r\nFrom: " + user + "<" +
		user + ">\r\nSubject: " + subject + "\r\n" +
		contentType + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	return smtp.SendMail(host, auth, user, send_to, msg)
}

// 发送邮件
func sendEmail(subject, body string) {
	host := viper.GetString("email.host")
	user := viper.GetString("email.user")
	pwd := viper.GetString("email.password")

	to := viper.GetStringSlice("email.to")
	for _, t := range to {
		err := sendToMail(user, pwd, host, t, subject, body, "html")
		if err != nil {
			log.Fatal("Send mail error!", err)
		} else {
			fmt.Println("Send mail success!")
		}
	}
}

// TimeSettle 定时结算（一天发一次）
func TimeSettle() {
	now := time.Now()
	hour, minute := viper.GetInt("timer.hour"), viper.GetInt("timer.minute")
	fore := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
	var dur time.Duration

	if now.Hour() < hour {
		// 今天需要预报
		dur = fore.Sub(now)
	} else if now.Hour() == hour && now.Minute() < minute {
		// 今天需要预报
		dur = fore.Sub(now)
	} else {
		// 明天预报
		fore.AddDate(0, 0, 1)
		dur = fore.Sub(now)
	}
	t := time.NewTimer(dur)
	for {
		select {
		case <-t.C:
			t.Reset(time.Hour * 24)
			go func() {
				timeSettle()
			}()
		}
	}
}

func timeSettle() {
	for _, city := range viper.GetStringSlice("cityCode") {
		subject, body, err := getWeather(city)
		if err != nil {
			log.Fatal(err)
		}
		sendEmail(subject, body)
	}
}

func numToStr(str string) string {
	switch str {
	case "1":
		return "一"
	case "2":
		return "二"
	case "3":
		return "三"
	case "4":
		return "四"
	case "5":
		return "五"
	case "6":
		return "六"
	case "7":
		return "日"
	default:
		return ""
	}
}
