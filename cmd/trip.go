/*
Copyright © 2022 Robert Sigler <sigler@improvisedscience.org>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"

	"github.com/davecgh/go-spew/spew"
	"github.com/rmrfslashbin/gomarta/pkg/buses"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// tripCmd represents the trip command
var tripCmd = &cobra.Command{
	Use:   "trip",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Catch errors
		var err error
		defer func() {
			if err != nil {
				log.WithFields(logrus.Fields{
					"error": err,
				}).Fatal("main crashed")
			}
		}()
		if err := getTrips(); err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("error")
		}
	},
}

func init() {
	rootCmd.AddCommand(tripCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tripCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tripCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getTrips() error {
	url := viper.GetString("gtfs.bus.trips")
	if url == "" {
		return errors.New("trips URL is empty- check the config file")
	}

	trips, err := buses.GetData(url, log)
	if err != nil {
		return err
	}

	for _, trip := range trips {
		spew.Dump(trip)
		//x := trip.GetTripUpdate().GetStopTimeUpdate()

	}

	return nil
}

/*
(
	id:"7118984"
	trip_update:{
		trip:{
			trip_id:"7118984"
			route_id:"17052"
			direction_id:9
			start_date:"20220821"}
		vehicle:{
			id:"4771"
			label:"1589"}
		stop_time_update:{
			stop_sequence:35
			stop_id:"58038"
			arrival:{
				delay:-199
				time:1661119296
				4:1661119495
			}
			departure:{
				delay:-199
				time:1661119296
				4:1661119495
			}
			6:"\x12\x08\n\x06\n\x00\x12\x02en"
		}
		stop_time_update:{
			stop_sequence:36
			stop_id:"58040" arrival:{delay:-193 time:1661119313 4:1661119506} departure:{delay:-193 time:1661119313 4:1661119506} 6:"\x12\x08\n\x06\n\x00\x12\x02en"} stop_time_update:{stop_sequence:37 stop_id:"58042" arrival:{delay:-201 time:1661119337 4:1661119538} departure:{delay:-201 time:1661119337 4:1661119538} 6:"\x12\x08\n\x06\n\x00\x12\x02en"} stop_time_update:{stop_sequence:38 stop_id:"58044" arrival:{delay:-215 time:1661119355 4:1661119570} departure:{delay:-215 time:1661119355 4:1661119570} 6:"\x12\x08\n\x06\n\x00\x12\x02en"} stop_time_update:{stop_sequence:39 stop_id:"68016" arrival:{delay:-230 time:1661119377 4:1661119607} departure:{delay:-230 time:1661119377 4:1661119607} 6:"\x12\x08\n\x06\n\x00\x12\x02en"} stop_time_update:{stop_sequence:40 stop_id:"68018" arrival:{delay:-236 time:1661119416 4:1661119652} departure:{delay:-236 time:1661119416 4:1661119652} 6:"\x12\x08\n\x06\n\x00\x12\x02en"} stop_time_update:{stop_sequence:41 stop_id:"68020" arrival:{delay:-255 time:1661119462 4:1661119717} departure:{delay:-255 time:1661119462 4:1661119717} 6:"\x12\x08\n\x06\n\x00\x12\x02en"} stop_time_update:{stop_sequence:42 stop_id:"68022" arrival:{delay:-263 time:1661119481 4:1661119744} departure:{delay:-263 time:1661119481 4:1661119744} 6:"\x12\x08\n\x06\n\x00\x12\x02en"} stop_time_update:{stop_sequence:43 stop_id:"68024" arrival:{delay:-279 time:1661119513 4:1661119792} departure:{delay:-279 time:1661119513 4:1661119792} 6:"\x12\x08\n\x06\n\x00\x12\x02en"} stop_time_update:{stop_sequence:44 stop_id:"68026" arrival:{delay:-291 time:1661119558 4:1661119849} departure:{delay:-291 time:1661119558 4:1661119849} 6:"\x12\x08\n\x06\n\x00\x12\x02en"} stop_time_update:{stop_sequence:45 stop_id:"68028" arrival:{delay:-299 time:1661119579 4:1661119878} departure:{delay:-299 time:1661119579 4:1661119878} 6:"\x12\x08\n\x06\n\x00\x12\x02en"} stop_time_update:{stop_sequence:46 stop_id:"68032" arrival:{delay:-336 time:1661119637 4:1661119973} departure:{delay:-336 time:1661119637 4:1661119973} 6:"\x12\x08\n\x06\n\x00\x12\x02en"} stop_time_update:{stop_sequence:47 stop_id:"68034" arrival:{delay:-380 time:1661119684 4:1661120064} departure:{delay:-380 time:1661119684 4:1661120064} 6:"\x12\x08\n\x06\n\x00\x12\x02en"} stop_time_update:{stop_sequence:48 stop_id:"68038" arrival:{delay:-382 time:1661119701 4:1661120083} departure:{delay:-382 time:1661119701 4:1661120083} 6:"\x12\x08\n\x06\n\x00\x12\x02en"} stop_time_update:{stop_sequence:49 stop_id:"68042" arrival:{delay:-407 time:1661119745 4:1661120152} departure:{delay:-407 time:1661119745 4:1661120152} 6:"\x12\x08\n\x06\n\x00\x12\x02en"} stop_time_update:{stop_sequence:50 stop_id:"68044" arrival:{delay:-404 time:1661119798 4:1661120202} departure:{delay:-404 time:1661119798 4:1661120202} 6:"\x12\x08\n\x06\n\x00\x12\x02en"} stop_time_update:{stop_sequence:51 stop_id:"99973485" arrival:{delay:-422 time:1661119834 4:1661120256} departure:{delay:-422 time:1661119834 4:1661120256} 6:"\x12\x08\n\x06\n\x00\x12\x02en"} stop_time_update:{stop_sequence:52 stop_id:"68900" arrival:{delay:-506 time:1661119834 4:1661120340} departure:{delay:-506 time:1661119834 4:1661120340} 6:"\x12\x08\n\x06\n\x00\x12\x02en\x18\x01"} timestamp:1661119737 6:"\n\x077118984\x12\x0820220821\"\x010*\x05170522\x06\n\x04\x12\x02en:\x08\n\x06\n\x00\x12\x02enB\x073303466" 7:"" 8:"0"})
*/

/*
(
	id:"7097159"
	trip_update:{
		trip:{
			trip_id:"7097159"
			route_id:"17013"
			direction_id:5
			start_time:"23:30:00"
			start_date:"20220821"
			schedule_relationship:CANCELED
		}
		stop_time_update:{
			stop_sequence:1
			stop_id:"68900"
			schedule_relationship:NO_DATA
		}
		stop_time_update:{
			stop_sequence:2
			stop_id:"68056"
			schedule_relationship:NO_DATA
		}
		stop_time_update:{
			stop_sequence:3
			stop_id:"68057"
			schedule_relationship:NO_DATA
		}
		stop_time_update:{
			stop_sequence:4
			stop_id:"69356"
			schedule_relationship:NO_DATA
		}
		stop_time_update:{
			stop_sequence:5
			stop_id:"99972303"
			schedule_relationship:NO_DATA
		}
		stop_time_update:{
			stop_sequence:6
			stop_id:"69325"
			schedule_relationship:NO_DATA
		}
		stop_time_update:{
			stop_sequence:7
			stop_id:"69326"
			schedule_relationship:NO_DATA
		}
		stop_time_update:{
			stop_sequence:8
			stop_id:"69328"
			schedule_relationship:NO_DATA
		}
		stop_time_update:{
			stop_sequence:9
			stop_id:"69330"
			schedule_relationship:NO_DATA
		}
		stop_time_update:{
			stop_sequence:10
			stop_id:"69332"
			schedule_relationship:NO_DATA
		}
		stop_time_update:{
			stop_sequence:11
			stop_id:"69334"
			schedule_relationship:NO_DATA
		}
		stop_time_update:{
			stop_sequence:12
			stop_id:"99973628"
			schedule_relationship:NO_DATA
		}
		stop_time_update:{
			stop_sequence:13
			stop_id:"99973941"
			schedule_relationship:NO_DATA
		}
		stop_time_update:{
			stop_sequence:14
			stop_id:"69341"
			schedule_relationship:NO_DATA
		}
		stop_time_update:{
			stop_sequence:15
			stop_id:"99973911"
			schedule_relationship:NO_DATA
		}
		stop_time_update:{
			stop_sequence:16
			stop_id:"69042"
			schedule_relationship:NO_DATA
		}
		stop_time_update:{stop_sequence:17 stop_id:"99973943" schedule_relationship:NO_DATA} stop_time_update:{stop_sequence:18 stop_id:"59042" schedule_relationship:NO_DATA} stop_time_update:{stop_sequence:19 stop_id:"99970171" schedule_relationship:NO_DATA} stop_time_update:{stop_sequence:20 stop_id:"99970172" schedule_relationship:NO_DATA} stop_time_update:{stop_sequence:21 stop_id:"59268" schedule_relationship:NO_DATA} stop_time_update:{stop_sequence:22 stop_id:"59262" schedule_relationship:NO_DATA} stop_time_update:{stop_sequence:23 stop_id:"99973924" schedule_relationship:NO_DATA} stop_time_update:{stop_sequence:24 stop_id:"99973925" schedule_relationship:NO_DATA} stop_time_update:{stop_sequence:25 stop_id:"99974102" schedule_relationship:NO_DATA} stop_time_update:{stop_sequence:26 stop_id:"59374" schedule_relationship:NO_DATA} stop_time_update:{stop_sequence:27 stop_id:"59334" schedule_relationship:NO_DATA} stop_time_update:{stop_sequence:28 stop_id:"59338" schedule_relationship:NO_DATA} stop_time_update:{stop_sequence:29 stop_id:"59371" schedule_relationship:NO_DATA} stop_time_update:{stop_sequence:30 stop_id:"59340" schedule_relationship:NO_DATA} stop_time_update:{stop_sequence:31 stop_id:"59342" schedule_relationship:NO_DATA} stop_time_update:{stop_sequence:32 stop_id:"59344" schedule_relationship:NO_DATA} stop_time_update:{stop_sequence:33 stop_id:"59346" schedule_relationship:NO_DATA} stop_time_update:{stop_sequence:34 stop_id:"59348" schedule_relationship:NO_DATA} stop_time_update:{stop_sequence:35 stop_id:"52220" schedule_relationship:NO_DATA} stop_time_update:{stop_sequence:36 stop_id:"52500" schedule_relationship:NO_DATA}
		timestamp:1661119815
		6:"\n\x077097159\x12\x0820220821\x1a\x0823:30:00\"\x010*\x05170132\x06\n\x04\x12\x02en:\x08\n\x06\n\x00\x12\x02enB\"TransitMaster.DataCube.BlockTripID"
		7:""
		8:"0"
	}
)
*/
