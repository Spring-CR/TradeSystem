package order_domain

import (
	"rhino-core/domain_cfg"
	"rhino-core/order_domain/order_archive"
	"rhino-core/order_domain/order_cache"
	"rhino-core/order_domain/order_position_manager"
	"rhino-core/order_domain/order_purge"
	"rhino-core/order_domain/schedule"
)

type orderOrderArchiverAndPurger struct {
	orderArchiver *order_archive.OrderArchiver
	orderPurger   *order_purge.OrderPurger
}

func (o *orderOrderArchiverAndPurger) start() {
	o.orderArchiver.Start()
	o.orderPurger.Start()
}

type OrderOrderArchivingAndPurgingManager struct {
	orderOrderArchiverAndPurgers []*orderOrderArchiverAndPurger
}

func NewOrderOrderArchivingAndPurgingManager(applicationCfg *domain_cfg.ApplicationCfg, orderCache *order_cache.OrderCache, scheduleAdapter schedule.ScheduleAdapter, positionManager *order_position_manager.PositionManager) *OrderOrderArchivingAndPurgingManager {
	inst := &OrderOrderArchivingAndPurgingManager{}

	appArchivingCfgItems := applicationCfg.GetApplicationArchivingCfgItems()
	if len(appArchivingCfgItems) > 0 {
		for _, appArchivingCfgItem := range appArchivingCfgItems {
			orderArchiver := order_archive.NewOrderArchiver(applicationCfg, appArchivingCfgItem, orderCache, scheduleAdapter)
			orderPurger := order_purge.NewOrderPurger(applicationCfg, orderCache, orderArchiver, scheduleAdapter, positionManager)
			inst.orderOrderArchiverAndPurgers = append(inst.orderOrderArchiverAndPurgers, &orderOrderArchiverAndPurger{orderArchiver: orderArchiver, orderPurger: orderPurger})
		}
	} else {
		orderArchiver := order_archive.NewOrderArchiver(applicationCfg, nil, orderCache, scheduleAdapter)
		orderPurger := order_purge.NewOrderPurger(applicationCfg, orderCache, orderArchiver, scheduleAdapter, positionManager)
		inst.orderOrderArchiverAndPurgers = append(inst.orderOrderArchiverAndPurgers, &orderOrderArchiverAndPurger{orderArchiver: orderArchiver, orderPurger: orderPurger})
	}

	return inst
}

func (o *OrderOrderArchivingAndPurgingManager) Start() {
	for _, item := range o.orderOrderArchiverAndPurgers {
		item.start()
	}
}

func (o *OrderOrderArchivingAndPurgingManager) ForceArchiving() {
	for _, item := range o.orderOrderArchiverAndPurgers {
		item.orderArchiver.ForceArchiving()
	}
}

func (o *OrderOrderArchivingAndPurgingManager) ForcePurging() {
	for _, item := range o.orderOrderArchiverAndPurgers {
		item.orderPurger.ForcePurging()
	}
}
