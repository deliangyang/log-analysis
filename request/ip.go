package request

import (
	"net/http"
	"errors"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

var host = "https://dm-81.data.aliyun.com/rest/160601/ip/getIpInfo.json"
var appCode = "40489272c9a24ba3b483f8ecc37cb846"

type Location struct {
	ip string
	country string
	area string
	region string
	city string
	county string
	isp string
	countryId string
	areaId string
	regionId string
	cityId string
	countyId string
	ispId string
}

type Response struct {
	Code int
	Data Location
}

func (local *Location) Query(ip string) (response Response, err error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", host + "?ip=" + ip, nil)
	if err != nil {
		return response, errors.New("fail")
	}

	request.Header.Set("Authorization", "APPCODE " + appCode)
	resp, _ := client.Do(request)
	defer resp.Body.Close()

	fmt.Println(resp.Body)
	fmt.Println(resp.StatusCode)
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		//response := Response{}
		json.Unmarshal(body, &response)
		fmt.Println(response)
	}
	return response, nil
}
