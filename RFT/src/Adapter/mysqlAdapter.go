package Adapter

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type SQLConfig struct {
	username string
	password string
	host     string
	port     int
	db       string
}

type Authentification struct {
	Valid     bool
	Username  string
	Firstname string
	Lastname  string
	Email			string
}

type TrainServices struct {
	S1_Toilet         bool
	S1_DisabledToilet bool
	S1_DiaperChange   bool
	S1_AirConditioner bool
	S1_Wifi           bool
	S1_PowerConnector bool
	S1_Restaurant     bool
	S1_BikeShed       bool
	S1_Bed            bool
	S2_Toilet         bool
	S2_DisabledToilet bool
	S2_DiaperChange   bool
	S2_AirConditioner bool
	S2_Wifi           bool
	S2_PowerConnector bool
	S2_Restaurant     bool
	S2_BikeShed       bool
	S2_Bed            bool
}

type TrainInformation struct {
	Station   []string
	Timetable []string
	Train     []string
	RouteID   []string
	Services  []TrainServices
}

type SQLData struct {
	From      string
	To        string
	Date      string
	Departure string
	Arrival   string
	Changes   string
	Duration  string
	Distance  string
	Price     string
	Info      TrainInformation
}

type SQLDataSlice []SQLData

type Data struct {
	From string
	To   string
	Data SQLDataSlice
}

type MapData struct {
	From      string
	To        string
	Departure string
	Arrival   string
	Stations  []string
}

type WagonServices struct {
	Toilet         bool
	DisabledToilet bool
	DiaperChange   bool
	AirConditioner bool
	Wifi           bool
	PowerConnector bool
	Restaurant     bool
	BikeShed       bool
	Bed            bool
}

type Wagon struct {
	WagonID       string
	NumberOfSeats string
	Class         string
	Services      WagonServices
	Seats         map[string]bool
}

type WagonData struct {
	Wagons        []Wagon
	Train         TrainType
	SelectedTrain string
}

type TrainType struct {
	From1       string
	To1         string
	Departure1  string
	Arrival1    string
	Train1ID    string
	Type1       bool
	SeatNumber1 string
	From2       string
	To2         string
	Departure2  string
	Arrival2    string
	Train2ID    string
	Type2       bool
	SeatNumber2 string
	Price       string
	Km          string
}

type Ticket struct {
	Firstname		string
	Lastname		string
	Date				string
	From1				string
	To1					string
	Departure1	string
	Arrival1		string
	Train1ID		string
	Seat1				string
	From2				string
	To2					string
	Departure2	string
	Arrival2		string
	Train2ID		string
	Seat2				string
	TicketID		string
	TicketPassw	string
	Price				string
	Km					string
}

type PurchaseInfo struct {
	TicketID		string
	TicketPassw	string
}

func SQLFactory(username, password, host, db string, port int) SQLConfig {
	return SQLConfig{username, password, host, port, db}
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func generateSalt() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (this SQLConfig) MysqlAuthentificate(username, password string) Authentification {
	var usern, passw, salt1, salt2, fname, lname, email string
	var auth Authentification

	inf := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", this.username, this.password, this.host, this.port, this.db)
	db, err := sql.Open("mysql", inf)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT NICKNAME, PASSWORD, SALT1, SALT2, FIRSTNAME, LASTNAME, EMAIL FROM USERS WHERE NICKNAME='" + username + "'")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&usern, &passw, &salt1, &salt2, &fname, &lname, &email)
		if err != nil {
			panic(err)
		}
	}

	saltedPassword := fmt.Sprintf(salt1 + password + salt2)
	hashedPassword := GetMD5Hash(saltedPassword)
	if hashedPassword == passw {
		auth.Valid = true
		auth.Username = username
		auth.Firstname = fname
		auth.Lastname = lname
		auth.Email = email
	} else {
		auth.Valid = false
		auth.Username = ""
	}

	return auth
}

func (this SQLConfig) MysqlRegistration(firstname, lastname, username, password, email string) int {
	var MySQLResultUser string

	inf := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", this.username, this.password, this.host, this.port, this.db)
	db, err := sql.Open("mysql", inf)
	if err != nil {
		panic(err)
		return -2
	}
	defer db.Close()

	rows, err := db.Query("SELECT NICKNAME FROM USERS WHERE NICKNAME='" + username + "'")
	if err != nil {
		panic(err)
		return -2
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&MySQLResultUser)
		if err != nil {
			panic(err)
			return -2
		}
	}

	if MySQLResultUser != "" {
		return -1
	}

	salt1 := generateSalt()
	salt2 := generateSalt()
	saltedPassword := fmt.Sprintf(salt1 + password + salt2)
	hashedPassword := GetMD5Hash(saltedPassword)

	_, err = db.Exec("INSERT INTO USERS (FIRSTNAME, LASTNAME, NICKNAME, PASSWORD, SALT1, SALT2, TYPE, EMAIL) VALUES ('" + firstname + "', '" + lastname + "', '" + username + "', '" +
		hashedPassword + "', '" + salt1 + "', '" + salt2 + "', '" + "user" + "', '" + email + "')")

	if err != nil {
		panic(err)
		return -2
	}

	return 0
}

func (this SQLConfig) MysqlSearchTimetable(from, to, date, discount string, extraTicket, change bool) Data {
	var finalResult Data
	var result SQLDataSlice
	var data SQLData
	var ch string
	var distanceAndTimeQuery string
	var d1, d2, d3, d4, d5, d6, d7, d8, d9, d10, d11, d12, d13, d14, d15, d16, d17, d18, d19, d20, d21, d22 string //adatbazis adatok
	pricePerKM := 18

	if !change {
		ch = "1"
	} else {
		ch = "0"
	}
	distanceAndTimeQuery = "SELECT * FROM ( SELECT '0' AS 'COUNT_OF_CHANGES' ,WAGONS.TRAINS_ID AS 'FIRST_TRAIN_ID' ,SPEED AS 'FIRST_TRAIN_SPEED'	,TYPE AS 'FIRST_TRAIN_TYPE'	,DISTANCE_BETWEEN_START_AND_STOP AS 'FIRST_TRAIN_DISTANCE_BETWEEN_START_AND_STOP'	,'0' AS 'CHANGE_STATION'	,CLASS1_SERVICES AS 'FIRST_TRAIN_CLASS1_SERVICES'	,CLASS2_SERVICES AS 'FIRST_TRAIN_CLASS2_SERVICES'	,TMP3_EXPLETIVE_TICKET AS 'FIRST_TRAIN_EXPLETIVE_TICKET'	,START + INTERVAL DISTANCE_BEFORE_START / SPEED * 60.0 MINUTE AS 'START_TIME1'	,START + INTERVAL(DISTANCE_BEFORE_START + DISTANCE_BETWEEN_START_AND_STOP) / SPEED * 60.0 MINUTE AS 'ARRIVE_TIME1' ,'0' AS 'SECOND_TRAIN_ID' ,'0' AS 'SECOND_TRAIN_SPEED'	,'0' AS 'SECOND_TRAIN_TYPE'	,'0' AS 'SECOND_TRAIN_DISTANCE_BETWEEN_START_AND_STOP' ,'0' AS 'SECOND_TRAIN_CLASS1_SERVICES'	,'0' AS 'SECOND_TRAIN_CLASS2_SERVICES' ,'0' AS 'SECOND_TRAIN_EXPLETIVE_TICKET' ,'0' AS 'START_TIME2' ,'0' AS 'ARRIVE_TIME2', ROUTES_ID AS 'FIRST_TRAIN_ROUTE_ID' ,'0' AS 'SECOND_TRAIN_ROUTE_ID' FROM WAGONS INNER JOIN ( SELECT * ,SUM(CASE WHEN STATION_INDEX_IN_ROUTE < START_INDEX_IN_ROUTE THEN DISTANCE ELSE 0 END) AS 'DISTANCE_BEFORE_START' ,SUM(CASE WHEN ( STATION_INDEX_IN_ROUTE >= START_INDEX_IN_ROUTE AND STATION_INDEX_IN_ROUTE < STOP_INDEX_IN_ROUTE ) THEN DISTANCE ELSE 0 END) AS 'DISTANCE_BETWEEN_START_AND_STOP' FROM ( SELECT * FROM ROUTE_STATIONS_CONNECTION INNER JOIN ( SELECT TRAINS.ID ,TRAINS.SPEED ,TRAINS.TYPE ,TRAIN_ROUTE_CONNECTION.ROUTES_ID AS 'TMP4_ROUTES_ID' ,TRAIN_ROUTE_CONNECTION.START ,TMP3.STATION_INDEX_IN_ROUTE AS 'START_INDEX_IN_ROUTE' ,TMP3.TMP2_STATION_INDEX_IN_ROUTE AS 'STOP_INDEX_IN_ROUTE' ,TMP3.STATIONS_ID AS 'START_STATIONS_ID' ,TMP3.TMP2_STATIONS_ID AS 'STOP_STATIONS_ID' ,TMP3_EXPLETIVE_TICKET FROM TRAINS INNER JOIN TRAIN_ROUTE_CONNECTION ON TRAINS.ID = TRAIN_ROUTE_CONNECTION.TRAINS_ID INNER JOIN ( SELECT ROUTE_STATIONS_CONNECTION.ROUTES_ID ,SUM(ROUTE_STATIONS_CONNECTION.EXPLETIVE_TICKET) AS 'TMP3_EXPLETIVE_TICKET' ,ROUTE_STATIONS_CONNECTION.TYPE_OF_TRAINS_STOPS AS 'FROM_TYPE_OF_TRAINS_STOPS' ,TMP.TO_TYPE_OF_TRAINS_STOPS ,TMP.STATIONS_ID ,TMP.TMP2_STATIONS_ID ,TMP.STATION_INDEX_IN_ROUTE,TMP.TMP2_STATION_INDEX_IN_ROUTE FROM ROUTE_STATIONS_CONNECTION INNER JOIN ( SELECT * FROM ROUTE_STATIONS_CONNECTION INNER JOIN STATIONS ON ROUTE_STATIONS_CONNECTION.STATIONS_ID = STATIONS.ID INNER JOIN ( SELECT STATIONS_ID AS 'TMP2_STATIONS_ID' ,STATION_INDEX_IN_ROUTE AS 'TMP2_STATION_INDEX_IN_ROUTE' ,ROUTES_ID AS 'TMP2_ROUTES_ID' ,TYPE_OF_TRAINS_STOPS AS 'TO_TYPE_OF_TRAINS_STOPS' FROM ROUTE_STATIONS_CONNECTION INNER JOIN STATIONS ON ROUTE_STATIONS_CONNECTION.STATIONS_ID = STATIONS.ID WHERE STATIONS.NAME = '" + to + "' ) TMP2 ON ROUTE_STATIONS_CONNECTION.ROUTES_ID = TMP2.TMP2_ROUTES_ID WHERE STATIONS.NAME = '" + from + "'	AND TMP2.TMP2_STATION_INDEX_IN_ROUTE > ROUTE_STATIONS_CONNECTION.STATION_INDEX_IN_ROUTE ) TMP ON ROUTE_STATIONS_CONNECTION.ROUTES_ID = TMP.ROUTES_ID WHERE ROUTE_STATIONS_CONNECTION.STATION_INDEX_IN_ROUTE BETWEEN TMP.STATION_INDEX_IN_ROUTE AND TMP.TMP2_STATION_INDEX_IN_ROUTE GROUP BY ROUTE_STATIONS_CONNECTION.ROUTES_ID ) TMP3 ON TRAIN_ROUTE_CONNECTION.ROUTES_ID = TMP3.ROUTES_ID WHERE TMP3.FROM_TYPE_OF_TRAINS_STOPS LIKE CONCAT ( '%' ,TRAINS.TYPE	,'%')	AND TMP3.TO_TYPE_OF_TRAINS_STOPS LIKE CONCAT ('%'	,TRAINS.TYPE,'%')) TMP4 ON ROUTE_STATIONS_CONNECTION.ROUTES_ID = TMP4.TMP4_ROUTES_ID) TMP8 INNER JOIN ( SELECT ROUTES_ID AS 'TMP9_ROUTES_ID' ,STATION_INDEX_IN_ROUTE AS 'TMP9_STATION_INDEX_IN_ROUTE'	,STATIONS_ID AS 'TMP9_STATIONS_ID' FROM ROUTE_STATIONS_CONNECTION) TMP9 ON TMP8.ROUTES_ID = TMP9.TMP9_ROUTES_ID	AND TMP8.STATION_INDEX_IN_ROUTE + 1 = TMP9.TMP9_STATION_INDEX_IN_ROUTE INNER JOIN (	SELECT NEIGHBOUR_ID	,DISTANCE	,STATIONS_ID AS 'TMP10_STATIONS_ID' FROM NEIGHBOURS) TMP10 ON TMP10.NEIGHBOUR_ID = TMP9.TMP9_STATIONS_ID AND TMP10.TMP10_STATIONS_ID = TMP8.STATIONS_ID GROUP BY ROUTES_ID) TMP7 ON TMP7.ID = WAGONS.TRAINS_ID INNER JOIN ( SELECT TRAINS_ID ,CASE WHEN ISNULL(GROUP_CONCAT(CASE WHEN CLASS = 1	THEN SERVICES	END separator ';'))	THEN 0 ELSE GROUP_CONCAT(CASE WHEN CLASS = 1 THEN SERVICES END separator ';')	END AS 'CLASS1_SERVICES' FROM WAGONS GROUP BY WAGONS.TRAINS_ID) TMP5 ON WAGONS.TRAINS_ID = TMP5.TRAINS_ID INNER JOIN ( SELECT TRAINS_ID ,CASE WHEN ISNULL(GROUP_CONCAT(CASE WHEN CLASS = 2 THEN SERVICES END separator ';')) THEN 0	ELSE GROUP_CONCAT(CASE WHEN CLASS = 2	THEN SERVICES	END separator ';') END AS 'CLASS2_SERVICES' FROM WAGONS GROUP BY WAGONS.TRAINS_ID) TMP6 ON WAGONS.TRAINS_ID = TMP6.TRAINS_ID WHERE DATE (START + INTERVAL DISTANCE_BEFORE_START / SPEED * 60.0 MINUTE) = DATE ('" + date + "') GROUP BY WAGONS.TRAINS_ID UNION ALL SELECT '1' AS 'COUNT_OF_CHANGES'	,FIRST_TRAIN_ID	,FIRST_TRAIN_SPEED ,FIRST_TRAIN_TYPE ,DISTANCE_BETWEEN_START_AND_STOP1 AS 'FIRST_TRAIN_DISTANCE_BETWEEN_START_AND_STOP'	,MIDDLE_STATION_NAME AS 'CHANGE_STATION' ,CLASS1_SERVICES1 AS 'FIRST_TRAIN_CLASS1_SERVICES'	,CLASS2_SERVICES1 AS 'FIRST_TRAIN_CLASS2_SERVICES' ,SUM(EXPLETIVE_TICKET1) AS 'FIRST_TRAIN_EXPLETIVE_TICKET' ,START_TIME1	,ARRIVE_TIME1	,SECOND_TRAIN_ID ,SECOND_TRAIN_SPEED ,SECOND_TRAIN_TYPE	,DISTANCE_BETWEEN_START_AND_STOP2 AS 'SECOND_TRAIN_DISTANCE_BETWEEN_START_AND_STOP'	,CLASS1_SERVICES2 AS 'SECOND_TRAIN_CLASS1_SERVICES'	,CLASS2_SERVICES2 AS 'SECOND_TRAIN_CLASS2_SERVICES'	,SUM(EXPLETIVE_TICKET2) AS 'SECOND_TRAIN_EXPLETIVE_TICKET' ,START_TIME2	,ARRIVE_TIME2 ,ID2 AS 'FIRST_TRAIN_ROUTE_ID' , ID3 AS 'SECOND_TRAIN_ROUTE_ID' FROM (SELECT * ,FIRST_TRAIN_START + INTERVAL DISTANCE_BEFORE_START1 / FIRST_TRAIN_SPEED * 60.0 MINUTE AS 'START_TIME1' ,FIRST_TRAIN_START + INTERVAL(DISTANCE_BEFORE_START1 + DISTANCE_BETWEEN_START_AND_STOP1) / FIRST_TRAIN_SPEED * 60.0 MINUTE AS 'ARRIVE_TIME1'	,SECOND_TRAIN_START + INTERVAL DISTANCE_BEFORE_START2 / SECOND_TRAIN_SPEED * 60.0 MINUTE AS 'START_TIME2'	,SECOND_TRAIN_START + INTERVAL(DISTANCE_BEFORE_START2 + DISTANCE_BETWEEN_START_AND_STOP2) / SECOND_TRAIN_SPEED * 60.0 MINUTE AS 'ARRIVE_TIME2' FROM (	SELECT * ,IF (START_STATION_INDEX_IN_ROUTE < MIDDLE_STATION_INDEX_IN_ROUTE1	,SUM(CASE WHEN TMP7_STATION_INDEX_IN_ROUTE < START_STATION_INDEX_IN_ROUTE	AND TMP7_ROUTES_ID = ID2 THEN DISTANCE ELSE 0	END) ,SUM(CASE WHEN TMP7_STATION_INDEX_IN_ROUTE >= START_STATION_INDEX_IN_ROUTE	AND TMP7_ROUTES_ID = ID2 THEN DISTANCE ELSE 0	END)) AS 'DISTANCE_BEFORE_START1'	,	IF (START_STATION_INDEX_IN_ROUTE < MIDDLE_STATION_INDEX_IN_ROUTE1	,SUM(CASE WHEN (TMP7_STATION_INDEX_IN_ROUTE >= START_STATION_INDEX_IN_ROUTE	AND TMP7_STATION_INDEX_IN_ROUTE < MIDDLE_STATION_INDEX_IN_ROUTE1 AND TMP7_ROUTES_ID = ID2	)	THEN DISTANCE	ELSE 0	END) ,SUM(CASE 	WHEN (TMP7_STATION_INDEX_IN_ROUTE < START_STATION_INDEX_IN_ROUTE AND TMP7_STATION_INDEX_IN_ROUTE >= MIDDLE_STATION_INDEX_IN_ROUTE1 AND TMP7_ROUTES_ID = ID2	)	THEN DISTANCE	ELSE 0 END)) AS 'DISTANCE_BETWEEN_START_AND_STOP1',SUM(CASE WHEN TMP7_STATION_INDEX_IN_ROUTE < MIDDLE_STATION_INDEX_IN_ROUTE2	AND TMP7_ROUTES_ID = ID3 THEN DISTANCE ELSE 0	END) AS 'DISTANCE_BEFORE_START2' ,SUM(CASE WHEN (	TMP7_STATION_INDEX_IN_ROUTE >= MIDDLE_STATION_INDEX_IN_ROUTE2	AND TMP7_STATION_INDEX_IN_ROUTE < STOP_STATION_INDEX_IN_ROUTE2 AND TMP7_ROUTES_ID = ID3) THEN DISTANCE ELSE 0	END) AS 'DISTANCE_BETWEEN_START_AND_STOP2' FROM (SELECT ROUTE_STATIONS_CONNECTION.EXPLETIVE_TICKET AS 'TMP7_EXPLETIVE_TICKET' ,ROUTE_STATIONS_CONNECTION.ROUTES_ID AS 'TMP7_ROUTES_ID' ,ROUTE_STATIONS_CONNECTION.STATIONS_ID AS 'TMP7_STATIONS_ID' ,ROUTE_STATIONS_CONNECTION.STATION_INDEX_IN_ROUTE AS 'TMP7_STATION_INDEX_IN_ROUTE' ,ROUTE_STATIONS_CONNECTION.TYPE_OF_TRAINS_STOPS AS 'TMP7_TYPE_OF_TRAINS_STOPS' ,TMP6.*	FROM ROUTE_STATIONS_CONNECTION INNER JOIN (	SELECT DISTINCT TRAINS.ID AS 'SECOND_TRAIN_ID' ,TRAINS.TYPE AS 'SECOND_TRAIN_TYPE' ,TRAINS.SPEED AS 'SECOND_TRAIN_SPEED' ,TRAINS.STATUS AS 'SECOND_TRAIN_STATUS' ,TRAIN_ROUTE_CONNECTION.START AS 'SECOND_TRAIN_START' ,TMP5.* FROM TRAINS INNER JOIN TRAIN_ROUTE_CONNECTION ON TRAINS.ID = TRAIN_ROUTE_CONNECTION.TRAINS_ID INNER JOIN ( SELECT TRAINS.ID AS 'FIRST_TRAIN_ID' ,TRAINS.TYPE AS 'FIRST_TRAIN_TYPE'	,TRAINS.SPEED AS 'FIRST_TRAIN_SPEED' ,TRAINS.STATUS AS 'FIRST_TRAIN_STATUS' ,TRAIN_ROUTE_CONNECTION.START AS 'FIRST_TRAIN_START' ,TMP4.* FROM TRAINS INNER JOIN TRAIN_ROUTE_CONNECTION ON TRAINS.ID = TRAIN_ROUTE_CONNECTION.TRAINS_ID INNER JOIN ( SELECT * ,ROUTE_STATIONS_CONNECTION.STATION_INDEX_IN_ROUTE AS 'STOP_STATION_INDEX_IN_ROUTE2' FROM ROUTE_STATIONS_CONNECTION INNER JOIN (SELECT ROUTE_STATIONS_CONNECTION.ROUTES_ID AS 'ID3' ,ROUTE_STATIONS_CONNECTION.STATION_INDEX_IN_ROUTE AS 'MIDDLE_STATION_INDEX_IN_ROUTE2' ,TMP2.MIDDLE_STATION_INDEX_IN_ROUTE1 ,MIDDLE_STATION_NAME ,TMP2.ID2 ,TMP2.START_STATION_INDEX_IN_ROUTE ,TMP2.START_TRAINS_STOP ,ROUTE_STATIONS_CONNECTION.TYPE_OF_TRAINS_STOPS AS 'MIDDLE_TRAINS_STOP' FROM ROUTE_STATIONS_CONNECTION INNER JOIN ( SELECT ROUTE_STATIONS_CONNECTION.STATIONS_ID ,ROUTES_ID AS 'ID2' ,TMP.FIRST_STATION_ID ,TMP.START_STATION_INDEX_IN_ROUTE ,TMP.START_TRAINS_STOP ,ROUTE_STATIONS_CONNECTION.STATION_INDEX_IN_ROUTE AS 'MIDDLE_STATION_INDEX_IN_ROUTE1' ,STATIONS.NAME AS 'MIDDLE_STATION_NAME' FROM ROUTE_STATIONS_CONNECTION INNER JOIN ( SELECT ROUTES_ID AS 'ID1' ,STATION_INDEX_IN_ROUTE AS 'START_STATION_INDEX_IN_ROUTE' ,TYPE_OF_TRAINS_STOPS AS 'START_TRAINS_STOP' ,STATIONS_ID AS 'FIRST_STATION_ID' FROM ROUTE_STATIONS_CONNECTION INNER JOIN STATIONS ON STATIONS.ID = ROUTE_STATIONS_CONNECTION.STATIONS_ID WHERE STATIONS.NAME = '" + from + "' ) TMP ON ROUTE_STATIONS_CONNECTION.ROUTES_ID = TMP.ID1 	INNER JOIN STATIONS ON STATIONS.ID = ROUTE_STATIONS_CONNECTION.STATIONS_ID ) TMP2 ON ROUTE_STATIONS_CONNECTION.STATIONS_ID = TMP2.STATIONS_ID WHERE ROUTE_STATIONS_CONNECTION.ROUTES_ID != TMP2.ID2 AND TMP2.MIDDLE_STATION_NAME != '" + to + "' AND FIRST_STATION_ID != ROUTE_STATIONS_CONNECTION.STATIONS_ID AND MIDDLE_STATION_INDEX_IN_ROUTE1 > START_STATION_INDEX_IN_ROUTE) TMP3 ON TMP3.ID3 = ROUTE_STATIONS_CONNECTION.ROUTES_ID INNER JOIN STATIONS ON STATIONS.ID = ROUTE_STATIONS_CONNECTION.STATIONS_ID	WHERE STATIONS.NAME = '" + to + "' AND ROUTE_STATIONS_CONNECTION.STATION_INDEX_IN_ROUTE > MIDDLE_STATION_INDEX_IN_ROUTE2 GROUP BY ID2 ,ID3 HAVING MIN(MIDDLE_STATION_INDEX_IN_ROUTE1)) TMP4 ON TRAIN_ROUTE_CONNECTION.ROUTES_ID = TMP4.ID2) TMP5 ON TRAIN_ROUTE_CONNECTION.ROUTES_ID = TMP5.ID3	WHERE START_TRAINS_STOP LIKE CONCAT ('%' ,TMP5.FIRST_TRAIN_TYPE	,'%')	AND MIDDLE_TRAINS_STOP LIKE CONCAT ('%',TRAINS.TYPE	,'%') AND TMP5.TYPE_OF_TRAINS_STOPS LIKE CONCAT ('%',TRAINS.TYPE,'%')) TMP6 ON ROUTE_STATIONS_CONNECTION.ROUTES_ID = TMP6.ID3	OR ROUTE_STATIONS_CONNECTION.ROUTES_ID = TMP6.ID2) TMP7 INNER JOIN (SELECT ROUTES_ID AS 'TMP8_ROUTES_ID',STATION_INDEX_IN_ROUTE AS 'TMP8_STATION_INDEX_IN_ROUTE' ,STATIONS_ID AS 'TMP8_STATIONS_ID' FROM ROUTE_STATIONS_CONNECTION) TMP8 ON TMP7.TMP7_ROUTES_ID = TMP8.TMP8_ROUTES_ID AND TMP7.TMP7_STATION_INDEX_IN_ROUTE + 1 = TMP8.TMP8_STATION_INDEX_IN_ROUTE INNER JOIN ( SELECT NEIGHBOUR_ID ,DISTANCE ,STATIONS_ID AS 'TMP9_STATIONS_ID' FROM NEIGHBOURS ) TMP9 ON TMP9.NEIGHBOUR_ID = TMP8.TMP8_STATIONS_ID AND TMP9.TMP9_STATIONS_ID = TMP7.TMP7_STATIONS_ID GROUP BY ID2 ,ID3 ) TMP10 INNER JOIN ( SELECT TRAINS_ID AS 'TMP11_TRAINS_ID' ,CASE WHEN ISNULL(GROUP_CONCAT(CASE WHEN CLASS = 1 THEN SERVICES END separator ';')) THEN 0 ELSE GROUP_CONCAT(CASE WHEN CLASS = 1 THEN SERVICES END separator ';') END AS 'CLASS1_SERVICES1' FROM WAGONS GROUP BY WAGONS.TRAINS_ID ) TMP11 ON TMP10.FIRST_TRAIN_ID = TMP11.TMP11_TRAINS_ID INNER JOIN ( SELECT TRAINS_ID AS 'TMP12_TRAINS_ID' ,CASE WHEN ISNULL(GROUP_CONCAT(CASE WHEN CLASS = 2 THEN SERVICES END separator ';')) THEN 0	ELSE GROUP_CONCAT(CASE WHEN CLASS = 2	THEN SERVICES	END separator ';') END AS 'CLASS2_SERVICES1' FROM WAGONS GROUP BY WAGONS.TRAINS_ID ) TMP12 ON TMP10.FIRST_TRAIN_ID = TMP12.TMP12_TRAINS_ID INNER JOIN ( SELECT TRAINS_ID AS 'TMP13_TRAINS_ID' ,CASE WHEN ISNULL(GROUP_CONCAT(CASE WHEN CLASS = 1 THEN SERVICES END separator ';')) THEN 0 ELSE GROUP_CONCAT(CASE WHEN CLASS = 1 THEN SERVICES END separator ';') END AS 'CLASS1_SERVICES2' FROM WAGONS GROUP BY WAGONS.TRAINS_ID ) TMP13 ON TMP10.SECOND_TRAIN_ID = TMP13.TMP13_TRAINS_ID INNER JOIN ( SELECT TRAINS_ID AS 'TMP14_TRAINS_ID' ,CASE  WHEN ISNULL(GROUP_CONCAT(CASE WHEN CLASS = 2 THEN SERVICES END separator ';')) THEN 0 ELSE GROUP_CONCAT(CASE WHEN CLASS = 2 THEN SERVICES END separator ';') END AS 'CLASS2_SERVICES2' FROM WAGONS GROUP BY WAGONS.TRAINS_ID ) TMP14 ON TMP10.SECOND_TRAIN_ID = TMP14.TMP14_TRAINS_ID WHERE DATE (FIRST_TRAIN_START + INTERVAL DISTANCE_BEFORE_START1 / FIRST_TRAIN_SPEED * 60.0 MINUTE) = DATE ('" + date + "') AND SECOND_TRAIN_START + INTERVAL DISTANCE_BEFORE_START2 / SECOND_TRAIN_SPEED * 60.0 MINUTE > FIRST_TRAIN_START + INTERVAL(DISTANCE_BEFORE_START1 + DISTANCE_BETWEEN_START_AND_STOP1) / SECOND_TRAIN_SPEED * 60.0 MINUTE ) TMP_ATSZALLAS INNER JOIN ( SELECT * ,EXPLETIVE_TICKET AS 'EXPLETIVE_TICKET1' FROM ROUTE_STATIONS_CONNECTION ) TMP15 ON TMP15.ROUTES_ID = ID2 INNER JOIN ( SELECT * ,EXPLETIVE_TICKET AS 'EXPLETIVE_TICKET2' FROM ROUTE_STATIONS_CONNECTION ) TMP16 ON TMP16.ROUTES_ID = ID3 WHERE " + ch + " = 1 AND TMP15.STATION_INDEX_IN_ROUTE BETWEEN TMP_ATSZALLAS.START_STATION_INDEX_IN_ROUTE AND TMP_ATSZALLAS.MIDDLE_STATION_INDEX_IN_ROUTE1 AND TMP16.STATION_INDEX_IN_ROUTE BETWEEN TMP_ATSZALLAS.MIDDLE_STATION_INDEX_IN_ROUTE2 AND TMP_ATSZALLAS.STOP_STATION_INDEX_IN_ROUTE2  GROUP BY ID2,ID3 ) TMP_ALL";

	if extraTicket {
			distanceAndTimeQuery = distanceAndTimeQuery + " WHERE FIRST_TRAIN_EXPLETIVE_TICKET = 0 AND SECOND_TRAIN_EXPLETIVE_TICKET = 0"
	}
	distanceAndTimeQuery = distanceAndTimeQuery + " ORDER BY START_TIME1"

	result = make(SQLDataSlice, 0)

	inf := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", this.username, this.password, this.host, this.port, this.db)
	db, err := sql.Open("mysql", inf)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(distanceAndTimeQuery)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&d1, &d2, &d3, &d4, &d5, &d6, &d7, &d8, &d9, &d10, &d11, &d12, &d13, &d14, &d15, &d16, &d17, &d18, &d19, &d20, &d21, &d22)
		if err != nil {
			panic(err)
		}

		var trainInfo TrainInformation
		changes := d1
		var date1 string
		if changes == "0" {
			date1 = d11
		} else {
			date1 = d20
		}
		departure := d10[11:19]
		arrival := date1[11:19]

		d, err := time.Parse("2006-01-02 15:04:00", d10)
		if err != nil {
			panic(err)
		}

		a, err := time.Parse("2006-01-02 15:04:00", date1)
		if err != nil {
			panic(err)
		}

		delta := a.Sub(d)
		h := int(delta.Seconds() / 3600)
		min := (int(delta.Seconds()) % 3600) / 60
		sec := (int(delta.Seconds()) % 3600) % 60
		var duration string
		if h < 10 {
			h := "0" + strconv.Itoa(h)
			duration = h + ":"
		} else {
			duration = fmt.Sprintf("%d:", h)
		}
		if min < 10 {
			min := "0" + strconv.Itoa(min)
			duration = duration + min + ":"
		} else {
			duration = duration + fmt.Sprintf("%d:", min)
		}
		if sec < 10 {
			sec := "0" + strconv.Itoa(sec)
			duration = duration + sec
		} else {
			duration = duration + fmt.Sprintf("%d", sec)
		}

		km1, err := strconv.Atoi(d5)
		km2, err := strconv.Atoi(d15)
		km := km1 + km2
		percent, err := strconv.Atoi(discount)

		var extra int
		if d9 != "0" || d18 != "0" { //potjegy
			extra = 300
		}
		if d4 == "IC" || d14 == "IC" { //helyijegy
			extra = extra + 250
		}

		var priceNUM float64
		priceNUM = (float64(km)*float64(pricePerKM) + float64(extra)) * ((100.0 - float64(percent)) / 100.0)
		price := strconv.FormatFloat(priceNUM, 'f', 0, 64)

		/*belso adatok*/
		trainInfo.Station = append(trainInfo.Station, from)
		if changes == "1" {
			trainInfo.Station = append(trainInfo.Station, d6)
			trainInfo.Station = append(trainInfo.Station, d6)
		}
		trainInfo.Station = append(trainInfo.Station, to)
		trainInfo.Timetable = append(trainInfo.Timetable, departure)
		if changes == "1" {
			trainInfo.Timetable = append(trainInfo.Timetable, d11[11:19])
			trainInfo.Timetable = append(trainInfo.Timetable, d19[11:19])
		}
		trainInfo.Timetable = append(trainInfo.Timetable, arrival)
		trainInfo.Train = append(trainInfo.Train, d2)
		if changes == "1" {
			trainInfo.Train = append(trainInfo.Train, d12)
		}
		trainInfo.RouteID = append(trainInfo.RouteID, d21)
		if changes == "1" {
			trainInfo.RouteID = append(trainInfo.RouteID, d22)
		}

		var service TrainServices
		str := d7
		for i := 0; i < len(str); i += 2 {
			if str[i] == '1' {
				service.S1_Toilet = true
			}
			if str[i] == '2' {
				service.S1_DisabledToilet = true
			}
			if str[i] == '3' {
				service.S1_DiaperChange = true
			}
			if str[i] == '4' {
				service.S1_AirConditioner = true
			}
			if str[i] == '5' {
				service.S1_Wifi = true
			}
			if str[i] == '6' {
				service.S1_PowerConnector = true
			}
			if str[i] == '7' {
				service.S1_Restaurant = true
			}
			if str[i] == '8' {
				service.S1_BikeShed = true
			}
			if str[i] == '9' {
				service.S1_Bed = true
			}
		}
		str = d8
		for i := 0; i < len(str); i += 2 {
			if str[i] == '1' {
				service.S2_Toilet = true
			}
			if str[i] == '2' {
				service.S2_DisabledToilet = true
			}
			if str[i] == '3' {
				service.S2_DiaperChange = true
			}
			if str[i] == '4' {
				service.S2_AirConditioner = true
			}
			if str[i] == '5' {
				service.S2_Wifi = true
			}
			if str[i] == '6' {
				service.S2_PowerConnector = true
			}
			if str[i] == '7' {
				service.S2_Restaurant = true
			}
			if str[i] == '8' {
				service.S2_BikeShed = true
			}
			if str[i] == '9' {
				service.S2_Bed = true
			}
		}
		trainInfo.Services = append(trainInfo.Services, service)

		if changes == "1" {
			var service2 TrainServices
			str := d16
			for i := 0; i < len(str); i += 2 {
				if str[i] == '1' {
					service2.S1_Toilet = true
				}
				if str[i] == '2' {
					service2.S1_DisabledToilet = true
				}
				if str[i] == '3' {
					service2.S1_DiaperChange = true
				}
				if str[i] == '4' {
					service2.S1_AirConditioner = true
				}
				if str[i] == '5' {
					service2.S1_Wifi = true
				}
				if str[i] == '6' {
					service2.S1_PowerConnector = true
				}
				if str[i] == '7' {
					service2.S1_Restaurant = true
				}
				if str[i] == '8' {
					service2.S1_BikeShed = true
				}
				if str[i] == '9' {
					service2.S1_Bed = true
				}
			}
			str = d17
			for i := 0; i < len(str); i += 2 {
				if str[i] == '1' {
					service2.S2_Toilet = true
				}
				if str[i] == '2' {
					service2.S2_DisabledToilet = true
				}
				if str[i] == '3' {
					service2.S2_DiaperChange = true
				}
				if str[i] == '4' {
					service2.S2_AirConditioner = true
				}
				if str[i] == '5' {
					service2.S2_Wifi = true
				}
				if str[i] == '6' {
					service2.S2_PowerConnector = true
				}
				if str[i] == '7' {
					service2.S2_Restaurant = true
				}
				if str[i] == '8' {
					service2.S2_BikeShed = true
				}
				if str[i] == '9' {
					service2.S2_Bed = true
				}
			}
			trainInfo.Services = append(trainInfo.Services, service2)
		}

		data.From = from
		data.To = to
		data.Date = date
		data.Departure = departure
		data.Arrival = arrival
		data.Changes = changes
		data.Duration = duration
		data.Distance = strconv.Itoa(km)
		data.Price = price
		data.Info = trainInfo

		result = append(result, data)
	}
	finalResult.From = from
	finalResult.To = to
	finalResult.Data = result
	return finalResult
}

func (this SQLConfig) MysqlListStationsByRouteID(from, to, departure, arrival, route_id string) MapData {

	var d string
	var mapData MapData
	query := "SELECT NAME FROM ROUTE_STATIONS_CONNECTION INNER JOIN STATIONS ON STATIONS_ID = ID WHERE ROUTES_ID = " + route_id + " ORDER BY STATION_INDEX_IN_ROUTE"
	inf := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", this.username, this.password, this.host, this.port, this.db)
	db, err := sql.Open("mysql", inf)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&d)
		if err != nil {
			panic(err)
		}
		mapData.Stations = append(mapData.Stations, d)
	}

	mapData.From = from
	mapData.To = to
	mapData.Departure = departure
	mapData.Arrival = arrival

	return mapData
}

func (this SQLConfig) MysqlGetTrainType(from1, to1, departure1, arrival1, train1ID, from2, to2, departure2, arrival2, train2ID, price, km string) TrainType {
	query1 := "SELECT TYPE FROM TRAINS WHERE ID = " + train1ID
	query2 := "SELECT TYPE FROM TRAINS WHERE ID = " + train2ID
	var d string
	var result TrainType
	inf := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", this.username, this.password, this.host, this.port, this.db)
	db, err := sql.Open("mysql", inf)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(query1)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&d)
		if err != nil {
			panic(err)
		}
		if d == "IC" {
			result.Type1 = true
		} else {
			result.Type1 = false
		}
	}

	result.From1 = from1
	result.To1 = to1
	result.Departure1 = departure1
	result.Arrival1 = arrival1
	result.Train1ID = train1ID
	result.SeatNumber1 = "0"
	result.Price = price
	result.Km = km

	result.From2 = from2
	result.To2 = to2
	result.Departure2 = departure2
	result.Arrival2 = arrival2
	result.Train2ID = train2ID
	result.SeatNumber2 = "0"
	if from2 == "0" {
		result.Type2 = false
	} else {
		rows, err := db.Query(query2)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&d)
			if err != nil {
				panic(err)
			}
			if d == "IC" {
				result.Type2 = true
			} else {
				result.Type2 = false
			}
		}
	}
	return result
}

func (this SQLConfig) MysqlSeatReserve(id, from1, to1, departure1, arrival1, train1ID, from2, to2, departure2, arrival2, train2ID, price, km, seat1, seat2 string) WagonData {

	var result WagonData
	var wagon Wagon
	wagon.Seats = make(map[string]bool)
	query := "SELECT SEATS_NUMBER, CLASS, SERVICES, WAGONS_ID, RESERVED ,NUMBER FROM WAGONS INNER JOIN WAGON_SEAT_CONNECTION ON WAGONS_ID = WAGONS.ID INNER JOIN SEATS ON SEATS.ID = WAGON_SEAT_CONNECTION.SEATS_ID WHERE TRAINS_ID = " + id
	var numOfSeats, class, services, wagonId, reserved, seatNumber string
	inf := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", this.username, this.password, this.host, this.port, this.db)
	db, err := sql.Open("mysql", inf)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	wagonID := ""
	for rows.Next() {
		err := rows.Scan(&numOfSeats, &class, &services, &wagonId, &reserved, &seatNumber)
		if err != nil {
			panic(err)
		}
		if wagonID == "" {
			wagonID = wagonId
			wagon.WagonID = wagonID
			wagon.NumberOfSeats = numOfSeats
			wagon.Class = class

			str := services
			for i := 0; i < len(str); i += 2 {
				if str[i] == '1' {
					wagon.Services.Toilet = true
				}
				if str[i] == '2' {
					wagon.Services.DisabledToilet = true
				}
				if str[i] == '3' {
					wagon.Services.DiaperChange = true
				}
				if str[i] == '4' {
					wagon.Services.AirConditioner = true
				}
				if str[i] == '5' {
					wagon.Services.Wifi = true
				}
				if str[i] == '6' {
					wagon.Services.PowerConnector = true
				}
				if str[i] == '7' {
					wagon.Services.Restaurant = true
				}
				if str[i] == '8' {
					wagon.Services.BikeShed = true
				}
				if str[i] == '9' {
					wagon.Services.Bed = true
				}
			}
		}
		if wagonID == wagonId {
			if reserved == "1" {
				wagon.Seats[seatNumber] = true
			} else {
				wagon.Seats[seatNumber] = false
			}
		} else {
			result.Wagons = append(result.Wagons, wagon)

			wagon.Seats = make(map[string]bool, 0)
			wagon.Services.Toilet = false
			wagon.Services.DisabledToilet = false
			wagon.Services.DiaperChange = false
			wagon.Services.AirConditioner = false
			wagon.Services.Wifi = false
			wagon.Services.PowerConnector = false
			wagon.Services.Restaurant = false
			wagon.Services.BikeShed = false
			wagon.Services.Bed = false

			wagonID = wagonId
			wagon.WagonID = wagonID
			wagon.NumberOfSeats = numOfSeats
			wagon.Class = class
			str := services
			for i := 0; i < len(str); i += 2 {
				if str[i] == '1' {
					wagon.Services.Toilet = true
				}
				if str[i] == '2' {
					wagon.Services.DisabledToilet = true
				}
				if str[i] == '3' {
					wagon.Services.DiaperChange = true
				}
				if str[i] == '4' {
					wagon.Services.AirConditioner = true
				}
				if str[i] == '5' {
					wagon.Services.Wifi = true
				}
				if str[i] == '6' {
					wagon.Services.PowerConnector = true
				}
				if str[i] == '7' {
					wagon.Services.Restaurant = true
				}
				if str[i] == '8' {
					wagon.Services.BikeShed = true
				}
				if str[i] == '9' {
					wagon.Services.Bed = true
				}
			}

			if reserved == "1" {
				wagon.Seats[seatNumber] = true
			} else {
				wagon.Seats[seatNumber] = false
			}
		}
	}
	result.Wagons = append(result.Wagons, wagon)

	query1 := "SELECT TYPE FROM TRAINS WHERE ID = " + train1ID
	query2 := "SELECT TYPE FROM TRAINS WHERE ID = " + train2ID
	var d string
	var traintype TrainType

	rows1, err := db.Query(query1)
	if err != nil {
		panic(err)
	}
	defer rows1.Close()

	for rows1.Next() {
		err := rows1.Scan(&d)
		if err != nil {
			panic(err)
		}
		if d == "IC" {
			traintype.Type1 = true
		} else {
			traintype.Type1 = false
		}
	}

	traintype.From1 = from1
	traintype.To1 = to1
	traintype.Departure1 = departure1
	traintype.Arrival1 = arrival1
	traintype.Train1ID = train1ID
	traintype.SeatNumber1 = seat1
	traintype.Price = price
	traintype.Km = km

	traintype.From2 = from2
	traintype.To2 = to2
	traintype.Departure2 = departure2
	traintype.Arrival2 = arrival2
	traintype.Train2ID = train2ID
	traintype.SeatNumber2 = seat2
	if from2 == "0" {
		traintype.Type2 = false
	} else {
		rows2, err := db.Query(query2)
		if err != nil {
			panic(err)
		}
		defer rows2.Close()

		for rows2.Next() {
			err = rows2.Scan(&d)
			if err != nil {
				panic(err)
			}
			if d == "IC" {
				traintype.Type2 = true
			} else {
				traintype.Type2 = false
			}
		}
	}
	result.Train = traintype
	result.SelectedTrain = id
	return result
}

func (this SQLConfig) MysqlCheckReservation(wagonID, seat string) bool {
	query := "SELECT RESERVED FROM SEATS INNER JOIN WAGON_SEAT_CONNECTION ON ID = SEATS_ID INNER JOIN WAGONS ON WAGONS.ID = WAGON_SEAT_CONNECTION.WAGONS_ID WHERE WAGONS.ID = '" + wagonID + "' AND SEATS.ID = " + seat
	var reserved string
	inf := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", this.username, this.password, this.host, this.port, this.db)
	db, err := sql.Open("mysql", inf)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&reserved)
		if err != nil {
			panic(err)
		}
		if reserved == "1" {
			return true
		}
	}
	return false
}

func (this SQLConfig) MysqlUpdateWagonReservation(wagonID, seat, from1, to1, departure1, arrival1, train1ID, from2, to2, departure2, arrival2, train2ID, price, km, seat1, seat2, selectedTrain string) TrainType {
	id := strings.Index(wagonID, " ")
	wagonId := wagonID[0:id]

	query := "UPDATE SEATS INNER JOIN WAGON_SEAT_CONNECTION ON SEATS.ID = SEATS_ID INNER JOIN WAGONS ON WAGONS.ID = WAGONS_ID SET SEATS.RESERVED=1 WHERE SEATS.NUMBER=" + seat + " AND WAGONS.ID = '" + wagonId + "'"
	inf := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", this.username, this.password, this.host, this.port, this.db)
	db, err := sql.Open("mysql", inf)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(query)

	if err != nil {
		panic(err)
	}

	query1 := "SELECT TYPE FROM TRAINS WHERE ID = " + train1ID
	query2 := "SELECT TYPE FROM TRAINS WHERE ID = " + train2ID
	var d string
	var result TrainType

	rows, err := db.Query(query1)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&d)
		if err != nil {
			panic(err)
		}
		if d == "IC" {
			result.Type1 = true
		} else {
			result.Type1 = false
		}
	}

	result.From1 = from1
	result.To1 = to1
	result.Departure1 = departure1
	result.Arrival1 = arrival1
	result.Train1ID = train1ID
	if selectedTrain == train1ID {
		result.SeatNumber1 = wagonId + " / " + seat
	} else {
		result.SeatNumber1 = seat1
	}
	result.Price = price
	result.Km = km

	result.From2 = from2
	result.To2 = to2
	result.Departure2 = departure2
	result.Arrival2 = arrival2
	result.Train2ID = train2ID
	if selectedTrain == train2ID {
		result.SeatNumber2 = wagonId + " / " + seat
	} else {
		result.SeatNumber2 = seat2
	}
	if from2 == "0" {
		result.Type2 = false
	} else {
		rows, err := db.Query(query2)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&d)
			if err != nil {
				panic(err)
			}
			if d == "IC" {
				result.Type2 = true
			} else {
				result.Type2 = false
			}
		}
	}
	return result
}

func (this SQLConfig) MysqlBuyTicket(firstname, lastname, from1, to1, departure1, arrival1, train1ID, seat1, from2, to2, departure2, arrival2, train2ID, seat2, price, km string) PurchaseInfo {
	var result PurchaseInfo
	passw := generateSalt()
	query1 := "INSERT INTO PURCHASES (FIRSTNAME, LASTNAME, FROM1, TO1, DEPARTURE1, ARRIVAL1, TRAIN1_ID, SEAT1, FROM2, TO2, DEPARTURE2, ARRIVAL2, TRAIN2_ID, SEAT2, PRICE, KM, PASSW) VALUES ('"+ firstname +"', '"+ lastname +"', '"+ from1 +"', '"+ to1 +"', '"+ departure1 +"', '"+ arrival1 +"', '"+ train1ID +"', '"+ seat1 +"', '"+ from2 +"', '"+ to2 +"', '"+ departure2 +"', '"+ arrival2 +"', '"+ train2ID +"', '"+ seat2 +"', '"+ price +"', '"+ km +"', '"+passw+"')"


	inf := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", this.username, this.password, this.host, this.port, this.db)
	db, err := sql.Open("mysql", inf)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Exec(query1)
	id, err := res.LastInsertId()

	if err != nil {
		panic(err)
	}

	result.TicketID = strconv.FormatInt(id, 10)
	result.TicketPassw = passw
	return result
}

func (this SQLConfig) MysqlCheckTicket(tickedID, passw string) bool {
	query := "SELECT PASSW FROM PURCHASES WHERE ID="+tickedID
	var password string
	inf := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", this.username, this.password, this.host, this.port, this.db)
	db, err := sql.Open("mysql", inf)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&password)
		if err != nil {
			panic(err)
		}
		if password == passw {
			return true
		}
	}
	return false
}

func (this SQLConfig) MysqlSetTicketInformation(id string) Ticket {
	var result Ticket
	query := "SELECT * FROM PURCHASES WHERE ID= "+id
	var firstname, lastname, date, from1, to1, departure1, arrival1, train1ID, seat1, from2, to2, departure2, arrival2, train2ID, seat2, price, km, ticketID, passw string
	inf := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", this.username, this.password, this.host, this.port, this.db)
	db, err := sql.Open("mysql", inf)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&ticketID, &firstname, &lastname, &from1, &to1, &departure1, &arrival1, &train1ID, &seat1, &from2, &to2, &departure2, &arrival2, &train2ID, &seat2, &price, &km, &passw, &date)
		if err != nil {
			panic(err)
		}
		result.Firstname = firstname
		result.Lastname = lastname
		result.Date = date[0:10]
		result.From1 = from1
		result.To1 = to1
		result.Departure1 = departure1
		result.Arrival1 = arrival1
		result.Train1ID = train1ID
		result.Seat1 = seat1
		result.From2 = from2
		result.To2 = to2
		result.Departure2 = departure2
		result.Arrival2 = arrival2
		result.Train2ID = train2ID
		result.Seat2 = seat2
		result.Price = price
		result.Km = km
		result.TicketID = ticketID
		result.TicketPassw = passw
	}

	return result
}
