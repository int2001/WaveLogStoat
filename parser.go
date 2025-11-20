package main

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// WSJT-X XML structure
type WSJTContactInfo struct {
	XMLName    xml.Name `xml:"contactinfo"`
	Timestamp  string   `xml:"timestamp"`
	Call       string   `xml:"call"`
	Mode       string   `xml:"mode"`
	TxFreq     string   `xml:"txfreq"`
	RxFreq     string   `xml:"rxfreq"`
	Rcv        string   `xml:"rcv"`
	Snt        string   `xml:"snt"`
	Power      string   `xml:"power"`
	Operator   string   `xml:"operator"`
	Comment    string   `xml:"comment"`
	Sntnr      string   `xml:"sntnr"`
	Rcvnr      string   `xml:"rcvnr"`
	MyCall     string   `xml:"mycall"`
	Gridsquare string   `xml:"gridsquare"`
}

func parseXMLMessage(message string) (QSO, error) {
	var contactInfo WSJTContactInfo
	if err := xml.Unmarshal([]byte(message), &contactInfo); err != nil {
		return QSO{}, fmt.Errorf("XML parsing failed: %v", err)
	}

	// Parse timestamp
	timestamp, err := time.Parse("2006-01-02T15:04:05", contactInfo.Timestamp)
	if err != nil {
		return QSO{}, fmt.Errorf("timestamp parsing failed: %v", err)
	}

	// Convert mode for TCADIF compatibility
	mode := contactInfo.Mode
	if mode == "USB" || mode == "LSB" {
		mode = "SSB"
	}

	// Convert frequency from Hz to MHz
	txFreq, err := strconv.ParseFloat(contactInfo.TxFreq, 64)
	if err != nil {
		return QSO{}, fmt.Errorf("TX frequency parsing failed: %v", err)
	}
	freqMHz := txFreq / 100000

	rxFreq, err := strconv.ParseFloat(contactInfo.RxFreq, 64)
	if err != nil {
		return QSO{}, fmt.Errorf("RX frequency parsing failed: %v", err)
	}
	freqRXMHz := rxFreq / 100000

	qso := QSO{
		CALL:             contactInfo.Call,
		MODE:             mode,
		QSO_DATE_OFF:     timestamp.Format("20060102"),
		QSO_DATE:         timestamp.Format("20060102"),
		TIME_OFF:         timestamp.Format("150405"),
		TIME_ON:          timestamp.Format("150405"),
		RST_RCVD:         contactInfo.Rcv,
		RST_SENT:         contactInfo.Snt,
		FREQ:             fmt.Sprintf("%.6f", freqMHz),
		FREQ_RX:          fmt.Sprintf("%.6f", freqRXMHz),
		OPERATOR:         contactInfo.Operator,
		COMMENT:          contactInfo.Comment,
		POWER:            contactInfo.Power,
		STX:              contactInfo.Sntnr,
		RTX:              contactInfo.Rcvnr,
		MYCALL:           contactInfo.MyCall,
		GRIDSQUARE:       contactInfo.Gridsquare,
		STATION_CALLSIGN: contactInfo.MyCall,
	}

	if verbose {
		logger.Printf("Parsed XML QSO: %s on %s MHz", qso.CALL, qso.FREQ)
	}

	return qso, nil
}

func parseADIFMessage(message string) (QSO, error) {
	qso := QSO{}

	// Simple ADIF parser - extract fields
	// ADIF format: <FIELD:length:data>
	re := regexp.MustCompile(`<([a-zA-Z_]+):(\d+)>`)
	matches := re.FindAllStringSubmatch(message, -1)

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		field := match[1]
		lengthStr := match[2]

		// Extract the data immediately after the field tag
		fieldStart := strings.Index(message, match[0]) + len(match[0])
		if fieldStart >= len(message) {
			continue
		}

		// Parse the length and extract that many characters
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			continue
		}

		// Extract the data of specified length
		fieldEnd := fieldStart + length
		if fieldEnd > len(message) {
			fieldEnd = len(message)
		}

		data := strings.TrimSpace(message[fieldStart:fieldEnd])

		// Map ADIF fields to QSO structure
		switch strings.ToUpper(field) {
		case "CALL":
			qso.CALL = data
		case "MODE":
			qso.MODE = data
		case "QSO_DATE_OFF":
			qso.QSO_DATE_OFF = data
			qso.QSO_DATE = data
		case "QSO_DATE":
			qso.QSO_DATE = data
		case "TIME_OFF":
			qso.TIME_OFF = data
			qso.TIME_ON = data
		case "TIME_ON":
			qso.TIME_ON = data
		case "RST_RCVD":
			qso.RST_RCVD = data
		case "RST_SENT":
			qso.RST_SENT = data
		case "FREQ":
			qso.FREQ = data
		case "FREQ_RX":
			qso.FREQ_RX = data
		case "OPERATOR":
			qso.OPERATOR = data
		case "COMMENT":
			qso.COMMENT = data
		case "TX_PWR":
			qso.POWER = data
		case "STX":
			qso.STX = data
		case "SRX":
			qso.SRX = data
		case "STX_STRING":
			qso.STX_STRING = data
		case "SRX_STRING":
			qso.SRX_STRING = data
		case "RTX":
			qso.RTX = data
		case "CONTEST_ID":
			qso.CONTEST_ID = data
		case "PREFIX":
			qso.PREFIX = data
		case "SUBMODE":
			qso.SUBMODE = data
		case "QSLMSG":
			qso.QSLMSG = data
		case "NOTES":
			qso.NOTES = data
		case "EMAIL":
			qso.EMAIL = data
		case "DARC_DOK":
			qso.DARC_DOK = data
		case "SOTA_REF":
			qso.SOTA_REF = data
		case "WWFF_REF":
			qso.WWFF_REF = data
		case "POTA_REF":
			qso.POTA_REF = data
		case "CNTY":
			qso.CNTY = data
		case "REGION":
			qso.REGION = data
		case "LAT":
			qso.LAT = data
		case "LON":
			qso.LON = data
		case "ANT_AZ":
			qso.ANT_AZ = data
		case "ANT_EL":
			qso.ANT_EL = data
		case "ANT_PATH":
			qso.ANT_PATH = data
		case "A_INDEX":
			qso.A_INDEX = data
		case "K_INDEX":
			qso.K_INDEX = data
		case "SFI":
			qso.SFI = data
		case "RX_PWR":
			qso.RX_PWR = data
		case "MY_CALL":
			qso.MYCALL = data
			qso.STATION_CALLSIGN = data
		case "MY_GRIDSQUARE":
			qso.MY_GRIDSQUARE = data
		case "NAME":
			qso.NAME = data
		case "QTH":
			qso.QTH = data
		case "STATE":
			qso.STATE = data
		case "COUNTRY":
			qso.COUNTRY = data
		case "CQZ":
			qso.CQZ = data
		case "ITUZ":
			qso.ITUZ = data
		case "CONT":
			qso.CONT = data
		case "IOTA":
			qso.IOTA = data
		case "DXCC":
			qso.DXCC = data
		case "PROP_MODE":
			qso.PROP_MODE = data
		case "SAT_NAME":
			qso.SAT_NAME = data
		case "SAT_MODE":
			qso.SAT_MODE = data
		case "GRIDSQUARE":
			qso.GRIDSQUARE = data
		case "STATION_CALLSIGN":
			qso.STATION_CALLSIGN = data
		}
	}

	// Validate required fields
	if qso.CALL == "" {
		return QSO{}, fmt.Errorf("missing required CALL field in ADIF")
	}

	if verbose {
		logger.Printf("Parsed ADIF QSO: %s on %s MHz", qso.CALL, qso.FREQ)
	}

	return qso, nil
}

func generateADIF(qso QSO) string {
	var adif strings.Builder

	// Add ADIF header if needed
	adif.WriteString("<ADIF_VER:5>5.0<EOH>\n")

	// Add QSO fields
	if qso.CALL != "" {
		adif.WriteString(fmt.Sprintf("<CALL:%d>%s ", len(qso.CALL), qso.CALL))
	}
	if qso.QSO_DATE != "" {
		adif.WriteString(fmt.Sprintf("<QSO_DATE:%d>%s ", len(qso.QSO_DATE), qso.QSO_DATE))
	}
	if qso.TIME_ON != "" {
		adif.WriteString(fmt.Sprintf("<TIME_ON:%d>%s ", len(qso.TIME_ON), qso.TIME_ON))
	}
	if qso.MODE != "" {
		adif.WriteString(fmt.Sprintf("<MODE:%d>%s ", len(qso.MODE), qso.MODE))
	}
	if qso.RST_RCVD != "" {
		adif.WriteString(fmt.Sprintf("<RST_RCVD:%d>%s ", len(qso.RST_RCVD), qso.RST_RCVD))
	}
	if qso.RST_SENT != "" {
		adif.WriteString(fmt.Sprintf("<RST_SENT:%d>%s ", len(qso.RST_SENT), qso.RST_SENT))
	}
	if qso.FREQ != "" {
		adif.WriteString(fmt.Sprintf("<FREQ:%d>%s ", len(qso.FREQ), qso.FREQ))
	}
	if qso.FREQ_RX != "" {
		adif.WriteString(fmt.Sprintf("<FREQ_RX:%d>%s ", len(qso.FREQ_RX), qso.FREQ_RX))
	}
	if qso.BAND != "" {
		adif.WriteString(fmt.Sprintf("<BAND:%d>%s ", len(qso.BAND), qso.BAND))
	}
	if qso.POWER != "" {
		adif.WriteString(fmt.Sprintf("<TX_PWR:%d>%s ", len(qso.POWER), qso.POWER))
	}
	if qso.OPERATOR != "" {
		adif.WriteString(fmt.Sprintf("<OPERATOR:%d>%s ", len(qso.OPERATOR), qso.OPERATOR))
	}
	if qso.MYCALL != "" {
		adif.WriteString(fmt.Sprintf("<MY_CALL:%d>%s ", len(qso.MYCALL), qso.MYCALL))
	}
	if qso.STATION_CALLSIGN != "" {
		adif.WriteString(fmt.Sprintf("<STATION_CALLSIGN:%d>%s ", len(qso.STATION_CALLSIGN), qso.STATION_CALLSIGN))
	}
	if qso.GRIDSQUARE != "" {
		adif.WriteString(fmt.Sprintf("<GRIDSQUARE:%d>%s ", len(qso.GRIDSQUARE), qso.GRIDSQUARE))
	}
	if qso.COMMENT != "" {
		adif.WriteString(fmt.Sprintf("<COMMENT:%d>%s ", len(qso.COMMENT), qso.COMMENT))
	}
	if qso.STX != "" {
		adif.WriteString(fmt.Sprintf("<STX:%d>%s ", len(qso.STX), qso.STX))
	}
	if qso.SRX != "" {
		adif.WriteString(fmt.Sprintf("<SRX:%d>%s ", len(qso.SRX), qso.SRX))
	}
	if qso.STX_STRING != "" {
		adif.WriteString(fmt.Sprintf("<STX_STRING:%d>%s ", len(qso.STX_STRING), qso.STX_STRING))
	}
	if qso.SRX_STRING != "" {
		adif.WriteString(fmt.Sprintf("<SRX_STRING:%d>%s ", len(qso.SRX_STRING), qso.SRX_STRING))
	}
	if qso.RTX != "" {
		adif.WriteString(fmt.Sprintf("<RTX:%d>%s ", len(qso.RTX), qso.RTX))
	}
	// ADIF-compliant contest fields
	if qso.CONTEST_ID != "" {
		adif.WriteString(fmt.Sprintf("<CONTEST_ID:%d>%s ", len(qso.CONTEST_ID), qso.CONTEST_ID))
	}
	if qso.PREFIX != "" {
		adif.WriteString(fmt.Sprintf("<PREFIX:%d>%s ", len(qso.PREFIX), qso.PREFIX))
	}
	if qso.MY_GRIDSQUARE != "" {
		adif.WriteString(fmt.Sprintf("<MY_GRIDSQUARE:%d>%s ", len(qso.MY_GRIDSQUARE), qso.MY_GRIDSQUARE))
	}
	if qso.NAME != "" {
		adif.WriteString(fmt.Sprintf("<NAME:%d>%s ", len(qso.NAME), qso.NAME))
	}
	if qso.QTH != "" {
		adif.WriteString(fmt.Sprintf("<QTH:%d>%s ", len(qso.QTH), qso.QTH))
	}
	if qso.STATE != "" {
		adif.WriteString(fmt.Sprintf("<STATE:%d>%s ", len(qso.STATE), qso.STATE))
	}
	if qso.COUNTRY != "" {
		adif.WriteString(fmt.Sprintf("<COUNTRY:%d>%s ", len(qso.COUNTRY), qso.COUNTRY))
	}
	if qso.CQZ != "" {
		adif.WriteString(fmt.Sprintf("<CQZ:%d>%s ", len(qso.CQZ), qso.CQZ))
	}
	if qso.ITUZ != "" {
		adif.WriteString(fmt.Sprintf("<ITUZ:%d>%s ", len(qso.ITUZ), qso.ITUZ))
	}
	if qso.CONT != "" {
		adif.WriteString(fmt.Sprintf("<CONT:%d>%s ", len(qso.CONT), qso.CONT))
	}
	if qso.IOTA != "" {
		adif.WriteString(fmt.Sprintf("<IOTA:%d>%s ", len(qso.IOTA), qso.IOTA))
	}
	if qso.DXCC != "" {
		adif.WriteString(fmt.Sprintf("<DXCC:%d>%s ", len(qso.DXCC), qso.DXCC))
	}
	if qso.PROP_MODE != "" {
		adif.WriteString(fmt.Sprintf("<PROP_MODE:%d>%s ", len(qso.PROP_MODE), qso.PROP_MODE))
	}
	if qso.SAT_NAME != "" {
		adif.WriteString(fmt.Sprintf("<SAT_NAME:%d>%s ", len(qso.SAT_NAME), qso.SAT_NAME))
	}
	if qso.SAT_MODE != "" {
		adif.WriteString(fmt.Sprintf("<SAT_MODE:%d>%s ", len(qso.SAT_MODE), qso.SAT_MODE))
	}
	if qso.SUBMODE != "" {
		adif.WriteString(fmt.Sprintf("<SUBMODE:%d>%s ", len(qso.SUBMODE), qso.SUBMODE))
	}
	if qso.QSLMSG != "" {
		adif.WriteString(fmt.Sprintf("<QSLMSG:%d>%s ", len(qso.QSLMSG), qso.QSLMSG))
	}
	if qso.NOTES != "" {
		adif.WriteString(fmt.Sprintf("<NOTES:%d>%s ", len(qso.NOTES), qso.NOTES))
	}
	if qso.EMAIL != "" {
		adif.WriteString(fmt.Sprintf("<EMAIL:%d>%s ", len(qso.EMAIL), qso.EMAIL))
	}
	if qso.DARC_DOK != "" {
		adif.WriteString(fmt.Sprintf("<DARC_DOK:%d>%s ", len(qso.DARC_DOK), qso.DARC_DOK))
	}
	if qso.SOTA_REF != "" {
		adif.WriteString(fmt.Sprintf("<SOTA_REF:%d>%s ", len(qso.SOTA_REF), qso.SOTA_REF))
	}
	if qso.WWFF_REF != "" {
		adif.WriteString(fmt.Sprintf("<WWFF_REF:%d>%s ", len(qso.WWFF_REF), qso.WWFF_REF))
	}
	if qso.POTA_REF != "" {
		adif.WriteString(fmt.Sprintf("<POTA_REF:%d>%s ", len(qso.POTA_REF), qso.POTA_REF))
	}
	if qso.CNTY != "" {
		adif.WriteString(fmt.Sprintf("<CNTY:%d>%s ", len(qso.CNTY), qso.CNTY))
	}
	if qso.REGION != "" {
		adif.WriteString(fmt.Sprintf("<REGION:%d>%s ", len(qso.REGION), qso.REGION))
	}
	if qso.LAT != "" {
		adif.WriteString(fmt.Sprintf("<LAT:%d>%s ", len(qso.LAT), qso.LAT))
	}
	if qso.LON != "" {
		adif.WriteString(fmt.Sprintf("<LON:%d>%s ", len(qso.LON), qso.LON))
	}
	if qso.ANT_AZ != "" {
		adif.WriteString(fmt.Sprintf("<ANT_AZ:%d>%s ", len(qso.ANT_AZ), qso.ANT_AZ))
	}
	if qso.ANT_EL != "" {
		adif.WriteString(fmt.Sprintf("<ANT_EL:%d>%s ", len(qso.ANT_EL), qso.ANT_EL))
	}
	if qso.ANT_PATH != "" {
		adif.WriteString(fmt.Sprintf("<ANT_PATH:%d>%s ", len(qso.ANT_PATH), qso.ANT_PATH))
	}
	if qso.A_INDEX != "" {
		adif.WriteString(fmt.Sprintf("<A_INDEX:%d>%s ", len(qso.A_INDEX), qso.A_INDEX))
	}
	if qso.K_INDEX != "" {
		adif.WriteString(fmt.Sprintf("<K_INDEX:%d>%s ", len(qso.K_INDEX), qso.K_INDEX))
	}
	if qso.SFI != "" {
		adif.WriteString(fmt.Sprintf("<SFI:%d>%s ", len(qso.SFI), qso.SFI))
	}
	if qso.RX_PWR != "" {
		adif.WriteString(fmt.Sprintf("<RX_PWR:%d>%s ", len(qso.RX_PWR), qso.RX_PWR))
	}

	// End of QSO
	adif.WriteString("<EOR>\n")

	return adif.String()
}