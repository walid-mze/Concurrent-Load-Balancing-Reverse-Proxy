package main
import(
	"fmt"
	"log"
	"net/http"
)
func handler(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("hello from backend server"))
}
func main(){
	http.HandleFunc("/",handler)
	fmt.Println("backend server running on :8081")
	log.Fatal(http.ListenAndServe(":8081",nil))
}