package main

import (
	"fmt"
	"time"
	"rrd"
	"strings"
	"io/ioutil"
	"regexp"
	
)
const (
		dbfile    = "/tmp/test.rrd"
		step      = 1
		heartbeat = 2 * step
		key		  = "p4p1"
		
)

var myregexp = regexp.MustCompile(`(.+?):\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)`)
var traffic_in_map map[string] string
var traffic_out_map map[string] string

func read_dev () {
	ra, _ := ioutil.ReadFile("/proc/net/dev")
	list := strings.Split(strings.TrimSpace(string(ra)),"\n")

	traffic_in_map = make(map[string]string)
	traffic_out_map = make(map[string]string)
    for index,value := range list {
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
	c := rrd.NewCreator(dbfile, time.Now(), step)
//	for key,_ := range traffic_in_map {
//		c.RRA("AVERAGE", 0.5, 1, 100)
//		c.RRA("AVERAGE", 0.5, 5, 100)
//		in_name := key + "_in"
//		out_name := key + "_out"
//		fmt.Println(in_name)
//		c.DS(in_name,"GAUGE",heartbeat,0,100)
//		c.DS(out_name,"GAUGE",heartbeat,0,100)
//	}
	c.RRA("AVERAGE", 0.5, 5, 600)
	c.RRA("AVERAGE", 0.5, 5, 600)
	
	c.RRA("AVERAGE",0.5, 5, 600)
	
	c.DS("p4p1_in", "COUNTER", heartbeat, 0, "U")
	c.DS("p4p1_out", "COUNTER", heartbeat, 0, "U")
	c.DS("cnt","COUNTER",heartbeat , 0 , "U")
	
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
		str = append(str,"N")
		str = append(str,traffic_in_map["p4p1"])
		str = append(str,traffic_out_map["p4p1"])
		str = append(str,fmt.Sprintf("%d",i))
				
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
	g.Def("v1", dbfile, "p4p1_in", "AVERAGE")
	g.Def("v2", dbfile, "p4p1_out", "AVERAGE")
	
	g.Def("v3", dbfile, "cnt" , "AVERAGE")
	
	g.VDef("max1", "v1,MAXIMUM")
	g.VDef("avg2", "v2,AVERAGE")
	g.VDef("avg3", "v3,AVERAGE") //avg3 是可变无特殊含义
	
	g.Line(1, "v1", "ff0000", "p4p1_in")
	g.Line(1, "v2", "0000ff", "p4p1_out")
	
	g.Line(1, "v3", "00ff00", "cnt")
	
	g.GPrintT("max1", "max1 at %c")
	g.GPrint("avg2", "avg2=%lf")
	g.PrintT("max1", "max1 at %c")
	g.Print("avg2", "avg2=%lf")

	now := time.Now()
	fmt.Println(now.Add(-1800*time.Second))
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

func main(){
	creat_rrd()
	update()
//	grapher()
}