package graphson

import (
	"encoding/json"
	"strconv"
)

func MakeDummyVertexProperty(label string, value interface{}) VertexProperty {
	return VertexProperty{
		Type: "g:VertexProperty",
		Value: VertexPropertyValue{
			ID: GenericValue{
				Type:  "Type",
				Value: 1,
			},
			Value: value,
			Label: label,
		},
	}
}

func MakeDummyVertex(vertexID, vertexLabel string, params map[string]interface{}) Vertex {
	properties := make(map[string][]VertexProperty)
	for label, value := range params {
		var vp []VertexProperty
		vSlice, ok := value.([]interface{})
		if ok {
			for _, p := range vSlice {
				vertexProperty := MakeDummyVertexProperty(label, p)
				vp = append(vp, vertexProperty)
			}
		} else {
			vertexProperty := MakeDummyVertexProperty(label, value)
			vp = append(vp, vertexProperty)
		}
		properties[label] = vp
	}
	vertexValue := VertexValue{
		ID:         vertexID,
		Label:      vertexLabel,
		Properties: properties,
	}
	return Vertex{
		Type:  "g:Vertex",
		Value: vertexValue,
	}
}

func MakeDummyProperty(label string, value GenericValue) EdgeProperty {
	jsonVal, err := json.Marshal(value)
	if err != nil {
		panic(err) // XXX
	}
	return EdgeProperty{
		Type: "g:Property",
		Value: EdgePropertyValue{
			Value: jsonVal,
			Label: label,
		},
	}
}

func MakeDummyEdge(edgeID int, edgeLabel, inVLabel, outVLabel string, inV, outV string, params map[string]int32) Edge {
	properties := make(map[string]EdgeProperty)
	for label, value := range params {
		properties[label] = MakeDummyProperty(label, MakeDummyGenericValue("g:Int32", value))
	}
	edgeIDStr := strconv.Itoa(edgeID)
	edgeValue := EdgeValue{
		ID:         edgeIDStr, // GenericValue{Type: "g:Int32", Value: edgeID},
		Label:      edgeLabel,
		InVLabel:   inVLabel,
		OutVLabel:  outVLabel,
		InV:        inV,  //GenericValue{Type: "g:Int32", Value: inV},
		OutV:       outV, //GenericValue{Type: "g:Int32", Value: outV},
		Properties: properties,
	}
	return Edge{
		Type:  "g:Edge",
		Value: edgeValue,
	}
}

func MakeDummyGenericValue(gvType string, value interface{}) GenericValue {
	return GenericValue{
		Type:  gvType,
		Value: value,
	}
}

// func MakeDummyMap(gvType string, value )
