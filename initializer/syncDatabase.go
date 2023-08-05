package initializer

import "github.com/aswinjithkukku/ecom-gingorm/models"

func SyncDatabase() {

	DB.AutoMigrate(&models.Products{}, &models.Users{})

}
