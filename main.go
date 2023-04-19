package main

import (
	"encoding/json"
	"log"
	"net/http"

	"custom-validate/validate"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type AddSKPPRequest struct {
	PrakarsaId     string `json:"prakarsaId" validate:"required,number"`
	BranchCode     string `json:"branchCode" validate:"required,name"`
	Month          string `json:"month" validate:"required,name"`
	Year           string `json:"year" validate:"required,number"`
	SkppDate       string `json:"skppDate" validate:"required,dateString"`
	CustomerStatus string `json:"customerStatus" validate:"required,name"`
	BranchName     string `json:"branchName" validate:"required,name"`
	UserId         string `json:"userId" validate:"required,number"`
}

func main() {
	// HTTP
	router := http.NewServeMux()
	router.HandleFunc("/add-skpp", addSKPPHandler)

	// FRAMEWORK GIN
	r := gin.Default()
	r.POST("/add-skpp", AddSKPPHandlers)

	// FRAMEWORK ECHO
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/skpp", CreateSKPPHandler)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))

	log.Println("Starting server at :8080")

	// http.ListenAndServe(":8080", router)

	//start server gin
	// r.Run(":8080")

	// Start server echo
	e.Logger.Fatal(e.Start(":8080"))
}

func addSKPPHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody AddSKPPRequest

	// Parse the incoming request
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Println("Error parsing request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the request parameters
	status, errMsg := validate.CustomeValidateRequest(requestBody)
	if status != 200 {
		log.Println("Error validating request:", errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	// Process the request
	// TODO: Implement your business logic here

	// Generate the response
	response := map[string]string{"message": "SKPP added successfully"}
	jsonBytes, _ := json.Marshal(response)

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func AddSKPPHandlers(c *gin.Context) {
	var requestBody AddSKPPRequest

	// Parsing request body
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validating request body
	if validateStatus, validateMsg := validate.CustomeValidateRequest(requestBody); validateStatus != 200 {
		c.JSON(int(validateStatus), gin.H{"error": validateMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data added successfully"})
}

func CreateSKPPHandler(c echo.Context) error {
	req := new(AddSKPPRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Validating request body
	status, errMsg := validate.CustomeValidateRequest(req)
	if status != 200 {
		return c.JSON(int(status), errMsg)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"statusCode": status,
		"Message":    "SUCCESS",
	})
}
