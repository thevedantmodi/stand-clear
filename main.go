package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thevedantmod/stand-clear/board"
)

func main() {
	router := gin.Default()

	router.GET("/arrivals", func(ctx *gin.Context) {
		fmt.Println(ctx.Params)
		line := ctx.Query("line")
		stop_id := ctx.Query("stop_id")
		Nstr := ctx.DefaultQuery("N", "10")
		N, err := strconv.Atoi(Nstr)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		arrivals := board.GetArrivals(line, stop_id, N)
		ctx.IndentedJSON(http.StatusOK, arrivals)
	})

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello from the board!")
	})
	router.Run("localhost:8080")
	// file, err := os.Open("subway_config/routes.txt")

	// defer file.Close()
	// // Create a new CSV reader
	// reader := csv.NewReader(file)

	// // Iterate over each record (line) in the CSV file
	// for {
	// 	record, err := reader.Read() // Read one record (line)
	// 	if err == io.EOF {
	// 		break // End of file reached
	// 	}
	// 	if err != nil {
	// 		log.Fatalf("Error reading CSV record: %s", err)
	// 	}

	// 	// Print the entire record (slice of strings)
	// 	fmt.Println(record[0], record[3])
	// }
	// board.GetArrivals("6", "627N", 10)
	// board.GetArrivals("7", "725N", 10)
	// x``
}
