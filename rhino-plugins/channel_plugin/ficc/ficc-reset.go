package ficc

import "rhino-common/domain_error"

func (c *FiccChannel) Reset(force bool) (de *domain_error.Error) { 
	c.kafkaClient.reset()
	return
}