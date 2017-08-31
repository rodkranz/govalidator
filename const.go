package govalidator

import "regexp"

const MAX_URL_RUNE_COUNT = 2083
const MIN_URL_RUNE_COUNT = 3

const (
	OMIT_EMPTY     string = "OmitEmpty"
	REQUIRED       string = "Required"
	ALPHA_DASH     string = "AlphaDash"
	ALPHA_DASH_DOT string = "AlphaDashDot"
	NUMERIC        string = "Numeric"
	SIZE           string = "Size"
	MIN_SIZE       string = "MinSize"
	MAX_SIZE       string = "MaxSize"
	URL            string = "Url"
	EMAIL          string = "Email"
	INCLUDE        string = "Include"
	EXCLUDE        string = "Exclude"
	DEFAULT        string = "Default"

	ERR_PARSE_VALUE_PARSED_INTEGER         string = "Value could not be parsed as integer"
	ERR_PARSE_VALUE_PARSED_UNSINED_INTEGER string = "Value could not be parsed as unsigned integer"
	ERR_PARSE_VALUE_PARSED_BOOLEAN         string = "Value could not be parsed as boolean"
	ERR_PARSE_VALUE_PARSED_INT_32          string = "Value could not be parsed as 32-bit float"
	ERR_PARSE_VALUE_PARSED_INT_64          string = "Value could not be parsed as 64-bit float"
)

var (
	AlphaDashPattern    = regexp.MustCompile("[^\\d\\w-_]")
	AlphaDashDotPattern = regexp.MustCompile("[^\\d\\w-_\\.]")
	NumericPattern      = regexp.MustCompile("[^\\d]")
	EmailPattern        = regexp.MustCompile("[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?")
)

var (
	urlSchemaRx    = `((ftp|tcp|udp|wss?|https?):\/\/)`
	urlUsernameRx  = `(\S+(:\S*)?@)`
	urlIPRx        = `([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))`
	ipRx           = `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
	urlSubdomainRx = `((www\.)|([a-zA-Z0-9]([-\.][-\._a-zA-Z0-9]+)*))`
	urlPortRx      = `(:(\d{1,5}))`
	urlPathRx      = `((\/|\?|#)[^\s]*)`
	URLPattern     = regexp.MustCompile(`^` + urlSchemaRx + `?` + urlUsernameRx + `?` + `((` + urlIPRx + `|(\[` + ipRx + `\])|(([a-zA-Z0-9]([a-zA-Z0-9-_]+)?[a-zA-Z0-9]([-\.][a-zA-Z0-9]+)*)|(` + urlSubdomainRx + `?))?(([a-zA-Z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-zA-Z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-zA-Z\x{00a1}-\x{ffff}]{1,}))?))\.?` + urlPortRx + `?` + urlPathRx + `?$`)
)
