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

// ReadAllPagination retrieves a paginated list of tests based on company ID, limit, and offset.
func ReadAllPaginationS(db *gorm.DB, model []domains.TestCandidats, modelID uuid.UUID, limit, offset int) ([]domains.TestCandidats, error) {
	err := db.Where("company_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadAllPaginationQ retrieves a paginated list of Questions based on test ID, limit, and offset.
func ReadAllPaginationQ(db *gorm.DB, model []domains.TestQuestions, modelID uuid.UUID, limit, offset int) ([]domains.TestQuestions, error) {
	err := db.Where("candidat_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadAllPagination retrieves a paginated list of tests based on company ID, limit, and offset.
func ReadAllPaginationC(db *gorm.DB, model []domains.Condidats, modelID uuid.UUID, limit, offset int) ([]domains.Condidats, error) {
	err := db.Where("id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadAllPagination retrieves a paginated list of questions based on test ID, limit, and offset.
func ReadAllPaginationQT(db *gorm.DB, model []domains.TestQuestions, modelID uuid.UUID, limit, offset int) ([]domains.TestQuestions, error) {
	err := db.Where("tests_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadAllPagination retrieves a paginated list of questions based on test ID, limit, and offset.
func ReadAllPaginationAnswers(db *gorm.DB, model []domains.TestQuestions, testID uuid.UUID, candidatID uuid.UUID, limit, offset int) ([]domains.TestQuestions, error) {
	err := db.Where("tests_id = ? and candidat_id = ? ", testID, candidatID).Limit(limit).Offset(offset).Find(&model).Error
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

	// Shuffle the questionIDs slice
	shuffledIDs := make(pq.StringArray, len(questionIDs))
	copy(shuffledIDs, questionIDs)
	rand.Shuffle(len(shuffledIDs), func(i, j int) {
		shuffledIDs[i], shuffledIDs[j] = shuffledIDs[j], shuffledIDs[i]
	})

	// Select the first numQuestions IDs from the shuffled slice
	selectedQuestions = shuffledIDs[:numQuestions]

	return selectedQuestions, nil
}

// ReadByID retrieves a test by their unique identifier.
func ReadByID(db *gorm.DB, model domains.Tests, id uuid.UUID) (domains.Tests, error) {
	err := db.First(&model, id).Error
	return model, err
}

// ReadCandidatByID retrieves a Candidat by their unique identifier.
func ReadCandidatByID(db *gorm.DB, model domains.Condidats, id uuid.UUID) (domains.Condidats, error) {
	err := db.First(&model, id).Error
	return model, err
}



// ReadProjectDetails retrive a project based on projectID.
func ReadProjectDetails(db *gorm.DB, model domains.Project, modelID uuid.UUID) (domains.Project, error) {
	err := db.Where("id = ? ", modelID).First(&model).Error
	return model, err
}



// ReadProjectIDByCandidatID retrieves the last projectID based on the candidatID.
func ReadProjectIDByCandidatID(db *gorm.DB, candidatID uuid.UUID) (uuid.UUID, error) {
	var projectsCondidats domains.ProjectsCondidats // Assuming ProjectsCondidats is your model struct
	err := db.Where("condidat_id = ?", candidatID).Last(&projectsCondidats).Error
	if err != nil {
		return uuid.Nil, err
	}
	return projectsCondidats.ProjectID, nil
}

// ReadQuestionDetails retrive a question based on questionID.
func ReadQuestionDetails(db *gorm.DB, model domains.Questions, modelID uuid.UUID) (domains.Questions, error) {
	err := db.Where("id = ? ", modelID).First(&model).Error
	return model, err
}
