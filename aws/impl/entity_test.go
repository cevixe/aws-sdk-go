package impl

import (
	"fmt"
	"github.com/cevixe/aws-sdk-go/aws/serdes/json"
	"testing"
)

func TestBug(t *testing.T) {
	jsonString := "{\"EventRecord\":null,\"StateRecord\":{\"type\":\"ChannelGroup\",\"id\":\"1\",\"version\":1,\"updated_at\":1635530624741,\"updated_by\":\"Unknown\",\"created_at\":1635530624741,\"created_by\":\"Unknown\",\"content_type\":\"application/json\",\"content_encoding\":\"gzip\",\"content\":\"H4sIAAAAAAAA/+xXXW/qOBD9K5afySXfgbwNiS9NFXBlQ7XV7qoyiWlRS8ImYVeo4r+vktDSrknp9r6StySeM+Mzxz6aFyz/lll1X+bbIpHYx/00X4tV1k8eRZbJZ+2hyLebvoF7h4WrFPtYV5/jiuRZlCX2cfj2pdptauighRzXiEEhRSXT45LVWmLfcC3HsXTXtD3beP2Vih32sambhmbomjl8ixHb6jEvsI/n2VOW/5PhY0QlsP+CD3sosf/7C14U+ZMsms95KlvEunDT1HEPp7JMitWmWuUZ9jHQKbohbI4CyhgJKSMchQRxMp4zyhH/AT/wft87CWo4DR0DFZRFwNEfW103XdS+dMMHnyRwm7I9JUEALKYcQTwibEZRAIwBDygaA2NRSM8AGkMVkE45uSZ1TXfva737IhPDGtczT+Ce2vZPiGFE4hjOEdC0TR+oBIzr+O+R2mKajoJ5RdiEcASc8MPev9e0VhW2gn8dz74vNbdDajEgmP8WxRGElEEXNARnytVVOibA+BVi5CokrCnuNHg3sNehtRtGbqPX0kjE4rOUd+cYNL20lByM3DWtpEFNzNtZ7KCHBhEJIUQwpdNoAiggjEEI57rhqmlpAOjnnLAZ+QVxOrZKGqOTzu5+5RIx1DPECTv2gUPwP7k3W0WqkrwFNgaOJhDDGNA1RBNy7rCo9/ItYTDlnwv6z/8GjZoEJZrlG9zDZSWqbW1PeN/DMqtW1a71tMbA2vcTfnX8WQPIi7tc3OXiLhd3ubjLxV1OuguZwigm4TuTSdqJ515U6pzzccVi92GmqYrVw4Ms3o1oSb5eiyzttwHah1FNK+RfW1m2s9VraDu0Oam3cBZCG8pkqNmD1NKEY9la4pmetBYitWyriRFZKZLD9gzNNbzE8jypWUuRLuxl6iVWknhL1156Rup4Ft7/GwAA//99xe+RRQ4AAA==\"}}"
	entity := &EntityImpl{}
	json.Unmarshall([]byte(jsonString), entity)

	entityMap := make(map[string]interface{})
	entity.State(&entityMap)
	entityMap["__typename"] = entity.Type()
	entityMap["_id"] = entity.ID()
	entityMap["_type"] = entity.Type()
	entityMap["_version"] = entity.Version()
	entityMap["_updatedAt"] = entity.UpdatedAt()
	entityMap["_updatedBy"] = entity.UpdatedBy()
	entityMap["_createdAt"] = entity.CreatedAt()
	entityMap["_createdBy"] = entity.CreatedBy()
	fmt.Println(string(json.Marshall(entityMap)))
}
