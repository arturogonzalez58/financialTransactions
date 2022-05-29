package main

import (
	"fmt"
	"github.com/arturogonzalez58/financialTransactions/dataGenerator/pkg/csv"
	"github.com/arturogonzalez58/financialTransactions/dataGenerator/pkg/generator"
	"github.com/arturogonzalez58/financialTransactions/dataGenerator/pkg/s3"
	"github.com/google/uuid"
	"os"
	"strconv"
	"time"
)

func printUsage() {
	fmt.Println("dataGenerator NUMBER_OF_DATA ERROR_PERCENTAGE")
	fmt.Println("example:")
	fmt.Println("dataGenerator 1000 20")
	fmt.Println("NUMBER_OF_DATA: The number of random data to generate")
	fmt.Println("ERROR_PERCENTAGE: The percentage of the data to be errors")
	fmt.Println("BUCKET: The bucket to save the output file")
}

func main() {
	arguments := os.Args[1:]
	if len(arguments) < 3 {
		printUsage()
		os.Exit(1)
	}

	NUMBER_OF_DATA, err := strconv.ParseInt(arguments[0], 10, 32)
	if err != nil {
		fmt.Println("There is a problem with the NUMBER_OF_DATA")
		os.Exit(2)
	}
	ERROR_PERCENTAGE, err := strconv.ParseInt(arguments[1], 10, 32)
	if err != nil {
		fmt.Println("There is a problem with the ERROR_PERCENTAGE: %w", err)
		os.Exit(2)
	}

	AwsRegion := os.Getenv("AWS_REGION")
	S3Bucket := arguments[2]

	s3Uploader, err := s3.Build(AwsRegion, S3Bucket)
	if err != nil {
		fmt.Println("Unable to create a session with aws: %w", err)
		os.Exit(2)
	}

	initialDate := time.Date(2021, 1, 0, 0, 0, 0, 0, time.UTC)
	finalDate := time.Date(2022, 1, 0, 0, 0, 0, 0, time.UTC)
	data := generator.Builder(int32(NUMBER_OF_DATA), float32(ERROR_PERCENTAGE/100.00), initialDate, finalDate).GenerateData()
	dataToSave, err := csv.Build(data).ToCsv()
	if err != nil {
		fmt.Println("there was a problem creating the csv data:", err)
		os.Exit(2)
	}

	fileName := fmt.Sprintf("%s.csv", uuid.New())
	err = s3Uploader.AddFileToS3(dataToSave, fileName)
	if err != nil {
		fmt.Println("there was a problem saving the data in the bucket:", err)
		os.Exit(2)
	}
	fmt.Println("The data has been generate and saved in the bucket")
	os.Exit(0)
}
