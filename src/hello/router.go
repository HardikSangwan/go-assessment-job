/*
Basic Application Router for parsing routing rules from config file and returning server routes for a given input based on the rules
Three Major Assumptions made for now just to showcase basic idea behind the program:
1. Configs would most likely be sent as json/yaml/ini/etc. file but implemented here as a txt file
2. Config File only contains the expected rules
3. Input String for findRouter function is as expected

Additional handling for errors and different types of formats would be needed for a proper application router
*/

package main

import (
	"os"
	"fmt"
	"bufio"
	"log"
	"strings"
)

//Declare config structure
type Config struct {
	Customer_id string
	Country string
	State string
	City string
	Server string
}
var configs []Config

//Parse and Load Routing Config File into Internal Data Structure
func loadconfig() {
	//Access Config File
	file, err := os.Open("configfile.txt") 
	if err != nil {  
		log.Fatal(err) 
	}
	defer file.Close()

	//Count number of lines/routes for config struct array size
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan(){
		count++
	}

	//Seek to file beginning again and read for configs
	file.Seek(0,0)
	configs = make([]Config, count)
	scanner = bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		split1 := strings.SplitN(scanner.Text(), "=", 2)
		split2 := strings.SplitN(split1[0], ".", 4)
		configs[i].Server = strings.TrimSpace(split1[1])
		configs[i].Customer_id = strings.TrimSpace(split2[0])
		configs[i].Country = strings.TrimSpace(split2[1])
		configs[i].State = strings.TrimSpace(split2[2])
		configs[i].City = strings.TrimSpace(split2[3])
		i++
	}
}
func findRoute(input string) string {
	//Separate out customer_id, country, state and city variables from given string
	split_input := strings.SplitN(input, ".", 4)
	customer_id := split_input[0]
	country := split_input[1]
	state := split_input[2]
	city := split_input[3]

	//These variables are used to check for the 'specificity of a match' IN ORDER, and store the respective rule index. 
	//for example, match0 would mean that none of the variables matched but there did exist a rule with *.*.*.* 
	//that can be used for allocating the appropriate server. 
	match_0 := false
	match_1 := false
	match_2 := false
	match_3 := false
	match_4 := false
	match_0id := -1
	match_1id := -1
	match_2id := -1
	match_3id := -1
	match_4id := -1

	//The for loop goes through all the stored rules. Breaks only if an exact match is found.
	for i := range configs {
		if (configs[i].City == city) && (configs[i].State == state) && (configs[i].Country == country) && (configs[i].Customer_id == customer_id) {
			match_4 = true
			match_4id = i
			break
		} else if (configs[i].City == "*") && (configs[i].State == state) && (configs[i].Country == country) && (configs[i].Customer_id == customer_id) {
			match_3 = true
			match_3id = i
		} else if (configs[i].City == "*") && (configs[i].State == "*") && (configs[i].Country == country) && (configs[i].Customer_id == customer_id) {
			match_2 = true
			match_2id = i
		} else if (configs[i].City == "*") && (configs[i].State == "*") && (configs[i].Country == "*") && (configs[i].Customer_id == customer_id) {
			match_1 = true
			match_1id = i
		} else if (configs[i].City == "*") && (configs[i].State == "*") && (configs[i].Country == "*") && (configs[i].Customer_id == "*") {
			match_0 = true
			match_0id = i
		}
	}
	if match_4 {
		return configs[match_4id].Server
	} else if match_3 {
		return configs[match_3id].Server
	} else if match_2 {
		return configs[match_2id].Server
	} else if match_1 {
		return configs[match_1id].Server
	} else if match_0 {
		return configs[match_0id].Server
	} else {
		return "No server match found"
	}
}

func main() {
	loadconfig()
	fmt.Println(findRoute("customer1.us.ca.sfo"))
	fmt.Println(findRoute("customer1.us.ca.sjc"))
	fmt.Println(findRoute("customer2.us.tx.dfw"))
	fmt.Println(findRoute("customer2.cn.tw.tai"))
	fmt.Println(findRoute("customer10.us.ny.nyc"))
}