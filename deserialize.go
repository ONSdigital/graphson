package graphson

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

// DeserializeVertices converts a graphson string to a slice of Vertex
func DeserializeVertices(rawResponse string) ([]Vertex, error) {
	// TODO: empty strings for property values will cause invalid json
	// make so it can handle that case
	if len(rawResponse) == 0 {
		return []Vertex{}, nil
	}
	return DeserializeVerticesFromBytes([]byte(rawResponse))
}

// DeserializeVerticesFromBytes returns a slice of Vertex from the graphson rawResponse list of vertex
func DeserializeVerticesFromBytes(rawResponse []byte) ([]Vertex, error) {
	// TODO: empty strings for property values will cause invalid json
	// make so it can handle that case
	var response []Vertex
	if len(rawResponse) == 0 {
		return response, nil
	}
	dec := json.NewDecoder(bytes.NewReader(rawResponse))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

// DeserializeListFromBytes validates that the provided RawMessage ([]byte) corresponds to a g:List type and returns a slice of values as RawMessages
func DeserializeListFromBytes(rawResponse json.RawMessage) ([]json.RawMessage, error) {

	if isEmptyResponse(rawResponse) {
		return []json.RawMessage{}, nil
	}

	var metaResponse RawSlice

	dec := json.NewDecoder(bytes.NewReader(rawResponse))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&metaResponse); err != nil {
		return nil, err
	}

	if metaResponse.Type != "g:List" {
		return nil, fmt.Errorf("DeserializeListFromBytes: Expected `g:List` type, but got %q", metaResponse.Type)
	}

	return metaResponse.Value, nil
}

// DeserializeListOfVerticesFromBytes returns a slice of Vertex from the graphson rawResponse g:List of vertex
func DeserializeListOfVerticesFromBytes(rawResponse []byte) ([]Vertex, error) {

	if isEmptyResponse(rawResponse) {
		return []Vertex{}, nil
	}

	var metaResponse ListVertices
	var response []Vertex

	dec := json.NewDecoder(bytes.NewReader(rawResponse))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&metaResponse); err != nil {
		return nil, err
	}

	if metaResponse.Type != "g:List" {
		return response, fmt.Errorf("DeserializeListOfVerticesFromBytes: Expected `g:List` type, but got %q", metaResponse.Type)
	}

	return metaResponse.Value, nil
}

func DeserializeListOfEdgesFromBytes(rawResponse []byte) (Edges, error) {

	if isEmptyResponse(rawResponse) {
		return Edges{}, nil
	}

	var metaResponse ListEdges
	var response Edges

	dec := json.NewDecoder(bytes.NewReader(rawResponse))
	dec.DisallowUnknownFields()
	err := dec.Decode(&metaResponse)
	if err != nil {
		return nil, err
	}

	if metaResponse.Type != "g:List" {
		return response, fmt.Errorf("DeserializeListOfEdgesFromBytes: Expected `g:List` type, but got %q", metaResponse.Type)
	}

	return metaResponse.Value, nil
}

func DeserializeMapFromBytes(rawResponse []byte) (resMap map[string]json.RawMessage, err error) {
	resMap = make(map[string]json.RawMessage)

	if isEmptyResponse(rawResponse) {
		return resMap, nil
	}

	var metaResponse RawSlice

	dec := json.NewDecoder(bytes.NewReader(rawResponse))
	dec.DisallowUnknownFields()
	if err = dec.Decode(&metaResponse); err != nil {
		return nil, err
	}

	if metaResponse.Type != "g:Map" {
		return nil, fmt.Errorf("DeserializeMapFromBytes: Expected `g:Map` type, but got %q", metaResponse.Type)
	}

	if len(metaResponse.Value)%2 != 0 {
		return nil, errors.New("odd number of values, maps should have an even number of values to construct key-value pairs")
	}

	for i := 0; i < len(metaResponse.Value); i += 2 {
		var resKey string
		if err := json.Unmarshal(metaResponse.Value[i], &resKey); err != nil {
			return nil, err
		}
		resMap[resKey] = metaResponse.Value[i+1]
	}

	return resMap, nil
}

// DeserializePropertiesFromBytes is for converting vertex .properties() results into a map
func DeserializePropertiesFromBytes(rawResponse []byte, resMap map[string][]interface{}) (err error) {
	var metaResponse Raw
	if len(rawResponse) == 0 {
		return
	}
	dec := json.NewDecoder(bytes.NewReader(rawResponse))
	dec.DisallowUnknownFields()
	if err = dec.Decode(&metaResponse); err != nil {
		return
	}

	if metaResponse.Type != "g:List" {
		return fmt.Errorf("DeserializePropertiesFromBytes: Expected `g:List` type, but got %q", metaResponse.Type)
	}
	var props []VertexProperty
	if err = json.Unmarshal(metaResponse.Value, &props); err != nil {
		return
	}

	for _, prop := range props {
		if _, ok := resMap[prop.Value.Label]; !ok {
			resMap[prop.Value.Label] = []interface{}{prop.Value.Value}
		} else {
			resMap[prop.Value.Label] = append(resMap[prop.Value.Label], prop.Value.Value)
		}
	}

	return
}

// DeserializeStringListFromBytes get a g:List value which should be a a list of strings, return those
func DeserializeStringListFromBytes(rawResponse []byte) (vals []string, err error) {

	if isEmptyResponse(rawResponse) {
		vals = []string{}
		return
	}

	dec := json.NewDecoder(bytes.NewReader(rawResponse))
	dec.DisallowUnknownFields()

	var metaResponse Raw
	if err = dec.Decode(&metaResponse); err != nil {
		return
	}

	if metaResponse.Type != "g:List" {
		err = fmt.Errorf("DeserializeStringListFromBytes: Expected `g:List` type, but got %q", metaResponse.Type)
		return
	}

	if err = json.Unmarshal(metaResponse.Value, &vals); err != nil {
		return
	}
	return
}

// DeserializeSingleFromBytes get a g:List value which should be a singular item, returns that item
func DeserializeSingleFromBytes(rawResponse []byte) (gV GenericValue, err error) {
	var metaResponse Raw
	if len(rawResponse) == 0 {
		err = errors.New("DeserializeSingleFromBytes: nothing to decode")
		return
	}
	dec := json.NewDecoder(bytes.NewReader(rawResponse))
	dec.DisallowUnknownFields()
	if err = dec.Decode(&metaResponse); err != nil {
		return
	}

	if metaResponse.Type != "g:List" {
		err = fmt.Errorf("DeserializeSingleFromBytes: Expected `g:List` type, but got %q", metaResponse.Type)
		return
	}

	var genVals GenericValues
	if genVals, err = DeserializeGenericValues(string(metaResponse.Value)); err != nil {
		return
	}

	if len(genVals) != 1 {
		err = fmt.Errorf("DeserializeSingleFromBytes: Expected single value, got %d", len(genVals))
		return
	}

	return genVals[0], nil
}

// DeserializeNumber returns the count from the g:List'd database response
func DeserializeNumber(rawResponse []byte) (count int64, err error) {
	var genVal GenericValue
	if genVal, err = DeserializeSingleFromBytes(rawResponse); err != nil {
		return
	}

	if genVal.Type != "g:Int64" {
		err = fmt.Errorf("DeserializeNumber: Expected `g:Int64` type, but got %q", genVal.Type)
		return
	}
	count = int64(genVal.Value.(float64))
	return
}

func DeserializeInt32(rawResponse json.RawMessage) (num int32, err error) {
	var genVal GenericValue
	if genVal, err = DeserializeSingleFromBytes(rawResponse); err != nil {
		return
	}

	if genVal.Type != "g:Int32" {
		err = fmt.Errorf("DeserializeNumber: Expected `g:Int32` type, but got %q", genVal.Type)
		return
	}
	num = int32(genVal.Value.(float64))
	return
}

func DeserializeEdges(rawResponse string) (Edges, error) {
	var response Edges
	if rawResponse == "" {
		return response, nil
	}
	err := json.Unmarshal([]byte(rawResponse), &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func DeserializeGenericValue(rawResponse string) (response GenericValue, err error) {
	if len(rawResponse) == 0 {
		return
	}
	if err = json.Unmarshal([]byte(rawResponse), &response); err != nil {
		return
	}
	return
}

func DeserializeGenericValues(rawResponse string) (GenericValues, error) {
	var response GenericValues
	if rawResponse == "" {
		return response, nil
	}
	err := json.Unmarshal([]byte(rawResponse), &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func ConvertToCleanVertices(vertices []Vertex) []CleanVertex {
	var responseVertices []CleanVertex
	for _, vertex := range vertices {
		responseVertices = append(responseVertices, CleanVertex{
			Id:    vertex.Value.ID,
			Label: vertex.Value.Label,
		})
	}
	return responseVertices
}

func ConvertToCleanEdges(edges Edges) []CleanEdge {
	var responseEdges []CleanEdge
	for _, edge := range edges {
		responseEdges = append(responseEdges, CleanEdge{
			Source: edge.Value.InV,
			Target: edge.Value.OutV,
		})
	}
	return responseEdges
}

func isEmptyResponse(rawResponse []byte) bool {
	return len(rawResponse) == 0 || isNullResponse(rawResponse)
}

func isNullResponse(rawResponse []byte) bool {
	return len(rawResponse) == 4 && string(rawResponse) == "null"
}
