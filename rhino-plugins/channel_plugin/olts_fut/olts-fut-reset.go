package olts_fut

import "rhino-common/domain_error"

func (c *OltsFutChannel) Reset(force bool) (de *domain_error.Error) { 
	c.kafkaClient.reset()
	return
}