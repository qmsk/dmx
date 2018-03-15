package api

type ColorID string

type Colors map[ColorID]Color

// Merge in colors from given map, overriding colors in this map
func (colors Colors) Merge(merge Colors) Colors {
	var merged = make(Colors)

	for colorID, color := range colors {
		merged[colorID] = color
	}

	for colorID, color := range merge {
		merged[colorID] = color
	}

	return merged
}

type Color struct {
	Red   Value
	Green Value
	Blue  Value
}

type ColorParams struct {
	ScaleIntensity *Value
	Color
}
