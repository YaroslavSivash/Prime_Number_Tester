package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"time"
)

type NumbersServiceProvider interface {
	CheckNumbersForPrimeNumber(ctx context.Context, numbersSlice []uint64) ([]bool, error)
}

type NumbersService struct {
	primesMap map[int64]bool
	maxValue  uint64
}

func NewNumbersService(filePath string) *NumbersService {
	ns := &NumbersService{}
	err := ns.parseJson(filePath)
	if err != nil {
		log.Println("cannot parse json")
	}
	return ns
}

func (s *NumbersService) CheckNumbersForPrimeNumber(ctx context.Context, numbersSlice []uint64) ([]bool, error) {

	// make copy slice
	numsSliceCp := make([]uint64, len(numbersSlice))
	copy(numsSliceCp, numbersSlice)

	// sort slice with numbers from smallest to largest year
	sort.Slice(numbersSlice, func(i, j int) bool {
		return numbersSlice[i] < numbersSlice[j]
	})

	// get max number in slice
	maxValue := numbersSlice[len(numbersSlice)-1]

	boolSlice := make([]bool, len(numbersSlice))
	// if max number in slice is less or equals max number in map, range on slice
	if maxValue <= s.maxValue {
		for idx, numSliceCp := range numsSliceCp {
			_, exist := s.primesMap[int64(numSliceCp)]
			if !exist {
				boolSlice[idx] = false
			} else {
				boolSlice[idx] = true
			}
		}
	} else {
		// else use an algorithm Eratosthene to find all prime numbers up to the number we need and store in map
		boolNumsSlice := make([]bool, maxValue)
		var i int64
		for i = 2; i <= int64(math.Sqrt(float64(maxValue)+1)); i++ {
			if boolNumsSlice[i] == false {
				for j := i * i; j < int64(maxValue); j += i {
					boolNumsSlice[j] = true
				}
			}
		}
		// add to map new prime numbers in range to max value
		for idx, isComposite := range boolNumsSlice {
			if idx > 1 && !isComposite {
				s.primesMap[int64(idx)] = true
			}
		}
		// add bool value for request numbers, if prime - true, else false
		for idx, numSliceCp := range numsSliceCp {
			_, exist := s.primesMap[int64(numSliceCp)]
			if !exist {
				boolSlice[idx] = false
			} else {
				boolSlice[idx] = true
			}
		}
		s.maxValue = maxValue
	}
	fmt.Println(boolSlice)
	return boolSlice, nil
}

func (s *NumbersService) parseJson(filePath string) error {
	now := time.Now()
	defer func() {
		log.Println("primes numbers download to map:", time.Now().Sub(now))
	}()

	// Open our json file
	jsonFile, err := os.Open(filePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}
	fmt.Println("Successfully opened primes.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened json file as a byte array.
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	// unmarshal bytes to map
	primesMapParse := make(map[string]int64)
	if err = json.Unmarshal(byteValue, &primesMapParse); err != nil {
		return err
	}

	maxValue := uint64(primesMapParse[strconv.Itoa(len(primesMapParse)-1)])
	// create map where is key is prime number
	primesMap := make(map[int64]bool, len(primesMapParse))
	for _, primeMapParse := range primesMapParse {
		primesMap[primeMapParse] = true
	}
	s.maxValue = maxValue
	s.primesMap = primesMap
	return nil
}
