package migrations

import (
	"github.com/pocketbase/pocketbase/core"
)

func init() {
	Register(func(app core.App) error {
		// Find the users collection
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			// If users collection doesn't exist, skip this migration
			return nil
		}

		// Check if phone field already exists
		phoneField := users.Fields.GetByName("phone")
		if phoneField != nil {
			// Phone field already exists, skip
			return nil
		}

		// Add phone field to users collection
		users.Fields.Add(&core.TextField{
			Name: "phone",
			Max:  20,
		})

		// Save the updated collection
		return app.Save(users)
	}, func(app core.App) error {
		// Find the users collection
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return nil
		}

		// Remove phone field if it exists
		phoneField := users.Fields.GetByName("phone")
		if phoneField != nil {
			users.Fields.RemoveByName("phone")
			return app.Save(users)
		}

		return nil
	})
}
