package envutil

func IsDebug()bool{
	envMap := Env()
	debug := envMap["DEBUG"]
	return debug == "1" || debug == "true"
}
