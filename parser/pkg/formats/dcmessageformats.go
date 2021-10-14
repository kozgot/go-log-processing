package formats

// This file contains regular expressions used to parse data from INFO level log lines in the dc_main.log
// eg.: <--[last index]--(DB) pod_uid[1658] time[1591776000] value[9130] (smart_meter_cabinet_facade.cc::315)

// IncomingMessageTypeRegex is used to identify the type of the incoming dc message, which is between the [].
//  eg.: <--[last index]--(DB)
const IncomingMessageTypeRegex = InComingArrow + AnythingBetweenBracketsRegex

// OutGoingMessageTypeRegex is used to identify the type of the outgoing dc message, which is between the [].
// eg.: --[settings]-->(DB).
const OutGoingMessageTypeRegex = "--" + AnythingBetweenBracketsRegex

// IncomingMessageSourceRegex is used to identify the source of the incoming dc message, which is between the ().
// eg.: <--[last index]--(DB).
const IncomingMessageSourceRegex = "--" + AnythingBetweenParenthesesRegex

// OutGoingMessageDestRegex is used to identify the destination of the outgoing dc message, which is between the ().
// eg.: --[settings]-->(DB).
const OutGoingMessageDestRegex = OutGoingArrow + AnythingBetweenParenthesesRegex

// PodUIDRegex matches the pod id in a dc message.
const PodUIDRegex = "pod_uid" + LongNumberBetweenBracketsRegex

// TimeTicksRegex matches the time field of a dc message, in ticks.
const TimeTicksRegex = "time" + LongNumberBetweenBracketsRegex

// ValueRegex matches the value field of a dc message.
const ValueRegex = "value" + LongNumberBetweenBracketsRegex

// ServiceLevelIDRegex matches the service level id field of a dc message,
// formatted like: service_level_id[9].
const ServiceLevelIDRegex = "service_level_id" + NumberBetweenBrackets

// TimeRangeFromRegex is used to identify the start fo the time range of the dc message.
// eg.: from[Wed Jun 10 09:18:33 2020].
const TimeRangeFromRegex = "from\\[" + anyCharsExceptOpeningParentheses + "\\]"

// TimeRangeToRegex is used to identify the end fo the time range of the dc message.
// eg.: to[Wed Jun 10 10:02:33 2020].
const TimeRangeToRegex = "to\\[" + anyCharsExceptOpeningParentheses + "\\]"

// TimeRangeStartTicksRegex is used to identify the start fo the time range of the dc message
// in this format: start[1591776600].
const TimeRangeStartTicksRegex = "start" + LongNumberBetweenBracketsRegex

// TimeRangeEndTicksRegex is used to identify the end fo the time range of the dc message
// in this format: end[1591777200].
const TimeRangeEndTicksRegex = "end" + LongNumberBetweenBracketsRegex

// ConnectOrDisconnectTypeRegex matches the type field of a log entry.
const ConnectOrDisconnectTypeRegex = "type" + NumberBetweenBrackets

// ClientIDRegex matches the client ID field of a log entry.
// eg.: client_id[dc18] url[tcp://172.30.31.20:9062] topic[dc.config.dc18] timeout[10000] connected[0].
const ClientIDRegex = "client_id" + AnythingBetweenBracketsRegex

// URLRegex matches the URL field of a log entry.
const URLRegex = "url" + AnythingBetweenBracketsRegex

// TopicRegex matches the topic field of a log entry.
const TopicRegex = "topic" + AnythingBetweenBracketsRegex

// TimeoutRegex matches the timeout field of a log entry.
const TimeoutRegex = "timeout" + LongNumberBetweenBracketsRegex

// ConnectedRegex matches the connected field of a log entry.
const ConnectedRegex = "connected" + NumberBetweenBrackets

// DateTimeFieldRegex matches the time field of a dc message, formatted like: time[Wed Jun 10 10:02:35 2020].
const DateTimeFieldRegex = "time\\[" + anyCharsExceptOpeningParentheses + "\\]"

// DLMSRequestTimeRegex matches the dlms request time field of a dc message, eg.: dlms_request_time[1591783353632].
const DLMSRequestTimeRegex = "dlms_request_time" + LongNumberBetweenBracketsRegex

// DLMSResponseTimeRegex matches the dlms response time field of a dc message, eg.: dlms_response_time[1591783355857].
const DLMSResponseTimeRegex = "dlms_response_time" + LongNumberBetweenBracketsRegex

// DLMSErrorRegex matches the dlms error field of a dc message, eg.: dlms_error[OK].
const DLMSErrorRegex = "dlms_error" + AnyLettersBetweenBrackets

// PreviousTimeRegex matches the previous time field of a dc message, formatted like: previous_time[1591776000].
const PreviousTimeRegex = "previous_time" + LongNumberBetweenBracketsRegex

// PreviousValueRegex matches the previous value field of a dc message, formatted like: previous_value[15986].
const PreviousValueRegex = "previous_value" + LongNumberBetweenBracketsRegex

// CurrentRegex matches the current field of a dc message, formatted like: current[0.1152].
const CurrentRegex = "current" + AnythingBetweenBracketsRegex

// TotalRegex matches the total field of a dc message, formatted like: total[0.1152].
const TotalRegex = "total" + AnythingBetweenBracketsRegex

// MeterModeRegex matches the meter mode field of a service_level message, formatted like: meter_mode[2].
const MeterModeRegex = "meter_mode" + NumberBetweenBrackets

// StartHourDailyCycleRegex matches the start hour daily cycle field of a service_level message,
// eg.: start_hour_daily_cycle[20h].
const StartHourDailyCycleRegex = "start_hour_daily_cycle" + AnythingBetweenBracketsRegex

// LoadSheddingDailyEnergyBudgetRegex matches the load shedding daily energy budget field of a service_level message.
// eg.: load_shedding_daily_energy_budget[0].
const LoadSheddingDailyEnergyBudgetRegex = "load_shedding_daily_energy_budget" + LongNumberBetweenBracketsRegex

// LocalSheddingDailyEnergyBudgetRegex matches the local shedding daily energy budget field of a service_level message.
// eg.: local_shedding_daily_energy_budget[0].
const LocalSheddingDailyEnergyBudgetRegex = "local_shedding_daily_energy_budget" + LongNumberBetweenBracketsRegex

// MaxActivePowerRegex matches the max active power field of a service_level message,
// formatted like: max_active_power[0].
const MaxActivePowerRegex = "max_active_power" + LongNumberBetweenBracketsRegex

// InServiceRegex matches the in service field of a service_level message, formatted like: in_service[1].
const InServiceRegex = "in_service" + NumberBetweenBrackets

// NameRegex matches the name field of a service_level message, formatted like: name[1. Suspension].
const NameRegex = "name" + AnythingBetweenBracketsRegex

// HourlyEnergyLimitsRegex matches the hourly energy limits field of a service_level message.
// eg.: hourly_energy_limits[[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]].
const HourlyEnergyLimitsRegex = "hourly_energy_limits" + AnythingBetweenDoubleBracketsRegex

// LocalHourlyEnergyLimitsRegex matches the local hourly energy limits field of a service_level message.
// eg.: local_hourly_energy_limits[[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]].
const LocalHourlyEnergyLimitsRegex = "local_hourly_energy_limits" + AnythingBetweenDoubleBracketsRegex

// LastJoiningDateRegex matches the last joining date field of a dc message.
// eg.: last_joining_date[Wed Jun 10 00:57:05 2020].
const LastJoiningDateRegex = "last_joining_date\\[" + anyCharsExceptOpeningParentheses + "\\]"

// CustomerSerialNumberRegex matches the customer serial number field of a dc message.
// eg.: customer_serial_number[SAG0980200000834].
const CustomerSerialNumberRegex = "customer_serial_number" + AnythingBetweenBracketsRegex

// SmcStatusRegex matches the smc status field of a dc message, formatted like: smc_status[].
const SmcStatusRegex = "smc_status" + AnythingBetweenBracketsRegex

// CurrentApp1FwRegex matches the name field of a dc message, formatted like: current_app1_fw[01.07].
const CurrentApp1FwRegex = "current_app1_fw" + AnythingBetweenBracketsRegex

// CurrentApp2FwRegex matches the current_app2_fw field of a dc message, formatted like: current_app2_fw[01.07].
const CurrentApp2FwRegex = "current_app2_fw" + AnythingBetweenBracketsRegex

// CurrentPlcFwRegex matches the current_plc_fw field of a dc message, formatted like: current_plc_fw[14.3.13.0].
const CurrentPlcFwRegex = "current_plc_fw" + AnythingBetweenBracketsRegex

// <--[smc configuration]--(DB)
// LastSuccessfulRespDateRegex matches the last successful dlms response date field of a dc message,
// formatted like: last_successful_dlms_response_date [Wed Jun 10 07:58:11 2020].
const LastSuccessfulRespDateRegex = "last_successful_dlms_response_date \\[" + anyCharsExceptOpeningParentheses + "\\]"

// NextHopRegex matches the next hop field of a dc message, formatted like: next_hop[0].
const NextHopRegex = "next_hop" + LongNumberBetweenBracketsRegex

// PhaseRegex matches the phase field of a dc message, formatted like: phase[2].
const PhaseRegex = "phase" + NumberBetweenBrackets

// PositionInSmcRegex matches the position in smc field of a dc message, formatted like: position_in_smc[1].
const PositionInSmcRegex = "position_in_smc" + LongNumberBetweenBracketsRegex

// SerialNumberRegex matches the serial number field of a dc message, formatted like: serial_number[98020067330].
const SerialNumberRegex = "serial_number" + LongNumberBetweenBracketsRegex

// SoftwareFirmwareVersionRegex matches the software firmware version field of a dc message,
// formatted like: software_firmware_version[IMETER190801].
const SoftwareFirmwareVersionRegex = "software_firmware_version" + AnythingBetweenBracketsRegex

// AnythingBetweenParenthesesRegex matches anything between parentheses.
const AnythingBetweenParenthesesRegex = "\\(([^\\)])*\\)"

// AnythingBetweenBracketsRegex matches anything between brackets.
const AnythingBetweenBracketsRegex = "\\[([^\\]])*\\]"

// AnythingBetweenDoubleBracketsRegex matches anything between two brackets.
const AnythingBetweenDoubleBracketsRegex = "\\[\\[([^\\]])*\\]\\]"

// ConnectToPLCIfaceRegex matches the interface field of a connect to plc log entry.
// eg.: --[connect]-->(PLC) Iface[plc0] - DestAddr[fe80::4021:ff:fe00:a:61616].
const ConnectToPLCIfaceRegex = "Iface" + AnyLettersBetweenBrackets

// ConnectToPLCDestAddressRegex matches the destination address field of a connect to plc log entry.
// eg.: --[connect]-->(PLC) Iface[plc0] - DestAddr[fe80::4021:ff:fe00:a:61616].
const ConnectToPLCDestAddressRegex = "DestAddr" + AnyLettersBetweenBrackets

// IndexProfileCapturePeriodRegex matches the capture period field of an [index high profile generic],
// or [index low profile generic] log entry.
// eg.: capture_period[86400].
const IndexProfileCapturePeriodRegex = "capture_period" + LongNumberBetweenBracketsRegex

// IndexProfileCapturePeriodRegex matches the capture objects field of an [index high profile generic],
// or [index low profile generic] log entry.
// eg.: capture_objects[0].
const IndexProfileCaptureObjectsRegex = "capture_objects" + LongNumberBetweenBracketsRegex
