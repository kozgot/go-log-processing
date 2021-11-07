package models

import "time"

// DcMessagePayload contains the parsed payload of info level messages that have been sent or received by the dc.
type DcMessagePayload struct {
	SmcUID         string
	PodUID         string
	ServiceLevelID int
	Value          int
	Time           time.Time

	TimeRange *TimeRange

	ConnectOrDisconnectPayload       *ConnectOrDisconnectPayload
	DLMSLogPayload                   *DLMSLogPayload
	IndexPayload                     *IndexPayload
	GenericIndexProfilePayload       *GenericIndexProfilePayload
	MessagePayload                   *MessagePayload
	SettingsPayload                  *SettingsPayload
	ServiceLevelPayload              *ServiceLevelPayload
	SmcAddressPayload                *SmcAddressParams
	SmcConfigPayload                 *SmcConfigPayload
	PodConfigPayload                 *PodConfigPayload
	ConnectToPLCPayload              *ConnectToPLCPayload
	StatisticsEntryPayload           *StatisticsEntryPayload
	ReadIndexLowProfilesEntryPayload *ReadIndexLowProfilesEntryPayload
	ReadIndexProfilesEntryPayload    *ReadIndexProfilesEntryPayload
}

// Equals checks equality.
func (d *DcMessagePayload) Equals(other DcMessagePayload) bool {
	if d.ServiceLevelID != other.ServiceLevelID ||
		d.PodUID != other.PodUID ||
		d.SmcUID != other.SmcUID ||
		d.Time != other.Time ||
		d.Value != other.Value {
		return false
	}

	// TimeRange
	if d.TimeRange != nil && other.TimeRange == nil ||
		d.TimeRange == nil && other.TimeRange != nil {
		return false
	}
	if d.TimeRange != nil && !d.TimeRange.Equals(*other.TimeRange) {
		return false
	}

	// ConnectOrDisconnectPayload
	if d.ConnectOrDisconnectPayload != nil && other.ConnectOrDisconnectPayload == nil ||
		d.ConnectOrDisconnectPayload == nil && other.ConnectOrDisconnectPayload != nil {
		return false
	}
	if d.ConnectOrDisconnectPayload != nil && !d.ConnectOrDisconnectPayload.Equals(*other.ConnectOrDisconnectPayload) {
		return false
	}

	// DLMSLogPayload
	if d.DLMSLogPayload != nil && other.DLMSLogPayload == nil ||
		d.DLMSLogPayload == nil && other.DLMSLogPayload != nil {
		return false
	}
	if d.DLMSLogPayload != nil && !d.DLMSLogPayload.Equals(*other.DLMSLogPayload) {
		return false
	}

	// IndexPayload
	if d.IndexPayload != nil && other.IndexPayload == nil ||
		d.IndexPayload == nil && other.IndexPayload != nil {
		return false
	}
	if d.IndexPayload != nil && !d.IndexPayload.Equals(*other.IndexPayload) {
		return false
	}

	// GenericIndexProfilePayload
	if d.GenericIndexProfilePayload != nil && other.GenericIndexProfilePayload == nil ||
		d.GenericIndexProfilePayload == nil && other.GenericIndexProfilePayload != nil {
		return false
	}
	if d.GenericIndexProfilePayload != nil && !d.GenericIndexProfilePayload.Equals(*other.GenericIndexProfilePayload) {
		return false
	}

	// MessagePayload
	if d.MessagePayload != nil && other.MessagePayload == nil ||
		d.MessagePayload == nil && other.MessagePayload != nil {
		return false
	}
	if d.MessagePayload != nil && !d.MessagePayload.Equals(*other.MessagePayload) {
		return false
	}

	// SettingsPayload
	if d.SettingsPayload != nil && other.SettingsPayload == nil ||
		d.SettingsPayload == nil && other.SettingsPayload != nil {
		return false
	}
	if d.SettingsPayload != nil && !d.SettingsPayload.Equals(*other.SettingsPayload) {
		return false
	}

	// ServiceLevelPayload
	if d.ServiceLevelPayload != nil && other.ServiceLevelPayload == nil ||
		d.ServiceLevelPayload == nil && other.ServiceLevelPayload != nil {
		return false
	}
	if d.ServiceLevelPayload != nil && !d.ServiceLevelPayload.Equals(*other.ServiceLevelPayload) {
		return false
	}

	// SmcAddressPayload
	if d.SmcAddressPayload != nil && other.SmcAddressPayload == nil ||
		d.SmcAddressPayload == nil && other.SmcAddressPayload != nil {
		return false
	}
	if d.SmcAddressPayload != nil && !d.SmcAddressPayload.Equals(*other.SmcAddressPayload) {
		return false
	}

	// SmcConfigPayload
	if d.SmcConfigPayload != nil && other.SmcConfigPayload == nil ||
		d.SmcConfigPayload == nil && other.SmcConfigPayload != nil {
		return false
	}
	if d.SmcConfigPayload != nil && !d.SmcConfigPayload.Equals(*other.SmcConfigPayload) {
		return false
	}

	// PodConfigPayload
	if d.PodConfigPayload != nil && other.PodConfigPayload == nil ||
		d.PodConfigPayload == nil && other.PodConfigPayload != nil {
		return false
	}
	if d.PodConfigPayload != nil && !d.PodConfigPayload.Equals(*other.PodConfigPayload) {
		return false
	}

	// ConnectToPLCPayload
	if d.ConnectToPLCPayload != nil && other.ConnectToPLCPayload == nil ||
		d.ConnectToPLCPayload == nil && other.ConnectToPLCPayload != nil {
		return false
	}
	if d.ConnectToPLCPayload != nil && !d.ConnectToPLCPayload.Equals(*other.ConnectToPLCPayload) {
		return false
	}

	// StatisticsEntryPayload
	if d.StatisticsEntryPayload != nil && other.StatisticsEntryPayload == nil ||
		d.StatisticsEntryPayload == nil && other.StatisticsEntryPayload != nil {
		return false
	}
	if d.StatisticsEntryPayload != nil && !d.StatisticsEntryPayload.Equals(*other.StatisticsEntryPayload) {
		return false
	}

	// ReadIndexLowProfilesEntryPayload
	if d.ReadIndexLowProfilesEntryPayload != nil && other.ReadIndexLowProfilesEntryPayload == nil ||
		d.ReadIndexLowProfilesEntryPayload == nil && other.ReadIndexLowProfilesEntryPayload != nil {
		return false
	}
	if d.ReadIndexLowProfilesEntryPayload != nil &&
		!d.ReadIndexLowProfilesEntryPayload.Equals(*other.ReadIndexLowProfilesEntryPayload) {
		return false
	}

	// ReadIndexProfilesEntryPayload
	if d.ReadIndexProfilesEntryPayload != nil && other.ReadIndexProfilesEntryPayload == nil ||
		d.ReadIndexProfilesEntryPayload == nil && other.ReadIndexProfilesEntryPayload != nil {
		return false
	}
	if d.ReadIndexProfilesEntryPayload != nil && !d.ReadIndexProfilesEntryPayload.Equals(*other.ReadIndexProfilesEntryPayload) {
		return false
	}

	return true
}

// ReadIndexLowProfilesEntryPayload contains the parsed --[read index low profiles]-->(SMC) entries.
type ReadIndexLowProfilesEntryPayload struct {
	SmcUID string
	From   time.Time
	To     time.Time
}

// Equals checks equality.
func (r *ReadIndexLowProfilesEntryPayload) Equals(other ReadIndexLowProfilesEntryPayload) bool {
	return r.SmcUID == other.SmcUID &&
		r.From == other.From &&
		r.To == other.To
}

// ReadIndexProfilesEntryPayload contains the parsed <--[read index profiles]--(SMC) entries.
type ReadIndexProfilesEntryPayload struct {
	SmcUID string
	Count  int
}

// Equals checks equality.
func (r *ReadIndexProfilesEntryPayload) Equals(other ReadIndexProfilesEntryPayload) bool {
	return r.SmcUID == other.SmcUID &&
		r.Count == other.Count
}

// StatisticsEntryPayload contains the parsed statistics log entry sent to the SVI.
type StatisticsEntryPayload struct {
	Type     string
	Value    float64
	Time     time.Time
	SourceID string
}

// Equals checks equality.
func (s *StatisticsEntryPayload) Equals(other StatisticsEntryPayload) bool {
	return s.Type == other.Type &&
		s.Value == other.Value &&
		s.Time == other.Time &&
		s.SourceID == other.SourceID
}

// GenericIndexProfilePayload contains the parsed index high/low profile generic payload.
type GenericIndexProfilePayload struct {
	CapturePeriod  int
	CaptureObjects int
}

// Equals checks equality.
func (g *GenericIndexProfilePayload) Equals(other GenericIndexProfilePayload) bool {
	return g.CapturePeriod == other.CapturePeriod &&
		g.CaptureObjects == other.CaptureObjects
}

// ConnectToPLCPayload contains the parsed connect to PLC payload.
type ConnectToPLCPayload struct {
	Interface          string
	DestinationAddress string
}

// Equals checks equality.
func (c *ConnectToPLCPayload) Equals(other ConnectToPLCPayload) bool {
	return c.Interface == other.Interface &&
		c.DestinationAddress == other.DestinationAddress
}

// SettingsPayload contains the parsed settings payload
// of info level messages that have been sent or received by the dc.
type SettingsPayload struct {
	DcUID                         string
	Locality                      string
	Region                        string
	Timezone                      string
	GlobalFtpAddress              string
	TargetFirmwareVersion         string
	IndexCollection               int
	DataPublish                   int
	LastServerCommunicationTime   time.Time
	DcDistroTargetFirmwareVersion string
	LastDcStartTime               time.Time // it might be ticks or something (eg. 1591780709)
	FrequencyBandChanged          bool
	FrequencyBandRollBackDone     bool
}

// Equals checks equality.
func (s *SettingsPayload) Equals(other SettingsPayload) bool {
	return s.DcUID == other.DcUID &&
		s.Locality == other.Locality &&
		s.Region == other.Region &&
		s.Timezone == other.Timezone &&
		s.GlobalFtpAddress == other.GlobalFtpAddress &&
		s.TargetFirmwareVersion == other.TargetFirmwareVersion &&
		s.IndexCollection == other.IndexCollection &&
		s.DataPublish == other.DataPublish &&
		s.LastServerCommunicationTime == other.LastServerCommunicationTime &&
		s.DcDistroTargetFirmwareVersion == other.DcDistroTargetFirmwareVersion &&
		s.LastDcStartTime == other.LastDcStartTime &&
		s.FrequencyBandChanged == other.FrequencyBandChanged &&
		s.FrequencyBandRollBackDone == other.FrequencyBandRollBackDone
}

// ServiceLevelPayload contains the parsed service level related data
// in info level log entries that describe messages sent or received by the dc.
type ServiceLevelPayload struct {
	MeterMode                      int
	StartHourDailyCycle            string // eg. 20h
	LoadSheddingDailyEnergyBudget  int
	LocalSheddingDailyEnergyBudget int
	MaxActivePower                 int
	InService                      bool
	Name                           string
	HourlyEnergyLimits             [24]HourlyEnergyLimit
	LocalHourlyEnergyLimits        [24]HourlyEnergyLimit
}

// Equals checks equality.
func (s *ServiceLevelPayload) Equals(other ServiceLevelPayload) bool {
	return s.MeterMode == other.MeterMode &&
		s.StartHourDailyCycle == other.StartHourDailyCycle &&
		s.LoadSheddingDailyEnergyBudget == other.LoadSheddingDailyEnergyBudget &&
		s.LocalSheddingDailyEnergyBudget == other.LocalSheddingDailyEnergyBudget &&
		s.MaxActivePower == other.MaxActivePower &&
		s.InService == other.InService &&
		s.Name == other.Name
	// todo: limits
}

// HourlyEnergyLimit contains the value and the hour number of an hourly energy limit.
type HourlyEnergyLimit struct {
	HourNumber int
	Limit      int
}

// Equals checks equality.
func (h *HourlyEnergyLimit) Equals(other HourlyEnergyLimit) bool {
	return h.HourNumber == other.HourNumber &&
		h.Limit == other.Limit
}

// SmcConfigPayload contains data related to the configuration of an SMC.
type SmcConfigPayload struct {
	CustomerSerialNumber           string
	PhysicalAddress                string
	SmcStatus                      string
	CurrentApp1Fw                  string
	CurrentApp2Fw                  string
	CurrentPlcFw                   string
	LastSuccessfulDlmsResponseDate time.Time
	NextHop                        int
}

// Equals checks equality.
func (s *SmcConfigPayload) Equals(other SmcConfigPayload) bool {
	return s.CustomerSerialNumber == other.CustomerSerialNumber &&
		s.PhysicalAddress == other.PhysicalAddress &&
		s.SmcStatus == other.SmcStatus &&
		s.CurrentApp1Fw == other.CurrentApp1Fw &&
		s.CurrentApp2Fw == other.CurrentApp2Fw &&
		s.CurrentPlcFw == other.CurrentPlcFw &&
		s.LastSuccessfulDlmsResponseDate == other.LastSuccessfulDlmsResponseDate &&
		s.NextHop == other.NextHop
}

// MessagePayload contains the parameters of a log entry related to a DC message.
type MessagePayload struct {
	Current float64
	Total   float64
	URL     string
	Topic   string
}

// Equals checks equality.
func (m *MessagePayload) Equals(other MessagePayload) bool {
	return m.Current == other.Current &&
		m.Total == other.Total &&
		m.URL == other.URL &&
		m.Topic == other.Topic
}

// PodConfigPayload contains the parameters of a log entry related to a pod configuration event.
type PodConfigPayload struct {
	SerialNumber            int
	Phase                   int
	PositionInSmc           int
	SoftwareFirmwareVersion string
}

// Equals checks equality.
func (p *PodConfigPayload) Equals(other PodConfigPayload) bool {
	return p.SerialNumber == other.SerialNumber &&
		p.Phase == other.Phase &&
		p.PositionInSmc == other.PositionInSmc &&
		p.SoftwareFirmwareVersion == other.SoftwareFirmwareVersion
}

// TimeRange represents a time range.
type TimeRange struct {
	From time.Time
	To   time.Time
}

// Equals checks equality.
func (t *TimeRange) Equals(other TimeRange) bool {
	return t.From == other.From && t.To == other.To
}

// DLMSLogPayload contains data of a log entry related to DLMS Log contents.
type DLMSLogPayload struct {
	DLMSRequestTime  time.Time
	DLMSResponseTime time.Time
	DLMSError        string
}

// Equals checks equality.
func (d *DLMSLogPayload) Equals(other DLMSLogPayload) bool {
	return d.DLMSRequestTime == other.DLMSRequestTime &&
		d.DLMSResponseTime == other.DLMSResponseTime &&
		d.DLMSError == other.DLMSError
}

// IndexPayload contains data of an index log entry.
type IndexPayload struct {
	PreviousTime  time.Time // it might be ticks or something (eg. 1591776000)
	PreviousValue int
	SerialNumber  int
}

// Equals checks equality.
func (i *IndexPayload) Equals(other IndexPayload) bool {
	return i.PreviousTime == other.PreviousTime &&
		i.PreviousValue == other.PreviousValue &&
		i.SerialNumber == other.SerialNumber
}

// ConnectOrDisconnectPayload contains information related to a connect or disconnect event.
type ConnectOrDisconnectPayload struct {
	Type      int
	ClientID  string
	URL       string
	Topic     string
	Timeout   int
	Connected bool
}

// Equals checks equality.
func (c *ConnectOrDisconnectPayload) Equals(other ConnectOrDisconnectPayload) bool {
	return c.Type == other.Type &&
		c.ClientID == other.ClientID &&
		c.URL == other.URL &&
		c.Topic == other.Topic &&
		c.Timeout == other.Timeout &&
		c.Connected == other.Connected
}

// Calendar contains data related to the calendar of a tariff settings log entry.
type Calendar struct {
	CalendarName CalendarName
}

// Equals checks equality.
func (c *Calendar) Equals(other Calendar) bool {
	return c.CalendarName.Equals(other.CalendarName)
}

// CalendarName contains data related to the calendar of a tariff settings log entry.
type CalendarName struct {
	IsActive      bool
	SeasonProfile string // for now, the exact type is unknown
	WeekProfile   string // for now, the exact type is unknown
	DayProfile    string // for now, the exact type is unknown
}

// Equals checks equality.
func (c *CalendarName) Equals(other CalendarName) bool {
	return c.IsActive == other.IsActive &&
		c.SeasonProfile == other.SeasonProfile &&
		c.DayProfile == other.DayProfile &&
		c.WeekProfile == other.WeekProfile
}
