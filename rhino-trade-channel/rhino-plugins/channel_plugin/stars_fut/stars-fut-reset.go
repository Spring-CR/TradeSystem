package stars_fut

import "rhino-common/domain_error"

func (c *StarsFutChannel) Reset(force bool) (de *domain_error.Error) { 
	c.kafkaClient.reset()
	return
}