package tests

import (
	"fmt"
	"labs/domains"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database represents the database instance for the tests package.
type Database struct {
	DB *gorm.DB
}

// NewUserRepository performs automatic migration of test-related structures in the database.
func NewTestRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.Tests{}, &domains.TestQuestions{}, &domains.TestCandidats{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the test structure. Error: ", err)
	}
}

// ReadAllPagination retrieves a paginated list of tests based on company ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.Tests, modelID uuid.UUID, limit, offset int) ([]domains.Tests, error) {
	err := db.Where("company_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadAllPaginationQ retrieves a paginated list of Questions based on test ID, limit, and offset.
func ReadAllPaginationQ(db *gorm.DB, model []domains.TestQuestions, modelID uuid.UUID, limit, offset int) ([]domains.TestQuestions, error) {
	err := db.Where("tests_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadAllPagination retrieves a paginated list of tests based on company ID, limit, and offset.
func ReadAllPaginationC(db *gorm.DB, model []domains.Condidats, modelID uuid.UUID, limit, offset int) ([]domains.Condidats, error) {
	err := db.Where("id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadAllList retrieves a list of tests based on company ID.
func ReadAllList(db *gorm.DB, model []domains.Tests, modelID uuid.UUID) ([]domains.Tests, error) {
	err := db.Where("company_id = ? ", modelID).Find(&model).Error
	return model, err
}

// ReadTechList retrieves the list of technologies associated with a test based on the test ID.
func ReadTechList(db *gorm.DB, testID uuid.UUID) (pq.StringArray, error) {
	var technologies pq.StringArray
	err := db.Model(&domains.Tests{}).Where("id = ?", testID).Pluck("technologies", &technologies).Error
	return technologies, err
}

// ReadQuestionsList retrieves a list of questions based on test ID.
func ReadQuestionsList(db *gorm.DB, model []domains.TestQuestions, modelID uuid.UUID) ([]domains.TestQuestions, error) {
	err := db.Where("test_id = ? ", modelID).Find(&model).Error
	return model, err
}

// GetAllQuestion IDs that associated with same technologies
func GetAllQuestions(db *gorm.DB, technologies pq.StringArray) (pq.StringArray, error) {
	// Define a slice to hold the IDs of the questions
	var questionIDs pq.StringArray

	// Iterate over each technology
	for _, techString := range technologies {
		// Split the techString into separate technologies
		techs := strings.Split(strings.Trim(techString, "{}"), ",")
		fmt.Println("techs:", techs)
		// Iterate over each individual technology
		for _, tech := range techs {
			// Define a slice to hold the IDs of questions for the current technology
			var techQuestionIDs pq.StringArray

			// Trim any extra spaces from the technology string
			tech = strings.TrimSpace(tech)

			// Query the database to find questions associated with the current technology
			err := db.Model(&domains.Questions{}).Where("associated_technology = ?", tech).Pluck("id", &techQuestionIDs).Error
			if err != nil {
				return nil, err
			}

			// Append the IDs of questions for the current technology to the main list
			questionIDs = append(questionIDs, techQuestionIDs...)
		}
	}

	return questionIDs, nil
}

// GetAllCandidats IDs that belongs to a specific test
func GetAllCandidates(testID uuid.UUID, db *gorm.DB) ([]string, error) {
	// Define a slice to hold the IDs of the candidates
	var candidateIDs []string

	// Query the database to find candidates associated with the current test ID
	err := db.Model(&domains.TestCandidats{}).
		Where("tests_id = ?", testID).Pluck("candidat_id", &candidateIDs).Error
	if err != nil {
		return nil, err
	}

	return candidateIDs, nil
}

// GetRandomQuestions retrieves a random selection of question IDs from the provided list.
func GetRandomQuestions(questionIDs pq.StringArray, numQuestions int) (selectedQuestions pq.StringArray, err error) {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Get the length of the questionIDs slice
	numIDs := len(questionIDs)

	// Check if the number of requested questions is greater than the available IDs
	if numQuestions > numIDs {
		return nil, fmt.Errorf("number of requested questions (%d) exceeds available IDs (%d)", numQuestions, numIDs)
	}

	// Create a map to track selected indices
	selectedIndices := make(map[int]bool)

	// Create a slice to hold the selected question IDs
	selectedQuestions = make(pq.StringArray, 0, numQuestions)

	// Loop until we have selected the desired number of questions
	for len(selectedQuestions) < numQuestions {
		// Generate a random index
		randIndex := rand.Intn(numIDs)

		// Check if the index is already selected
		if _, ok := selectedIndices[randIndex]; !ok {
			// Add the index to the map of selected indices
			selectedIndices[randIndex] = true

			// Append the corresponding question ID to the selected questions slice
			selectedQuestions = append(selectedQuestions, questionIDs[randIndex])
		}
	}

	return selectedQuestions, nil
}

// ReadByID retrieves a test by their unique identifier.
func ReadByID(db *gorm.DB, model domains.Tests, id uuid.UUID) (domains.Tests, error) {
	err := db.First(&model, id).Error
	return model, err
}

/*func GetAllQuestionsDetails(db *gorm.DB, questionIDs pq.StringArray) ([]domains.Questions, error) {
	// Define a slice to hold the details of the questions
	var questions []domains.Questions

	// Iterate over each string containing question IDs
	for _, idQString := range questionIDs {
		// Split the idQString into separate question IDs
		ids := strings.Split(strings.Trim(idQString, "{}"), ",")

		// Iterate over each individual question ID
		for _, id := range ids {
			// Trim any extra spaces from the question ID string
			id = strings.TrimSpace(id)

			// Define a variable to hold the details of the current question
			var question domains.Questions

			// Query the database to find details of the question with the current ID
			err := db.Where("id = ?", id).Find(&question).Error
			if err != nil {
				return nil, err
			}

			// Append the details of the current question to the main list
			questions = append(questions, question)
		}
	}

	return questions, nil
}*/
