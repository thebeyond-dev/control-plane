package config

import (
	"github.com/thebeyond-net/control-plane/internal/core/domain"
)

func ToDomainPeriods(dto []PeriodDTO) []domain.Period {
	periods := make([]domain.Period, len(dto))
	for i, period := range dto {
		periods[i] = domain.Period{
			Days:     period.Days,
			Discount: period.Discount,
			NameKey:  period.NameKey,
		}
	}
	return periods
}

func ToDomainPlans(dto []PlanDTO) []domain.Plan {
	plans := make([]domain.Plan, len(dto))
	for i, plan := range dto {
		plans[i] = domain.Plan{
			ID:              i,
			Devices:         plan.Devices,
			Discount:        plan.Discount,
			BandwidthPrices: calculateDailyBandwidthPrices(plan.Prices),
		}
	}
	return plans
}

func calculateDailyBandwidthPrices(prices map[int]map[string]float64) map[int]map[string]float64 {
	dailyBandwidthPrices := make(map[int]map[string]float64, len(prices))

	for bandwidth, currencyPrices := range prices {
		dailyPrices := make(map[string]float64, len(currencyPrices))
		for currencyCode, monthlyPrice := range currencyPrices {
			dailyPrices[currencyCode] = monthlyPrice / 30
		}
		dailyBandwidthPrices[bandwidth] = dailyPrices
	}

	return dailyBandwidthPrices
}

func ToDomain(items map[string]ItemDTO) domain.Items {
	itemsMap := domain.NewItems()
	for code, item := range items {
		itemsMap.Add(domain.Item{
			Code:   code,
			Symbol: item.Symbol,
			Name:   item.Name,
			Icon:   item.Icon,
		})
	}
	return itemsMap
}
