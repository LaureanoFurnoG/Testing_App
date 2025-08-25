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
	initializers.DB.AutoMigrate(&models.Users{}) //migrate schemas
}