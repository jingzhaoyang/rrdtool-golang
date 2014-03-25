package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"rrd"
	"strings"
	"time"
)

const (
	dbfile    = "/tmp/test.rrd"
	step      = 1
	heartbeat = 2 * step
	// key		  = "p4p1"

)

var myregexp = regexp.MustCompile(`(.+?):\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)`)
var traffic_in_map map[string]string
var traffic_out_map map[string]string

func read_dev() {
	ra, _ := ioutil.ReadFile("/proc/net/dev")
	list := strings.Split(strings.TrimSpace(string(ra)), "\n")

	traffic_in_map = make(map[string]string)
	traffic_out_map = make(map[string]string)
	for index, value := range list {
		if index > 1 {
			list_a := myregexp.FindStringSubmatch(strings.TrimSpace(value))
			//			fmt.Println(list_a)
			traffic_in_map[list_a[1]] = list_a[2]
			traffic_out_map[list_a[1]] = list_a[10]
			//			fmt.Println(traffic_in_map[list_a[1]])
			//			fmt.Println(traffic_out_map[list_a[1]])
		}
	}
	//    return traffic_in_map,traffic_out_map
}
func creat_rrd() {
	// Create
	read_dev()
	c := rrd.NewCreator(dbfile, time.Now(), step)
	for key, _ := range traffic_in_map {

		c.RRA("AVERAGE", 0.5, 1, 600)
		fmt.Println("122")
		c.RRA("AVERAGE", 0.5, 5, 600)
		in_name := key + "_in"
		out_name := key + "_out"
		fmt.Println(in_name)
		c.DS(in_name, "COUNTER", heartbeat, 0, "U")
		c.DS(out_name, "COUNTER", heartbeat, 0, "U")
	}
	// c.RRA("AVERAGE", 0.5, 5, 600)
	// c.RRA("AVERAGE", 0.5, 5, 600)
	// c.DS("p4p1_in", "COUNTER", heartbeat, 0, "U")
	// c.DS("p4p1_out", "COUNTER", heartbeat, 0, "U")

	err := c.Create(true)
	if err != nil {
		fmt.Println(err)
	}
}
func update() {

	// Update
	u := rrd.NewUpdater(dbfile)
	for i := 0; i < 1800; i++ {
		read_dev()
		time.Sleep(step * time.Second)
		str := []string{}
		str = append(str, "N")
		for key, value := range traffic_in_map {
			str = append(str, value)
			str = append(str, traffic_out_map[key])
		}

		//		fmt.Println(time.Now(), i, 1.5*float64(i))
		//		for key,value := range traffic_in_map {
		//			fmt.Println(key,value)
		//		}
		new_str := make([]interface{}, len(str)) // 将str的slicen string类型转化为interface类型
		for i, v := range str {
			new_str[i] = interface{}(v)
		}

		fmt.Println(new_str...)
		//		err := u.Update("N",  traffic_in_map["p4p1"],traffic_out_map["p4p1"],i)
		err := u.Update(new_str...)
		if err != nil {
			fmt.Println(err)
		}
		grapher()
	}
}

func grapher() {

	g := rrd.NewGrapher()
	g.SetTitle("Test")
	g.SetVLabel("some variable")
	g.SetSize(800, 300)
	g.SetWatermark("some watermark")
	for key, _ := range traffic_in_map {
		in_name := key + "_in"
		out_name := key + "_out"
		max_in_name := key + "_max"
		avg_out_name := key + "_avg"
		g.Def(in_name, dbfile, in_name, "AVERAGE")
		g.Def(out_name, dbfile, out_name, "AVERAGE")
		g.VDef(max_in_name, fmt.Sprintf("%s,MAXIMUM", in_name))
		g.VDef(avg_out_name, fmt.Sprintf("%s,AVERAGE", out_name))
		g.Line(1, in_name, "ff0000", in_name)
		g.Line(1, out_name, "0000ff", out_name)
	}

	//	g.GPrintT("max1", "max1 at %c")
	//	g.GPrint("avg2", "avg2=%lf")
	//	g.PrintT("max1", "max1 at %c")
	//	g.Print("avg2", "avg2=%lf")

	now := time.Now()
	fmt.Println(now.Add(-1800 * time.Second))
	i, err := g.SaveGraph("/tmp/test_rrd1.png", now.Add(-1800*time.Second), now)
	fmt.Printf("%+v\n", i)
	if err != nil {
		fmt.Println(err)
		//	i, buf, err := g.Graph(now.Add(-1800*time.Second), now)
		//	i, _, err := g.Graph(now.Add(-1800*time.Second), now)
		//	fmt.Printf("%+v\n", i)
		//	if err != nil {
		//		fmt.Println(err)
		//	}
		//	err = ioutil.WriteFile("/tmp/test_rrd2.png", buf, 0666)
		//	if err != nil {
		//		fmt.Println(err)
		//	}
	}
}

func main() {
	creat_rrd()
	update()
	//	grapher()
}
