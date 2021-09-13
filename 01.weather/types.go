/**
 * @Author: betterfor
 * @File: types
 * @Date: 2021/09/13 10:08
 * @Description:
 */

package main

// Weather 返回的天气嘻嘻
type Weather struct {
	Status    string     `json:"status"`    // 返回状态
	Count     string     `json:"count"`     // 返回结果总条数
	Info      string     `json:"info"`      // 返回的状态信息
	Infocode  string     `json:"infocode"`  // 返回状态说明
	Forecasts []Forecast `json:"forecasts"` // 预报天气信息数据
}

// Forecast 天气预报
type Forecast struct {
	City       string `json:"city"`       // 城市名称
	Adcode     string `json:"adcode"`     // 城市编码
	Province   string `json:"province"`   // 省份
	Reporttime string `json:"reporttime"` // 预报时间
	Casts      []Cast `json:"casts"`      // 预报数据
}

// Cast 详细天气
type Cast struct {
	Date         string `json:"date"`         // 日期
	Week         string `json:"week"`         // 星期
	Dayweather   string `json:"dayweather"`   // 白天天气
	Nightweather string `json:"nightweather"` // 晚上天气
	Daytemp      string `json:"daytemp"`      // 白天温度
	Nighttemp    string `json:"nighttemp"`    // 晚上温度
	Daywind      string `json:"daywind"`      // 白天风向
	Nightwind    string `json:"nightwind"`    // 晚上风向
	Daypower     string `json:"daypower"`     // 白天风力
	Nightpower   string `json:"nightpower"`   // 晚上风力
}
