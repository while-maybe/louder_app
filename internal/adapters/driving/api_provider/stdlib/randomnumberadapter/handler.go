package randomnumberadapter

import (
	"encoding/json"
	"errors"
	"log"
	stdlibapiadapter "louder/internal/adapters/driving/api_provider/stdlib"
	"louder/internal/core/domain"
	"louder/internal/core/service/randomnumbers"
	"strconv"
	"strings"

	"net/http"
)

// RandomNumber Handler
type RandomNumberHandler struct {
	RandomNumberService randomnumbers.Port // inject core service
}

// RandomDice Handler
type DiceRollHandler struct {
	RandomDiceService randomnumbers.Port
}

func NewRandomNumberHandler(service randomnumbers.Port) *RandomNumberHandler {
	return &RandomNumberHandler{RandomNumberService: service}
}

func NewRandomDiceHandler(service randomnumbers.Port) *DiceRollHandler {
	return &DiceRollHandler{RandomDiceService: service}
}

// HandleGetRandomNumber is an http.HandlerFunc for the /random route
func (h *RandomNumberHandler) HandleGetRandomNumber(w http.ResponseWriter, r *http.Request) {
	log.Println("stdlib API adapter: GET for /random")

	randomNumber := h.RandomNumberService.GetRandomNumber()
	response := RandomNumberResponse{RandomNumber: randomNumber}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode random number response %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// HandleGetDiceRoll is an http.HandlerFunc for the /diceroll route
func (h *DiceRollHandler) HandleGetDiceRoll(w http.ResponseWriter, r *http.Request) {
	log.Println("stdlib API adapter: GET for /diceroll")

	params := r.URL.Query()
	validationErrors := make([]string, 0)

	// convert the params to string first to remove \" if the user has provided a string as opposed to a number
	numDiceParam := strings.Trim(params.Get("numdice"), "\"` ")
	numSidesParam := strings.Trim(params.Get("numsides"), "\"` ")
	var numDice, numSides uint

	// check if required params exist
	if numDiceParam == "" {
		validationErrors = append(validationErrors, ErrMissingNumDice.Error())
	} else {
		val, err := strconv.Atoi(numDiceParam)
		switch {
		case err != nil:
			validationErrors = append(validationErrors, ErrFormatNumDice.Error())
		case val <= 0:
			// this test is needed in case of a negative int being later converted to uint
			validationErrors = append(validationErrors, ErrValueNumDice.Error())
		default:
			numDice = uint(val)
		}
	}

	if numSidesParam == "" {
		validationErrors = append(validationErrors, ErrMissingNumSides.Error())
	} else {
		val, err := strconv.Atoi(numSidesParam)
		switch {
		case err != nil:
			validationErrors = append(validationErrors, ErrFormatNumSides.Error())
		case val <= 0:
			validationErrors = append(validationErrors, ErrValueNumSides.Error())
		default:
			numSides = uint(val)
		}
	}

	if len(validationErrors) > 0 {
		stdlibapiadapter.RespondWithJSON(w, http.StatusBadRequest, stdlibapiadapter.ErrorResponse{ErrorMsgs: validationErrors})
		return
	}

	// Initially I was validating everything before giving any feedback to the user but AI suggested separating Transport and Domain validation separate so lets go with that

	diceRoll, err := h.RandomDiceService.RollDice(uint(numDice), uint(numSides))
	if err != nil {
		log.Printf("Service error during RollDice: %v", err)

		if errors.Is(err, domain.ErrInvalidNumDice) {
			validationErrors = append(validationErrors, domain.ErrInvalidNumDice.Error())
		}
		if errors.Is(err, domain.ErrInvalidNumSides) {
			validationErrors = append(validationErrors, domain.ErrInvalidNumSides.Error())
		}

		switch {
		case len(validationErrors) > 0:
			stdlibapiadapter.RespondWithJSON(w, http.StatusBadRequest, stdlibapiadapter.ErrorResponse{ErrorMsgs: validationErrors})
		default:
			stdlibapiadapter.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	// victory!
	response := DiceRollResponse{DiceRoll: *toRandomNumberDTO(diceRoll)}
	stdlibapiadapter.RespondWithJSON(w, http.StatusOK, response)
}
