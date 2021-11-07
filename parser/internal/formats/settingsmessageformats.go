package formats

// This file contains regular expressions needed to
// parse log entries that start like this (after the timestamp and log level):
// --[settings]-->(DB) dc_uid[dc18] locality[Tanambao Daoud] region[Madagascar] timezone[Indian/Antananarivo]

// DcUIDRegex matches the dc uid of a dc message, formatted like: dc_uid[dc18].
const DcUIDRegex = "dc_uid" + AnythingBetweenBracketsRegex

// LocalityRegex matches the locality field of a dc message, formatted like: locality[Tanambao Daoud].
const LocalityRegex = "locality" + AnythingBetweenBracketsRegex

// RegionRegex matches the region field of a dc message, formatted like: region[Madagascar].
const RegionRegex = "region" + AnythingBetweenBracketsRegex

// TimezoneRegex matches the timezone field of a dc message, formatted like: timezone[Indian/Antananarivo].
const TimezoneRegex = "timezone" + AnythingBetweenBracketsRegex

// GlobalFtpAddressRegex matches the global ftp address field of a dc message log entry.
// eg.: global_ftp_address[sftp://sagemcom@172.30.31.20:firmwares].
const GlobalFtpAddressRegex = "global_ftp_address" + AnythingBetweenBracketsRegex

// TargetFirmwareVersionRegex matches the target firmware version field of a dc message log entry.
// eg.: target_firmware_version[].
const TargetFirmwareVersionRegex = "target_firmware_version" + AnythingBetweenBracketsRegex

// IndexCollectionRegex matches the index collection field of a dc message, formatted like: index_collection[600].
const IndexCollectionRegex = "index_collection" + LongNumberBetweenBracketsRegex

// DataPublishRegex matches the data publish field of a dc message, formatted like: data_publish[2400].
const DataPublishRegex = "data_publish" + LongNumberBetweenBracketsRegex

// LastServerCommunicationTimeRegex matches the last server communication time field of a dc message.
// eg.: last_server_communication_time[1591775824].
const LastServerCommunicationTimeRegex = "last_server_communication_time" + LongNumberBetweenBracketsRegex

// LastDcStartTimeRegex matches the last dc start time field of a dc message, eg.: last_dc_start_time[1591737413].
const LastDcStartTimeRegex = "last_dc_start_time" + LongNumberBetweenBracketsRegex

// DcDistroTargetFirmwareVersionRegex matches the dc distro target firmware version field of a dc message.
// eg.: dc_distro_target_firmware_version[].
const DcDistroTargetFirmwareVersionRegex = "dc_distro_target_firmware_version" + AnythingBetweenBracketsRegex

// FrequencyBandChangedRegex matches the frequency band changed field of a dc message, eg.: frequency_band_changed[0].
const FrequencyBandChangedRegex = "frequency_band_changed" + SingleNumberBetweenBrackets

// FrequencyBandRollbackDoneRegex matches the frequency band rollback done field of a dc message.
// eg.: frequency_band_rollback_done[0].
const FrequencyBandRollbackDoneRegex = "frequency_band_rollback_done" + SingleNumberBetweenBrackets
