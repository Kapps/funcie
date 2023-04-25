package main

func main() {
	panic("not implemented")
	/*config := bastion.NewConfigFromEnvironment()
	redisClient := &redis.Client{}
	//publisher := r.NewPublisher(redisClient, config.RequestChannel)
	handler := bastion.NewRequestHandler(publisher, config.RequestTtl)
	server := bastion.NewServer(config.ListenAddress, handler)

	err := server.Listen()
	if err != nil {
		panic("stopped listening on server")
	}*/
}
