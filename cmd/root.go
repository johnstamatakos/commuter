package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/marioharper/commuter/directions"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var From string
var To string
var Locations []directions.Location
var Increments int32

var RootCmd = &cobra.Command{
	Use:   "commuter",
	Short: "Tool to get travel time",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		home := directions.Location{
			Name:    "home",
			Address: "4424 Gaines Ranch Loop, Austin, TX",
		}

		work := directions.Location{
			Name:    "work",
			Address: "1835 Kramer Ln, Austin, TX 78758",
		}

		Locations = append(Locations, home)
		Locations = append(Locations, work)

	},
	Run: func(cmd *cobra.Command, args []string) {

		from := Locations[getLocationByName(Locations, From)]
		to := Locations[getLocationByName(Locations, To)]
		currTime := time.Now().Unix()
		minute := 60
		increment := int64(15 * minute)
		var traveTime int64 // time leaving
		var info directions.CommuteInfo

		commute := directions.Commute{
			From: from,
			To:   to,
		}

		fmt.Printf("\nCommute from %s to %s\n", commute.From.Name, commute.To.Name)
		for i := 0; i < 5; i++ {
			var printTime string
			traveTime = currTime + (increment * int64(i)) // time leaving
			info = commute.GetInfo(traveTime)
			hr, min, sec := time.Unix(traveTime, 0).Clock()
			amPm := "AM"

			if hr > 12 {
				hr -= 12
				amPm = "PM"
			} else if hr == 0 {
				hr = 12
			}

			if i == 0 {
				printTime = "Now"
			} else {
				printTime = fmt.Sprintf("%d:%d:%d %s", hr, min, sec, amPm)
			}
			fmt.Printf("\n %s: %d minutes \n", printTime, int(info.TotalDuration))
		}

	},
}

func getLocationByName(locations []directions.Location, name string) int {

	for idx, location := range locations {
		if location.Name == name {
			return idx
		}
	}

	return -1
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&From, "from", "f", "work", "Starting location name")
	RootCmd.PersistentFlags().StringVarP(&To, "to", "t", "home", "Destination location name")
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.commuter-cli.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".commuter") // name of config file (without extension)
	viper.AddConfigPath("$HOME")     // adding home directory as first search path
	viper.AutomaticEnv()             // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}