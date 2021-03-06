package main

type jsonHash map[string]interface{}

func getBpRequest() jsonHash {
	json := jsonHash{}
	json["dataStreamName"] = "BloodPressure"
	json["type"] = "raw"
	json["application"] = jsonHash{
		"detailsUrl": "",
		"name":       "",
		"version":    "1",
	}
	json["dataType"] = jsonHash{
		"name": "om.google.blood_pressure",
	}
	return json
}
