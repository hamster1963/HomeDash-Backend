package ha_utils

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"home-network-watcher/manifest"
	"time"
)

type uHaUtils struct{}

var HaUtils = &uHaUtils{}

// 卧室空调
const (
	AirConditionerEntityPattern = "climate.lumi_mcn02_1f58_air_conditioner"
	AirConditionerStatePattern  = "attributes.air_conditioner.on"
	AirConditionerTempPattern   = "attributes.temperature"
)

// 加湿器
const (
	HumidifierEntityPattern   = "humidifier.deerma_jsq2g_d916_humidifier"
	HumidifierStatePattern    = "attributes.humidifier.on"
	HumidifierHumidityPattern = "attributes.environment.relative_humidity"
)

// 空气净化器
const (
	AirPurifierEntityPattern = "fan.zhimi_mp4a_20cb_air_purifier"
	AirPurifierStatePattern  = "attributes.air_purifier.on"
	AirPurifierPM25Pattern   = "attributes.environment.pm2_5_density"
)

// 卧室床头灯
const (
	LightEntityPattern     = "light.yeelink_bslamp2_0b02_light"
	LightStatePattern      = "attributes.light.on"
	LightBrightnessPattern = "attributes.light.brightness"
)

type HaEntitiesInfo struct {
	AirConditioner g.Map
	Humidifier     g.Map
	AirPurifier    g.Map
	Light          g.Map
}

// GetEntitiesInfo
//
//	@dc: 获取实体信息
//	@author: hamster   @date:2023/9/22 17:54:25
func (u *uHaUtils) GetEntitiesInfo(ctx context.Context) (err error) {
	haClient := g.Client().SetHeaderMap(manifest.HomeAssistantAuthMap)
	baseUrl := manifest.HomeAssistantBaseUrl
	// 获取空调状态
	airConditionerMap := g.Map{}
	response, err := haClient.Get(ctx, baseUrl+AirConditionerEntityPattern)
	if err != nil {
		return err
	}
	defer response.Close()
	responseJson := gjson.New(response.ReadAllString())
	responseJson.SetViolenceCheck(true)
	airConditionerMap["state"] = responseJson.Get(AirConditionerStatePattern)
	airConditionerMap["temp"] = responseJson.Get(AirConditionerTempPattern)

	// 获取加湿器状态
	humidifierMap := g.Map{}
	response, err = haClient.Get(ctx, baseUrl+HumidifierEntityPattern)
	if err != nil {
		return err
	}
	defer response.Close()
	responseJson = gjson.New(response.ReadAllString())
	responseJson.SetViolenceCheck(true)
	humidifierMap["state"] = responseJson.Get(HumidifierStatePattern)
	humidifierMap["humidity"] = responseJson.Get(HumidifierHumidityPattern)

	// 获取空气净化器状态
	airPurifierMap := g.Map{}
	response, err = haClient.Get(ctx, baseUrl+AirPurifierEntityPattern)
	if err != nil {
		return err
	}
	defer response.Close()
	responseJson = gjson.New(response.ReadAllString())
	responseJson.SetViolenceCheck(true)
	airPurifierMap["state"] = responseJson.Get(AirPurifierStatePattern)
	airPurifierMap["pm25"] = responseJson.Get(AirPurifierPM25Pattern)

	// 获取床头灯状态
	lightMap := g.Map{}
	response, err = haClient.Get(ctx, baseUrl+LightEntityPattern)
	if err != nil {
		return err
	}
	defer response.Close()
	responseJson = gjson.New(response.ReadAllString())
	responseJson.SetViolenceCheck(true)
	lightMap["state"] = responseJson.Get(LightStatePattern)
	lightMap["brightness"] = responseJson.Get(LightBrightnessPattern)

	// 汇总信息
	HaEntitiesInfo := &HaEntitiesInfo{
		AirConditioner: airConditionerMap,
		Humidifier:     humidifierMap,
		AirPurifier:    airPurifierMap,
		Light:          lightMap,
	}
	err = gcache.Set(ctx, manifest.HaEntitiesCacheKey, HaEntitiesInfo, 1*time.Minute)
	if err != nil {
		return err
	}

	return
}
