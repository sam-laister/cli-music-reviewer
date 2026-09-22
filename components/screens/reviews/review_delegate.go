package reviews

import (
	"cli-music-reviewer/models/entities"
	"cli-music-reviewer/services"
	"cli-music-reviewer/styles"
	"fmt"
)

type reviewDelegate struct{}

func (reviewDelegate) Render(item *entities.EntryRow, selected bool) string {
	label := item.Title
	value := fmt.Sprintf("[ %s ]", services.DateToString(item.UpdatedAt))

	if selected {
		label = styles.SelectedItemStyle.Render(label)
		value = styles.SelectedItemStyle.Render(value)
	}

	return fmt.Sprintf("%s %s", label, value)
}
