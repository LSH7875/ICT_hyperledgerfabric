/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import ( //외부 패키지 목록
	"encoding/json" //json 마샬링
	"fmt"           //log
	"strconv"       //str 형태의 데이터를 다양한 자료형으로 변환 strconvert

	//위의 패키지는 golang 기본 패키지

	"github.com/hyperledger/fabric-contract-api-go/contractapi" //스마트 컨트랙트를 다루게 해주는 api가 모여져있다.
)

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}

/**2
* contractapi.Contract 라는 한마디를 씀으로써 Coontract를 모두 상속하는 것과 같다고 보면 됨.
* 이 한줄로 인해서 interface를 따르게 된다.
* contractapiContract설명참조
*func (c *Contract) GetAfterTransaction() interface{} // interface{} 뒤에 붙은 애들이 interface역할을 하는 함수이다.
func (c *Contract) GetBeforeTransaction() interface{}
func (c *Contract) GetInfo() metadata.InfoMetadata
func (c *Contract) GetName() string
func (c *Contract) GetTransactionContextHandler() SettableTransactionContextInterface
func (c *Contract) GetUnknownTransaction() interface{}
위의 함수들은 신경쓸 필요가 없고 InitLedger같은 우리가 직접 작성한 함수를 신경 쓰면 된다.
*/

// Car describes basic details of what makes up a car
type Car struct {
	Make   string `json:"make"`
	Model  string `json:"model"`
	Colour string `json:"colour"`
	Owner  string `json:"owner"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Car
}

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	cars := []Car{
		Car{Make: "Toyota", Model: "Prius", Colour: "blue", Owner: "Tomoko"},
		Car{Make: "Ford", Model: "Mustang", Colour: "red", Owner: "Brad"},
		Car{Make: "Hyundai", Model: "Tucson", Colour: "green", Owner: "Jin Soo"},
		Car{Make: "Volkswagen", Model: "Passat", Colour: "yellow", Owner: "Max"},
		Car{Make: "Tesla", Model: "S", Colour: "black", Owner: "Adriana"},
		Car{Make: "Peugeot", Model: "205", Colour: "purple", Owner: "Michel"},
		Car{Make: "Chery", Model: "S22L", Colour: "white", Owner: "Aarav"},
		Car{Make: "Fiat", Model: "Punto", Colour: "violet", Owner: "Pari"},
		Car{Make: "Tata", Model: "Nano", Colour: "indigo", Owner: "Valeria"},
		Car{Make: "Holden", Model: "Barina", Colour: "brown", Owner: "Shotaro"},
	}
	//총 10개의 구조체를 만들어주고 있다. slice로ㅠ 만들어서 Car라는 변수에 저장하고 있다.

	for i, car := range cars {
		carAsBytes, _ := json.Marshal(car)
		//carAsBytes에 직렬화된 데이터코드를 만드는 마샬링하는 과정임. putState를 넣는데 사용한다.
		//json.Marshal()함수가 Car구조체 타입을 알 필요가 없나요? 무슨 구조인지 알 수 있다.
		//뭐가 넣어질 지 알 수 밖에 없다. tag까지 달아서 아려준다. type Car struct{...}에서
		//태그가 키 값이다.
		err := ctx.GetStub().PutState("CAR"+strconv.Itoa(i), carAsBytes)
		//Car0, Make: "Toyota", Model: "Prius", Colour: "blue", Owner: "Tomoko"
		/**
		5
		원장에 값을 기록하는 역할을 한다.
		InitLedger는 err := ctx.GetStub().PutState("CAR"+strconv.Itoa(i), carAsBytes) 을 위한 함수라고도 볼 수 있다.
		ctx는  contractapi.TransactionContextInterface이다.
		InitLedger는 Dapp이 peer에게 요청해서 peer가 InitLedger를 호출하는 것임.
		*/

		//worldStateDB에 저장이 된다.
		//첫번째는 key가 value이다.
		//getState일 때는 ws값만 가지고 옴. 블록을 생성할 일이 없음.

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateCar adds a new car to the world state with given details
func (s *SmartContract) CreateCar(ctx contractapi.TransactionContextInterface, carNumber string, make string, model string, colour string, owner string) error {
	car := Car{
		Make:   make,
		Model:  model,
		Colour: colour,
		Owner:  owner,
	}
	//뒤에 인자값으로 받는다. peer가 준다. Dapp이 인자값을 넣어서 호출한다.

	carAsBytes, _ := json.Marshal(car)

	return ctx.GetStub().PutState(carNumber, carAsBytes)
}

// QueryCar returns the car stored in the world state with given id
func (s *SmartContract) QueryCar(ctx contractapi.TransactionContextInterface, carNumber string) (*Car, error) {
	carAsBytes, err := ctx.GetStub().GetState(carNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if carAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", carNumber)
	}

	car := new(Car)
	_ = json.Unmarshal(carAsBytes, car)

	return car, nil
}

// QueryAllCars returns all cars found in world state
func (s *SmartContract) QueryAllCars(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		car := new(Car)
		_ = json.Unmarshal(queryResponse.Value, car)
		//2줄은 언마샬한 결과물을 car에 넣어주는 것임.

		queryResult := QueryResult{Key: queryResponse.Key, Record: car}
		results = append(results, queryResult)
	}

	return results, nil
}

// ChangeCarOwner updates the owner field of car with given id in world state
func (s *SmartContract) ChangeCarOwner(ctx contractapi.TransactionContextInterface, carNumber string, newOwner string) error {
	car, err := s.QueryCar(ctx, carNumber)

	if err != nil {
		return err
	}

	car.Owner = newOwner

	carAsBytes, _ := json.Marshal(car)

	return ctx.GetStub().PutState(carNumber, carAsBytes)
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract)) //초기화 및 선언
	//1
	//contractapi.NewChaincode :peer가 자신의 메모리 공간 위에 올려서 돌려야할 cc를 실행 가능한 객체 형식으로 만들어 주는 역할을 한다.
	//실행 가능한 체인 코드를 만들어주는 역할

	//err: 만약 체인코드 실행 중에 오류가 난다면  어떤 에러가 났는지 err로 밖으로 뱉어줌.
	// 잘 실행되면 err에 값이 없음. err 값이 있으면 실행가능한 cc로 만드는데 에러가 있다는 것임.

	//new:  빈껍데기인 구조체 생성하는 것임.
	//contractapi.NewChaincode는 ibm이 개발함. 그런데 SmartContract는 cc개발자가 만들었다.
	//cc개발자가 미리 newChaincode의 자료형을 알고 있어야함? 그래서 인터페이스 형식으로 되어있어서 넘길 수 있음.
	//패브릭 개발자가 미리 정의한 인터페이스를 구현하고 있으니까 인자로 전달이 가능하다.  그 인터페이스를 따르면 된다.
	//인터페이스를 따르는 어떠한 구조체든지 다 넣을 수 있다.
	//https://pkg.go.dev/github.com/hyperledger/fabric-contract-api-go/contractapi#NewChaincode

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}
	/**
	3
	go 언어에서 nil이 null
	에러 체크 해주겠다. */

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
	/**
	4
	만들어진 chaincode.Start라는 함수를 때리면 start를 하는 순간 peer위에 cc가 올라가서 돌아간다.
	이 때 실패해서 에러가 발생해서 오류가 발생할 수 있다.
	위의 err != nil을 통과했으면 err엔 nil값이어서 err := 를 쓸 수 있다.
	우리가 결국 만들 수 잇는 건 SmartContract밖에 없다. 누가 짜든 main함수는 변하지 않는다.
	cc에 이름을 정해 배포하고 실행파일을 실행시키면start()가동작해서 호출을 리스닝하고 있다고 생각하면 됨.
	chaincode 가 sdk */

	//5가지 함수
}
