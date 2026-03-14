package plans

import (
	"fmt"

	"github.com/thebeyond-net/control-plane/internal/core/domain"
	"github.com/thebeyond-net/control-plane/internal/i18n"
)

func formatPrice(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}

func (uc *UseCase) formatPlanDetails(plan domain.Plan, languageCode, currencyCode string, bandwidth, days int) string {
	periodDiscount := int(uc.getPeriodDiscountFraction(days) * 100)
	price, err := plan.GetPrice(currencyCode, bandwidth, days, periodDiscount)
	if err != nil {
		return i18n.Get(languageCode, "PriceUnavailable", nil, nil)
	}

	currency, _ := uc.currencies.Get(currencyCode)

	period := i18n.Get(languageCode, "DaysCount", map[string]any{
		"Count": days,
	}, days)

	return i18n.Get(languageCode, "PlanDescription", map[string]any{
		"Period":         period,
		"Devices":        plan.Devices,
		"Bandwidth":      formatBandwidth(bandwidth),
		"Price":          formatPrice(price),
		"CurrencySymbol": currency.Symbol,
	}, nil)
}

func formatBandwidth(bw int) string {
	if bw >= 1000 {
		return fmt.Sprintf("%g Gbps", float64(bw)/1000)
	}
	return fmt.Sprintf("%d Mbps", bw)
}

func (uc *UseCase) getPeriodDiscountFraction(days int) float64 {
	for _, p := range uc.periods {
		if p.Days == days {
			return float64(p.Discount) / 100.0
		}
	}
	return 0
}
