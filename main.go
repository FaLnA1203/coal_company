package main

import (
	"coal_company/api"
	"coal_company/company"
	"context"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	company := company.NewCompany(ctx)
	handler := api.NewAPIHandler(company)
	server := api.NewServer(&handler)
	if err := server.Run(); err != nil {
		println(err.Error())
	}
}
