package main

func main() {
	//cfg, err := GetConfig()
	cfg := Config{
		update: true,
		asgID:  "-sandbox-",
		debug:  true,
	}
	//if err != nil {
	//	panic(err)
	//}
	log := cfg.GetLogger()
	log.Info().Msgf("config: %+v", cfg)
}
