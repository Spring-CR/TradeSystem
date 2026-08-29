package rhino_plugins

import (
	_ "rhino-plugins/api_plugin/atp_cbf"
	_ "rhino-plugins/api_plugin/ficc"
	_ "rhino-plugins/api_plugin/stars_fut"
	_ "rhino-plugins/api_plugin/titans" // titans API插件
	_ "rhino-plugins/channel_plugin/atp_cbf"
	_ "rhino-plugins/channel_plugin/citic_qfii" // 中信qfii交易通道插件
	_ "rhino-plugins/channel_plugin/ficc"
	_ "rhino-plugins/channel_plugin/gfgfix"
	_ "rhino-plugins/channel_plugin/stars_fut"
	_ "rhino-plugins/data_sync_plugin/ficc"
	_ "rhino-plugins/executor_plugin/ficc"
	_ "rhino-plugins/fix_api_plugin/ficc"
	_ "rhino-plugins/order_capital_plugin/ficc"
	_ "rhino-plugins/order_position_plugin/ficc_v2"
	_ "rhino-plugins/order_status_plugin/ficc"
	_ "rhino-plugins/schedule_plugin/ficc"
)
