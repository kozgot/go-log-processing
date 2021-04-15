package formats

// RoutingTableRegex represents the regular expression that matches a log line that contains data regarding the Rounting Table
const RoutingTableRegex = "(Routing Table: )"

// RoutingAddressRegex matches the Addr[...] field
const RoutingAddressRegex = "Addr" + AnyLettersBetweenBrackets

// NextHopAddressRegex matches the NextHopAddr[...] field
const NextHopAddressRegex = "NextHopAddr" + AnyLettersBetweenBrackets

// RouteCostRegex matches the RouteCost[...] field
const RouteCostRegex = "RouteCost" + AnyLettersBetweenBrackets // todo number

// HopCountRegex matches the HopCount[...] field
const HopCountRegex = "HopCount" + AnyLettersBetweenBrackets

// WeakLinkRegex matches the WeakLink[...] field
const WeakLinkRegex = "WeakLink" + AnyLettersBetweenBrackets // todo number

// ValidTimeRegex matches the ValidTime[...] field
const ValidTimeRegex = "ValidTime" + AnyLettersBetweenBrackets // todo number
