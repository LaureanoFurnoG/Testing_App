package main
import(
	"Go-API-T/initializers"
	"Go-API-T/models"
)
//initializers db and .env
func init(){
	initializers.InitEnv()
	initializers.ConnectToDB()
}

func main(){
	initializers.DB.AutoMigrate(&models.Users{})
	initializers.DB.AutoMigrate(&models.BackendTests{})
	initializers.DB.AutoMigrate(&models.Groups{})
	initializers.DB.AutoMigrate(&models.GroupsRelation{})
	initializers.DB.AutoMigrate(&models.SaveEndpointResult{}) //migrate schemas
}