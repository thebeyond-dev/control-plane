package domain

import (
	"fmt"
)

type Plan struct {
	ID              int
	Devices         int
	Discount        int
	BandwidthPrices map[int]map[string]float64
}

func (plan Plan) GetPrice(currency string, bandwidth, days, periodDiscount int) (float64, error) {
	currencyPrices, ok := plan.BandwidthPrices[bandwidth]
	if !ok {
		return 0, fmt.Errorf("bandwidth %d not found in plan", bandwidth)
	}

	dailyPrice, ok := currencyPrices[currency]
	if !ok {
		return 0, fmt.Errorf("price not found for currency %s", currency)
	}

	discount := float64(plan.Discount+periodDiscount) / 100
	totalPrice := dailyPrice * float64(days) * (1 - discount)
	return totalPrice, nil
}

func GetNeighborIDs(currentID, totalPlans int) (prevID, nextID int, hasPrev, hasNext bool) {
	return currentID - 1,
		currentID + 1,
		currentID > 0,
		currentID < totalPlans-1
}

func GetNeighborBandwidths(current int, allBandwidths []int) (prev, next int, hasPrev, hasNext bool) {
	for i, bw := range allBandwidths {
		if bw == current {
			if i > 0 {
				prev = allBandwidths[i-1]
				hasPrev = true
			}
			if i < len(allBandwidths)-1 {
				next = allBandwidths[i+1]
				hasNext = true
			}
			break
		}
	}
	return
}
