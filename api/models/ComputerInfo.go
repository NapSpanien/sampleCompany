package models

import "fmt"

type ComputerInfo struct {
	MAC          string `json:"MAC"`
	ComputerName string `json:"computerName"`
	IP           string `json:"ip"`
	Employee     string `json:"employee"`
	Description  string `json:"description"`
}

var AssignedComputersPerUser map[string][]string

var Computers []ComputerInfo

func init() {
	AssignedComputersPerUser = make(map[string][]string)
}

//remove removes element without caring about the order
func remove(s []ComputerInfo, i int) []ComputerInfo {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

//remove removes element without caring about the order
func removeString(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// InsertInfo inserts new entry by MAC
func (computerInfo ComputerInfo) InsertInfo() bool {
	//This allows duplicates, with DB doesn't happen because of primary key
	Computers = append(Computers, computerInfo)
	deviceList, exist := AssignedComputersPerUser[computerInfo.Employee]
	if !exist {
		AssignedComputersPerUser[computerInfo.Employee] = []string{computerInfo.MAC}
		return false
	}
	deviceList = append(deviceList, computerInfo.MAC)
	AssignedComputersPerUser[computerInfo.Employee] = deviceList
	return len(deviceList) >= 3
}

//GetInfoByMac returns matching computer to a MAC
func (computerInfo ComputerInfo) GetInfoByMac() ComputerInfo {
	var matchingData ComputerInfo

	for index, info := range Computers {
		if info.MAC == computerInfo.MAC {
			return Computers[index]
		}
	}
	return matchingData
}

//GetComputerInfoByEmployee returns matching computers for an employee
func (computerInfo ComputerInfo) GetComputerInfoByEmployee() []ComputerInfo {
	var matchingData []ComputerInfo

	for index, info := range Computers {
		if info.Employee == computerInfo.Employee {
			matchingData = append(matchingData, Computers[index])
		}
	}
	return matchingData
}

//UpdateInfoByMac Overwrittes saved data with entry data
func (computerInfo ComputerInfo) UpdateInfoByMac() {
	for index, info := range Computers {
		if info.MAC == computerInfo.MAC {
			Computers[index] = computerInfo
			return
		}
	}
}

//DeleteInfoByMac deletes entry from list
func (computerInfo *ComputerInfo) DeleteInfoByMac() {
	for index, info := range Computers {
		if info.MAC == computerInfo.MAC {
			Computers = remove(Computers, index)
			for index, value := range AssignedComputersPerUser[computerInfo.Employee] {
				if value == computerInfo.MAC {
					AssignedComputersPerUser[computerInfo.Employee] = removeString(AssignedComputersPerUser[computerInfo.Employee], index)
					return
				}
			}
		}
	}
}

//GetAllComputers gets all computers stored
func GetAllComputers() []ComputerInfo {
	return Computers
}

//GetAllSavedComputersFromEmployee GetAllSavedComputersFromEmployee
func GetAllSavedComputersFromEmployee(employee string) []ComputerInfo {
	var matchingData []ComputerInfo

	for index, info := range Computers {
		if employee == info.Employee {
			matchingData = append(matchingData, Computers[index])
		}
	}
	return matchingData
}

//AssignComputer AssignComputer
func AssignComputer(mac, employee string) bool {
	deviceList, exist := AssignedComputersPerUser[employee]
	if !exist {
		AssignedComputersPerUser[employee] = []string{mac}
		fmt.Printf("\nassigned computers: %v\n", AssignedComputersPerUser[employee])
		return false
	}
	AssignedComputersPerUser[employee] = append(deviceList, mac)
	fmt.Printf("\nassigned computers: %v\n", AssignedComputersPerUser[employee])
	return len(AssignedComputersPerUser[employee]) >= 3
}

//DisassignComputer DisassignComputer
func DisassignComputer(mac, employee string) {
	for index, value := range AssignedComputersPerUser[employee] {
		if value == mac {
			AssignedComputersPerUser[employee] = removeString(AssignedComputersPerUser[employee], index)
			return
		}
	}
}

//Example with database...
// //GETALLCOMPUTERS ...
// func (computerinfo *ComputerInfo) GetAllComputers() ([]ComputerInfo,error) {
// 	var info []ComputerInfo
// 	var employee sql.NullString
// 	var description sql.NullString
// 	stmt, err := db.Prepare("SELECT MAC, computerName, ip, employee, description FROM sampleCompany.computeres")
// 	if err != nil {
// 		return err
// 	}
// 	defer stmt.Close()
// 	rows, err := stmt.Query()
// 	if err != nil {
// 		return err
// 	}
// 	for rows.Next() {
// 		var mac, computerName, ip string
// 		rows.Scan(&mac, &computerName, &ip)
// 		computerInfo := ComputerInfo{
// 			MAC: mac,
// 			ComputerName: computerName,
// 			IP: ip,
// 		}
// 		if employee.Valid(){
// 			computerInfo.Employee = employee.String
// 		}
// 		if description.Valid(){
// 			computerInfo.Description = description.String
// 		}
// 		info = append(info, computerInfo)
// 	}
// 	defer rows.Close()
// 	return info, err
// }
