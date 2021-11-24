package contentparser

import (
	"github.com/kozgot/go-log-processing/parser/internal/common"
	"github.com/kozgot/go-log-processing/parser/internal/formats"
	"github.com/kozgot/go-log-processing/parser/pkg/models"
)

type RoutingMessageParser struct {
	line models.EntryWithLevelAndTimestamp
}

func (r *RoutingMessageParser) Parse() *models.RoutingTableParams {
	isRoutingTableLine := common.ParseFieldAsString(r.line.Rest, formats.RoutingTableRegex) != ""
	if !isRoutingTableLine {
		return nil
	}

	routingtableLine := models.RoutingTableParams{}
	routingtableLine.Address = common.ParseFieldInBracketsAsString(r.line.Rest, formats.RoutingAddressRegex)
	routingtableLine.NextHopAddress = common.ParseFieldInBracketsAsString(r.line.Rest, formats.NextHopAddressRegex)
	routingtableLine.RouteCost = common.TryParseIntFromString(
		common.ParseFieldInBracketsAsString(r.line.Rest, formats.RouteCostRegex))
	routingtableLine.ValidTimeMins = common.TryParseIntFromString(
		common.ParseFieldInBracketsAsString(r.line.Rest, formats.ValidTimeRegex))
	routingtableLine.WeakLink = common.TryParseIntFromString(
		common.ParseFieldInBracketsAsString(r.line.Rest, formats.WeakLinkRegex))
	routingtableLine.HopCount = common.TryParseIntFromString(
		common.ParseFieldInBracketsAsString(r.line.Rest, formats.HopCountRegex))

	return &routingtableLine
}
