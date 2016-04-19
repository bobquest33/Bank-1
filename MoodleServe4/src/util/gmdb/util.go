package gmdb

import (
	"strconv"
	"errors"
)

func MapI2MapS(fv map[string]interface{}) (map[string]string, error) {
	res := make(map[string]string)
	for k, v := range fv {
		switch v.(type) {
		case string:
			res[k] = v.(string)
			break
		case int:
			tyv := strconv.Itoa(v.(int))
			res[k] = tyv
			break
		case float32:
			tyv := strconv.FormatFloat(float64(v.(float32)), 'e', -1, 32)
			res[k] = tyv
			break
		case float64:
			tyv := strconv.FormatFloat(v.(float64), 'e', -1, 32)
			res[k] = tyv
			break
		case bool:
			tyv := strconv.FormatBool(v.(bool))
			res[k] = tyv
			break
		default:
			return nil, errors.New("Unknow type of the value")
		}
	}
	return res, nil
}

func collectFieldValue()  {

}