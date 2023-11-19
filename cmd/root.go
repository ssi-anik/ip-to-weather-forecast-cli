package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"github.com/ssi-anik/ip-to-weather-forecast-cli/types"
)

func init() {
	rootCmd.Flags().StringVarP(&ip, "ip", "i", "", "For the IP")
	// rootCmd.MarkFlagRequired("ip")
}

var ip string
var now = time.Now()

var rootCmd = &cobra.Command{
	Use:   "ip-to-weather-forecast-cli",
	Short: "How's the weather?",
	Long:  `Show weather of your area or of an IP`,
	Run:   run,
}

func getGeoLocation() (*types.IpGeoLocation, error) {
	url := `https://json.geoiplookup.io/`
	if ip != "" {
		addr, err := netip.ParseAddr(ip)
		if err != nil {
			panic(err)
		}
		url = fmt.Sprintf("%s?ip=%s", url, addr.String())
	}

	res, err := http.Get(url)
	if nil != err {
		panic(fmt.Errorf("cannot get geolocation. %s", err))
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		panic(fmt.Errorf("response status code is not 200"))
	}

	c, err := io.ReadAll(res.Body)
	if nil != err {
		panic(fmt.Errorf("cannot read response body. %s", err))
	}

	geo := types.NewIpGeoLocation()
	err = json.Unmarshal(c, geo)
	if nil != err {
		panic(fmt.Errorf("cannot unmarshal response. %s", err))
	}

	return geo, nil
}

func getForecast(lat, lng float64) *types.Forecast {
	url := fmt.Sprintf(
		`https://api.brightsky.dev/weather?lat=%f&lon=%f&date=%s`,
		lat,
		lng,
		time.Now().Format("2006-01-02"),
	)

	res, err := http.Get(url)
	if nil != err {
		panic(fmt.Errorf("cannot get geolocation. %s", err))
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		panic(fmt.Errorf("response status code is not 200"))
	}

	c, err := io.ReadAll(res.Body)
	if nil != err {
		panic(fmt.Errorf("cannot read response body. %s", err))
	}

	forecast := types.NewForecast()
	err = json.Unmarshal(c, forecast)
	if nil != err {
		panic(fmt.Errorf("cannot unmarshal response. %s", err))
	}

	return forecast
}

func printGeo(geo *types.IpGeoLocation) {
	tbl := table.NewWriter()
	tbl.SetTitle("Geo information")
	tbl.SetStyle(table.StyleColoredCyanWhiteOnBlack)
	tbl.SetOutputMirror(os.Stdout)
	tbl.Style().Title.Align = text.AlignCenter
	tbl.AppendRows([]table.Row{
		{"IP", geo.Ip},
		{"Latitude", geo.Latitude},
		{"Longitude", geo.Longitude},
		{"Country", fmt.Sprintf("%v (%v)", geo.CountryName, geo.CountryCode)},
		{"Currency", fmt.Sprintf("%v (%v)", geo.CurrencyName, geo.CurrencyCode)},
		{"Timezone", geo.TimezoneName},
		{"Connection Type", geo.ConnectionType},
		{"Date", now.Format("2006-02-01")},
	})
	tbl.Render()
}

func printForecast(forecast *types.Forecast, timezone string) {
	tbl := table.NewWriter()
	tbl.SetTitle("Weather forecast")
	tbl.SetStyle(table.StyleColoredYellowWhiteOnBlack)
	tbl.Style().Title.Align = text.AlignCenter
	tbl.SetOutputMirror(os.Stdout)
	tbl.AppendHeader(table.Row{"Time", "Precipitation", "Temperature", "Condition"})
	location, _ := time.LoadLocation(timezone)
	for _, weather := range forecast.Weathers {
		weather.Timestamp = weather.Timestamp.In(location)

		// Past hours
		if now.Sub(weather.Timestamp).Minutes() > 59 {
			continue
		}

		// Next day
		if now.Day() != weather.Timestamp.Day() {
			continue
		}

		tbl.AppendRow(table.Row{weather.Timestamp, weather.Precipitation, weather.Temperature, weather.Condition})
	}

	tbl.Render()
}

func run(cmd *cobra.Command, args []string) {
	geo, e := getGeoLocation()
	if e != nil {
		fmt.Printf("Cannot fetch geolocation. %s", e)
		return
	}

	forecast := getForecast(geo.Latitude, geo.Longitude)

	printGeo(geo)

	fmt.Println()

	printForecast(forecast, geo.TimezoneName)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
