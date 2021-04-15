package formats

// This file contains regular expressions used to parse data from INFO level log lines in the dc_main.log
// eg.: <--[last index]--(DB) pod_uid[1658] time[1591776000] value[9130] (smart_meter_cabinet_facade.cc::315)

// IncomingMessageTypeRegex is used to identify the type of the incoming dc message, which is between the []. eg.: <--[last index]--(DB)
const IncomingMessageTypeRegex = InComingArrow + AnythingBetweenBracketsRegex

// OutGoingMessageTypeRegex is used to identify the type of the outgoing dc message, which is between the []. eg.: --[last index]-->(DB)
const OutGoingMessageTypeRegex = "--" + AnythingBetweenBracketsRegex

// IncomingMessageSourceRegex is used to identify the source of the incoming dc message, which is between the (). eg.: <--[last index]--(DB)
const IncomingMessageSourceRegex = "--" + AnythingBetweenParenthesesRegex

// OutGoingMessageDestRegex is used to identify the destination of the outgoing dc message, which is between the (). eg.: --[last index]-->(DB)
const OutGoingMessageDestRegex = OutGoingArrow + AnythingBetweenParenthesesRegex

// PodUidRegex matches the pod id in a dc message.
const PodUidRegex = "pod_uid" + LongNumberBetweenBracketsRegex

// TimeTicksRegex matches the time field of a dc message, in ticks.
const TimeTicksRegex = "time" + LongNumberBetweenBracketsRegex

// ValueRegex matches the value field of a dc message.
const ValueRegex = "value" + LongNumberBetweenBracketsRegex

// ServiceLevelIdRegex matches the service level id field of a dc message, formatted like: service_level_id[9].
const ServiceLevelIdRegex = "service_level_id" + NumberBetweenBrackets

// TimeRangeFromRegex is used to identify the start fo the time range of the dc message. eg.: from[Wed Jun 10 09:18:33 2020]
const TimeRangeFromRegex = "from\\[" + anyCharsExceptOpeningParentheses + "\\]"

// TimeRangeToRegex is used to identify the end fo the time range of the dc message. eg.: to[Wed Jun 10 10:02:33 2020]
const TimeRangeToRegex = "to\\[" + anyCharsExceptOpeningParentheses + "\\]"

// TimeRangeStartTicksRegex is used to identify the start fo the time range of the dc message in this format: start[1591776600]
const TimeRangeStartTicksRegex = "start" + LongNumberBetweenBracketsRegex

// TimeRangeEndTicksRegex is used to identify the end fo the time range of the dc message in this format: end[1591777200]
const TimeRangeEndTicksRegex = "end" + LongNumberBetweenBracketsRegex

const ConnectOrDisconnectTypeRegex = "type" + NumberBetweenBrackets

// comes next
// client_id[dc18] url[tcp://172.30.31.20:9062] topic[dc.config.dc18] timeout[10000] connected[0]
const ClientIdRegex = "client_id" + AnythingBetweenBracketsRegex

const URLRegex = "url" + AnythingBetweenBracketsRegex

const TopicRegex = "topic" + AnythingBetweenBracketsRegex

const TimeoutRegex = "timeout" + LongNumberBetweenBracketsRegex

const ConnectedRegex = "connected" + NumberBetweenBrackets

// DLMSLogPayload comes next
// DateTimeRegexFieldRegex matches the time field of a dc message, formatted like: time[Wed Jun 10 10:02:35 2020].
const DateTimeFieldRegex = "time\\[" + anyCharsExceptOpeningParentheses + "\\]"

// DLMSRequestTimeRegex matches the dlms request time field of a dc message, formatted like: dlms_request_time[1591783353632].
const DLMSRequestTimeRegex = "dlms_request_time" + LongNumberBetweenBracketsRegex

// DLMSResponseTimeRegex matches the dlms response time field of a dc message, formatted like: dlms_response_time[1591783355857].
const DLMSResponseTimeRegex = "dlms_response_time" + LongNumberBetweenBracketsRegex

// DLMSErrorRegex matches the dlms error field of a dc message, formatted like: dlms_error[OK].
const DLMSErrorRegex = "dlms_error" + AnyLettersBetweenBrackets

// IndexPayload comes next
// PreviousTimeRegex matches the previous time field of a dc message, formatted like: previous_time[1591776000].
const PreviousTimeRegex = "previous_time" + LongNumberBetweenBracketsRegex

// PreviousValueRegex matches the previous value field of a dc message, formatted like: previous_value[15986].
const PreviousValueRegex = "previous_value" + LongNumberBetweenBracketsRegex

// SerailNumberRegex matches the serial number field of a dc message, formatted like: serial_number[98020066928].
const SerailNumberRegex = "serial_number" + LongNumberBetweenBracketsRegex

// CurrentRegex matches the current field of a dc message, formatted like: current[0.1152].
const CurrentRegex = "current" + AnythingBetweenBracketsRegex

// TotalRegex matches the total field of a dc message, formatted like: total[0.1152].
const TotalRegex = "total" + AnythingBetweenBracketsRegex

// SettingsPayload comes next
// DcUidRegex matches the dc uid of a dc message, formatted like: dc_uid[dc18].
const DcUidRegex = "dc_uid" + AnythingBetweenBracketsRegex

// LocalityRegex matches the locality field of a dc message, formatted like: locality[Tanambao Daoud].
const LocalityRegex = "locality" + AnythingBetweenBracketsRegex

// RegionRegex matches the region field of a dc message, formatted like: region[Madagascar].
const RegionRegex = "region" + AnythingBetweenBracketsRegex

// TimezoneRegex matches the timezone field of a dc message, formatted like: timezone[Indian/Antananarivo].
const TimezoneRegex = "timezone" + AnythingBetweenBracketsRegex

// GlobalFtpAddressRegex matches the global ftp address field of a dc message, formatted like: global_ftp_address[sftp://sagemcom@172.30.31.20:firmwares].
const GlobalFtpAddressRegex = "global_ftp_address" + AnythingBetweenBracketsRegex

// TargetFirmwareVersionRegex matches the target firmware version field of a dc message, formatted like: target_firmware_version[].
const TargetFirmwareVersionRegex = "target_firmware_version" + AnythingBetweenBracketsRegex

// IndexCollectionRegex matches the index collection field of a dc message, formatted like: index_collection[600].
const IndexCollectionRegex = "index_collection" + LongNumberBetweenBracketsRegex

// DataPublishRegex matches the data publish field of a dc message, formatted like: data_publish[2400].
const DataPublishRegex = "data_publish" + LongNumberBetweenBracketsRegex

// LastServerCommunicationTimeRegex matches the last server communication time field of a dc message, formatted like: last_server_communication_time[1591775824].
const LastServerCommunicationTimeRegex = "last_server_communication_time" + LongNumberBetweenBracketsRegex

// LastDcStartTimeRegex matches the last dc start time field of a dc message, formatted like: last_dc_start_time[1591737413].
const LastDcStartTimeRegex = "last_dc_start_time" + LongNumberBetweenBracketsRegex

// DcDistroTargetFirmwareVersionRegex matches the dc distro target firmware version field of a dc message, formatted like: dc_distro_target_firmware_version[].
const DcDistroTargetFirmwareVersionRegex = "dc_distro_target_firmware_version" + AnythingBetweenBracketsRegex

// FrequencyBandChangedRegex matches the frequency band changed field of a dc message, formatted like: frequency_band_changed[0].
const FrequencyBandChangedRegex = "frequency_band_changed" + NumberBetweenBrackets

// FrequencyBandRollbackDoneRegex matches the frequency band rollback done field of a dc message, formatted like: frequency_band_rollback_done[0].
const FrequencyBandRollbackDoneRegex = "frequency_band_rollback_done" + NumberBetweenBrackets

// ServiceLevelPayload coming up
// MeterModeRegex matches the meter mode field of a dc message, formatted like: meter_mode[2].
const MeterModeRegex = "meter_mode" + NumberBetweenBrackets

// StartHourDailyCycleRegex matches the start hour daily cycle field of a dc message, formatted like: start_hour_daily_cycle[20h].
const StartHourDailyCycleRegex = "start_hour_daily_cycle" + AnythingBetweenBracketsRegex

// LoadSheddingDailyEnergyBudgetRegex matches the load shedding daily energy budget field of a dc message, formatted like: load_shedding_daily_energy_budget[0].
const LoadSheddingDailyEnergyBudgetRegex = "load_shedding_daily_energy_budget" + LongNumberBetweenBracketsRegex

// LocalSheddingDailyEnergyBudgetRegex matches the local shedding daily energy budget field of a dc message, formatted like: local_shedding_daily_energy_budget[0].
const LocalSheddingDailyEnergyBudgetRegex = "local_shedding_daily_energy_budget" + LongNumberBetweenBracketsRegex

// MaxActivePowerRegex matches the max active power field of a dc message, formatted like: max_active_power[0].
const MaxActivePowerRegex = "max_active_power" + LongNumberBetweenBracketsRegex

// InServiceRegex matches the in service field of a dc message, formatted like: in_service[1].
const InServiceRegex = "in_service" + NumberBetweenBrackets

// NameRegex matches the name field of a dc message, formatted like: name[1. Suspension].
const NameRegex = "name" + AnythingBetweenBracketsRegex

// HourlyEnergyLimitsRegex matches the hourly energy limits field of a dc message, formatted like: hourly_energy_limits[[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]].
const HourlyEnergyLimitsRegex = "hourly_energy_limits" + AnythingBetweenDoubleBracketsRegex

// LocalHourlyEnergyLimitsRegex matches the local hourly energy limits field of a dc message, formatted like: local_hourly_energy_limits[[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]].
const LocalHourlyEnergyLimitsRegex = "local_hourly_energy_limits" + AnythingBetweenDoubleBracketsRegex

// SmcAddressParams, the rest is in the smcjoinformats file
// LastJoiningDateRegex matches the last joining date field of a dc message, formatted like: last_joining_date[Wed Jun 10 00:57:05 2020].
const LastJoiningDateRegex = "last_joining_date\\[" + anyCharsExceptOpeningParentheses + "\\]"

// SmcConfig coming up
// CustomerSerialNumberRegex matches the customer serial number field of a dc message, formatted like: customer_serial_number[SAG0980200000834].
const CustomerSerialNumberRegex = "customer_serial_number" + AnythingBetweenBracketsRegex

// SmcStatusRegex matches the smc status field of a dc message, formatted like: smc_status[]
const SmcStatusRegex = "smc_status" + AnythingBetweenBracketsRegex

// CurrentApp1FwRegex matches the name field of a dc message, formatted like: current_app1_fw[01.07].
const CurrentApp1FwRegex = "current_app1_fw" + AnythingBetweenBracketsRegex

// CurrentApp2FwRegex matches the current_app2_fw field of a dc message, formatted like: current_app2_fw[01.07].
const CurrentApp2FwRegex = "current_app2_fw" + AnythingBetweenBracketsRegex

// CurrentPlcFwRegex matches the current_plc_fw field of a dc message, formatted like: current_plc_fw[14.3.13.0].
const CurrentPlcFwRegex = "current_plc_fw" + AnythingBetweenBracketsRegex

// LastSuccessfulDlmsResponseDateRegex matches the last successful dlms response date field of a dc message, formatted like: last_successful_dlms_response_date [Wed Jun 10 07:58:11 2020].
const LastSuccessfulDlmsResponseDateRegex = "last_successful_dlms_response_date \\[" + anyCharsExceptOpeningParentheses + "\\]"

// NextHopRegex matches the next hop field of a dc message, formatted like: next_hop[0].
const NextHopRegex = "next_hop" + LongNumberBetweenBracketsRegex

// PodConfigPayload rest props coming up
// PhaseRegex matches the phase field of a dc message, formatted like: phase[2].
const PhaseRegex = "phase" + NumberBetweenBrackets

// PositionInSmcRegex matches the position in smc field of a dc message, formatted like: position_in_smc[1].
const PositionInSmcRegex = "position_in_smc" + LongNumberBetweenBracketsRegex

// SerialNumberRegex matches the serial number field of a dc message, formatted like: serial_number[98020067330].
const SerialNumberRegex = "serial_number" + LongNumberBetweenBracketsRegex

// SoftwareFirmwareVersionRegex matches the software firmware version field of a dc message, formatted like: software_firmware_version[IMETER190801].
const SoftwareFirmwareVersionRegex = "software_firmware_version" + AnythingBetweenBracketsRegex

// AnythingBetweenParenthesesRegex matches anything between parentheses.
const AnythingBetweenParenthesesRegex = "\\(([^\\)])*\\)"

// AnythingBetweenBracketsRegex matches anything between brackets.
const AnythingBetweenBracketsRegex = "\\[([^\\]])*\\]"

// AnythingBetweenDoubleBracketsRegex matches anything between two brackets.
const AnythingBetweenDoubleBracketsRegex = "\\[\\[([^\\]])*\\]\\]"
