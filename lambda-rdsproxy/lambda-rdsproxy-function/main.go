package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db               *sql.DB
	err              error
	connectionString string
	dbUser           string
	dbPass           string
	dataSource       string
)

type Employee struct {
	EmployeeNo int    `json:"emp_no"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
}

func init() {
	connectionString = os.Getenv("CONN")
	dbUser = os.Getenv("DBUSER")
	dbPass = os.Getenv("DBPASS")
	dataSource = dbUser + ":" + dbPass + "@tcp(" + connectionString + "/employees"
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	searchCriteria := request.Body

	db, err = sql.Open("mysql", dataSource)
	if err != nil {
		panic(err.Error()) // Just for example purpose. Use proper error handling instead of panic
	}

	defer db.Close()

	results, err := db.Query("select e.emp_no, e.first_name, e.last_name " +
		"from employees e, departments d, dept_emp de " +
		"where de.emp_no = e.emp_no " +
		"and de.dept_no = d.dept_no " +
		"and d.dept_name = 'Marketing' " +
		"and e.last_name LIKE '" + searchCriteria + "%';")

	if err != nil {
		log.Fatal(err)
	}
	defer results.Close()

	employees := make([]Employee, 0)

	for results.Next() {
		var e Employee

		err := results.Scan(&e.EmployeeNo, &e.FirstName, &e.LastName)
		if err != nil {
			log.Fatal(err)
		}
		employees = append(employees, e)
	}

	data, _ := json.Marshal(employees)

	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Body:            string(data),
		IsBase64Encoded: false,
	}, nil
}

func main() {
	lambda.Start(handler)
}
